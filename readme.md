### Table of Contents

- [GORM Relation belongs to](#gorm-relation-belongs-to)
  - [Handlers](#handlers)
  - [Repository](#repository)
  - [Routes](#routes)

---

# GORM Relation Belongs to

Reference: [Official GORM Website](https://gorm.io/docs/belongs_to.html)

## Relation

For this section, example Belongs To relation:

- `Profile` &rarr; `User`: to get Profile User
- `Product` &rarr; `User`: to get Product Owner

## Handlers

- Inside `handlers` folder, create `profile.go` file, and write this below code

  > File: `handlers/profile.go`

  ```go
  package handlers

  import (
    profiledto "dumbmerch/dto/profile"
    dto "dumbmerch/dto/result"
    "dumbmerch/models"
    "dumbmerch/repositories"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
  )

  type handlerProfile struct {
    ProfileRepository repositories.ProfileRepository
  }

  func HandlerProfile(ProfileRepository repositories.ProfileRepository) *handlerProfile {
    return &handlerProfile{ProfileRepository}
  }

  func (h *handlerProfile) GetProfile(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    var profile models.Profile
    profile, err := h.ProfileRepository.GetProfile(id)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    w.WriteHeader(http.StatusOK)
    response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseProfile(profile)}
    json.NewEncoder(w).Encode(response)
  }

  func convertResponseProfile(u models.Profile) profiledto.ProfileResponse {
    return profiledto.ProfileResponse{
      ID:      u.ID,
      Phone:   u.Phone,
      Gender:  u.Gender,
      Address: u.Address,
      UserID:  u.UserID,
      User:    u.User,
    }
  }
  ```

- Inside `handlers` folder, create `product.go` file, and write this below code

  > File: `handlers/product.go`

  ```go
  package handlers

  import (
    productdto "dumbmerch/dto/product"
    dto "dumbmerch/dto/result"
    "dumbmerch/models"
    "dumbmerch/repositories"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-playground/validator/v10"
    "github.com/gorilla/mux"
  )

  type handlerProduct struct {
    ProductRepository repositories.ProductRepository
  }

  func HandlerProduct(ProductRepository repositories.ProductRepository) *handlerProduct {
    return &handlerProduct{ProductRepository}
  }

  func (h *handlerProduct) FindProducts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    products, err := h.ProductRepository.FindProducts()
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    w.WriteHeader(http.StatusOK)
    response := dto.SuccessResult{Code: http.StatusOK, Data: products}
    json.NewEncoder(w).Encode(response)
  }

  func (h *handlerProduct) GetProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    id, _ := strconv.Atoi(mux.Vars(r)["id"])

    var product models.Product
    product, err := h.ProductRepository.GetProduct(id)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    w.WriteHeader(http.StatusOK)
    response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseProduct(product)}
    json.NewEncoder(w).Encode(response)
  }

  func (h *handlerProduct) CreateProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    request := new(productdto.ProductRequest)
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
      w.WriteHeader(http.StatusBadRequest)
      response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
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
      Image:  request.Image,
      Qty:    request.Qty,
      UserID: request.UserID,
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

  func convertResponseProduct(u models.Product) models.ProductResponse {
    return models.ProductResponse{
      Name:     u.Name,
      Desc:     u.Desc,
      Price:    u.Price,
      Image:    u.Image,
      Qty:      u.Qty,
      User:     u.User,
      Category: u.Category,
    }
  }
  ```

## Repository

- Inside `repositories` folder, create `profile.go` file, and write this below code

  > File: `repositories/profile.go`

  ```go
  package repositories

  import (
    "dumbmerch/models"

    "gorm.io/gorm"
  )

  type ProfileRepository interface {
    GetProfile(ID int) (models.Profile, error)
  }

  func RepositoryProfile(db *gorm.DB) *repository {
    return &repository{db}
  }

  func (r *repository) GetProfile(ID int) (models.Profile, error) {
    var profile models.Profile
    err := r.db.Preload("User").First(&profile, ID).Error

    return profile, err
  }
  ```

- Inside `repositories` folder, create `product.go` file, and write this below code

  > File: `repositories/product.go`

  ```go
  package repositories

  import (
    "dumbmerch/models"

    "gorm.io/gorm"
  )

  type ProductRepository interface {
    FindProducts() ([]models.Product, error)
    GetProduct(ID int) (models.Product, error)
    CreateProduct(product models.Product) (models.Product, error)
  }

  func RepositoryProduct(db *gorm.DB) *repository {
    return &repository{db}
  }

  func (r *repository) FindProducts() ([]models.Product, error) {
    var products []models.Product
    err := r.db.Preload("User").Find(&products).Error

    return products, err
  }

  func (r *repository) GetProduct(ID int) (models.Product, error) {
    var product models.Product
    // not yet using category relation, cause this step doesnt Belong to Many
    err := r.db.Preload("User").First(&product, ID).Error

    return product, err
  }

  func (r *repository) CreateProduct(product models.Product) (models.Product, error) {
    err := r.db.Create(&product).Error

    return product, err
  }
  ```

## Routes

- Inside `routes` folder, create `profile.go` file, and write this below code

  > File: `routes/profile.go`

  ```go
  package routes

  import (
    "dumbmerch/handlers"
    "dumbmerch/pkg/mysql"
    "dumbmerch/repositories"

    "github.com/gorilla/mux"
  )

  func ProfileRoutes(r *mux.Router) {
    profileRepository := repositories.RepositoryProfile(mysql.DB)
    h := handlers.HandlerProfile(profileRepository)

    r.HandleFunc("/profile/{id}", h.GetProfile).Methods("GET")
  }
  ```

- Inside `routes` folder, create `profile.go` file, and write this below code

  > File: `routes/product.go`

  ```go
  package routes

  import (
    "dumbmerch/handlers"
    "dumbmerch/pkg/mysql"
    "dumbmerch/repositories"

    "github.com/gorilla/mux"
  )

  func ProductRoutes(r *mux.Router) {
    productRepository := repositories.RepositoryProduct(mysql.DB)
    h := handlers.HandlerProduct(productRepository)

    r.HandleFunc("/products", h.FindProducts).Methods("GET")
    r.HandleFunc("/product/{id}", h.GetProduct).Methods("GET")
    r.HandleFunc("/product", h.CreateProduct).Methods("POST")
  }
  ```

- On `routes.go` file, write `ProfileRoutes` and `ProductRoutes`

  > File: `routes/routes.go`

  ```go
  package routes

  import (
    "github.com/gorilla/mux"
  )

  func RouteInit(r *mux.Router) {
    UserRoutes(r)
    ProfileRoutes(r) // Add this code
    ProductRoutes(r) // Add this code
  }
  ```
