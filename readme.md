### Table of Contents

- [Handle Upload File](#handle-upload-file)
  - [Introduction](#introduction)
  - [Package](#Package)
  - [Routes](#routes)
  - [Handler](#Handler)
  - [Folder Store File](#folder-store-file)
  - [DotEnv](#dotenv)

---

# Handle Upload File

## Introduction

For this section:

- Handle File Upload for `Create Product` data

## Package

- Inside `pkg` folder, in `middleware` folder, inside it create `uploadFile.go` file, and write this below code

  > File: `pkg/middleware/uploadFile.go`

  ```go
  package middleware

  import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
  )

  func UploadFile(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      // Upload file
      // FormFile returns the first file for the given key `myFile`
      // it also returns the FileHeader so we can get the Filename,
      // the Header and the size of the file
      file, _, err := r.FormFile("image")

      if err != nil {
        fmt.Println(err)
        json.NewEncoder(w).Encode("Error Retrieving the File")
        return
      }
      defer file.Close()
      // fmt.Printf("Uploaded File: %+v\n", handler.Filename)
      // fmt.Printf("File Size: %+v\n", handler.Size)
      // fmt.Printf("MIME Header: %+v\n", handler.Header)
      const MAX_UPLOAD_SIZE = 10 << 20 // 10MB
      // Parse our multipart form, 10 << 20 specifies a maximum
      // upload of 10 MB files.
      r.ParseMultipartForm(MAX_UPLOAD_SIZE)
      if r.ContentLength > MAX_UPLOAD_SIZE {
        w.WriteHeader(http.StatusBadRequest)
        response := Result{Code: http.StatusBadRequest, Message: "Max size in 1mb"}
        json.NewEncoder(w).Encode(response)
        return
      }

      // Create a temporary file within our temp-images directory that follows
      // a particular naming pattern
      tempFile, err := ioutil.TempFile("uploads", "image-*.png")
      if err != nil {
        fmt.Println(err)
        fmt.Println("path upload error")
        json.NewEncoder(w).Encode(err)
        return
      }
      defer tempFile.Close()

      // read all of the contents of our uploaded file into a
      // byte array
      fileBytes, err := ioutil.ReadAll(file)
      if err != nil {
        fmt.Println(err)
      }

      // write this byte array to our temporary file
      tempFile.Write(fileBytes)

      data := tempFile.Name()
      filename := data[8:] // split uploads/

      // add filename to ctx
      ctx := context.WithValue(r.Context(), "dataFile", filename)
      next.ServeHTTP(w, r.WithContext(ctx))
    })
  }
  ```

## Routes

- In `routes` folder, inside `product.go` file, write `uploadFile` middleware on `/product` route

  > File: `routes/product.go`

  ```go
  package routes

  import (
    "dumbmerch/handlers"
    "dumbmerch/pkg/middleware"
    "dumbmerch/pkg/mysql"
    "dumbmerch/repositories"

    "github.com/gorilla/mux"
  )

  func ProductRoutes(r *mux.Router) {
    productRepository := repositories.RepositoryProduct(mysql.DB)
    h := handlers.HandlerProduct(productRepository)

    r.HandleFunc("/products", middleware.Auth(h.FindProducts)).Methods("GET")
    r.HandleFunc("/product/{id}", h.GetProduct).Methods("GET")
    r.HandleFunc("/product", middleware.Auth(middleware.UploadFile(h.CreateProduct))).Methods("POST") // add this code
  }
  ```

## Handler

- In `handlers` folder, inside `product.go` file, write get `filename` and store like this below code

  > File: `handlers/product.go`

  ```go
  func (h *handlerProduct) CreateProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // get data user token
    userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
    userId := int(userInfo["id"].(float64))

    dataContex := r.Context().Value("dataFile") // add this code
    filename := dataContex.(string) // add this code

    price, _ := strconv.Atoi(r.FormValue("price"))
    qty, _ := strconv.Atoi(r.FormValue("qty"))
    category_id, _ := strconv.Atoi(r.FormValue("category_id"))
    request := productdto.ProductRequest{
      Name:       r.FormValue("name"),
      Desc:       r.FormValue("desc"),
      Price:      price,
      Qty:        qty,
      CategoryID: category_id,
    }

    validation := validator.New()
    err := validation.Struct(request)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    product := models.Product{
      Name:   request.Name,
      Desc:   request.Desc,
      Price:  request.Price,
      Image:  filename, // add this code
      Qty:    request.Qty,
      UserID: userId,
    }

    product, err = h.ProductRepository.CreateProduct(product)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    product, _ = h.ProductRepository.GetProduct(product.ID)

    w.WriteHeader(http.StatusOK)
    response := dto.SuccessResult{Code: http.StatusOK, Data: product}
    json.NewEncoder(w).Encode(response)
  }
  ```

- Embed Path file in `FindProducts` and `GetProduct` method
  > File: `handlers/product.go`
  - Create `path_file` Global variable
    ```go
    var path_file = "http://localhost:5000/uploads/"
    ```
  - `FindProducts` method
    ```go
    for i, p := range products {
      products[i].Image = path_file + p.Image
    }
    ```
  - `GetProduct` method
    ```go
    product.Image = path_file + product.Image
    ```

## Folder Store File

- Create `uploads` folder

  > File: `./uploads`

- Add this below code to make `uploads` can be used another client

  > File: `main.go`

  ```go
  package main

  import (
    "dumbmerch/database"
    "dumbmerch/pkg/mysql"
    "dumbmerch/routes"
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
  )

  func main() {
    // initial DB
    mysql.DatabaseInit()

    // run migration
    database.RunMigration()

    r := mux.NewRouter()

    routes.RouteInit(r.PathPrefix("/api/v1").Subrouter())

    //path file
    r.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads")))) // add this code

    fmt.Println("server running localhost:5000")
    http.ListenAndServe("localhost:5000", r)
  }
  ```

## DotEnv

- Installation

  ```bash
  go get github.com/joho/godotenv
  ```

- Create `.env` file and write this below code

  > File: `.env`

  ```env
  SECRET_KEY=bolehapaaja
  ```

- In `main.go` file import `godotenv` and Init `godotenv` inside `main` function like this below code

  > File: `main.go`

  - Import `godotenv` package
    ```go
    import (
      // another package here ...
      "github.com/joho/godotenv" // import this package
    )
    ```
  - Init `godotenv`

    ```go
    func main() {

      	// env
        errEnv := godotenv.Load()
        if errEnv != nil {
          panic("Failed to load env file")
        }

        // Another code on this below ...
    }
    ```

- How to use Environment Variable, write this below code inside `jwt.go` file

  > File: `pkg/jwt/jwt.go`

  - Import `os` package
    ```go
    import (
      "fmt"
      "os" // import this package
      "github.com/golang-jwt/jwt/v4"
    )
    ```
  - Modify `SecretKey` variable like this below code
    ```go
    var SecretKey = os.Getenv("SECRET_KEY")
    ```
