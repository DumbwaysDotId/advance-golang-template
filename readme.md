### Table of Contents

- [Authentication JWT](#authentication-jwt)
  - [Introduction](#introduction)
  - [Installation](#installation)
  - [Package](#Package)
  - [Handler](#Handler)
  - [Repository](#repository)
  - [Routes](#routes)

---

# Authentication JWT

Reference: [Golang JWT](https://github.com/golang-jwt/jwt)

## Introduction

For this section:

- Generate Token using JWT if `User Login`
- Verify Token and Get User Data if `Create Product Data`

## Installation

- Golang Json Web Token (JWT)

  ```bash
  go get -u github.com/golang-jwt/jwt/v4
  ```

## Package

- Inside `pkg` folder, create `jwt` folder, inside it create `jwt.go` file, and write this below code

  > File: `pkg/jwt/jwt.go`

  ```go
  package jwtToken

  import (
    "fmt"

    "github.com/golang-jwt/jwt/v4"
  )

  var SecretKey = "SECRET_KEY"

  func GenerateToken(claims *jwt.MapClaims) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    webtoken, err := token.SignedString([]byte(SecretKey))
    if err != nil {
      return "", err
    }

    return webtoken, nil
  }

  func VerifyToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
      if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
      }
      return []byte(SecretKey), nil
    })

    if err != nil {
      return nil, err
    }
    return token, nil
  }

  func DecodeToken(tokenString string) (jwt.MapClaims, error) {
    token, err := VerifyToken(tokenString)
    if err != nil {
      return nil, err
    }

    claims, isOk := token.Claims.(jwt.MapClaims)
    if isOk && token.Valid {
      return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
  }
  ```

- Inside `pkg` folder, create `middleware` folder, inside it create `auth.go` file, and write this below code

  > File: `pkg/middleware/auth.go`

  ```go
  package middleware

  import (
    "context"
    dto "dumbmerch/dto/result"
    jwtToken "dumbmerch/pkg/jwt"
    "encoding/json"
    "net/http"
    "strings"
  )

  type Result struct {
    Code    int         `json:"code"`
    Data    interface{} `json:"data"`
    Message string      `json:"message"`
  }

  func Auth(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "application/json")

      token := r.Header.Get("Authorization")

      if token == "" {
        w.WriteHeader(http.StatusUnauthorized)
        response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "unauthorized"}
        json.NewEncoder(w).Encode(response)
        return
      }

      token = strings.Split(token, " ")[1]
      claims, err := jwtToken.DecodeToken(token)

      if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        response := Result{Code: http.StatusUnauthorized, Message: "unauthorized"}
        json.NewEncoder(w).Encode(response)
        return
      }

      ctx := context.WithValue(r.Context(), "userInfo", claims)
      r = r.WithContext(ctx)
      next.ServeHTTP(w, r.WithContext(ctx))
    })
  }
  ```

## Handler

- Inside `handlers` folder, On `auth.go` file and write `Login` Function like this below code

  > File: `handlers/auth.go`

  ```go
  func (h *handlerAuth) Login(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    request := new(authdto.LoginRequest)
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
      w.WriteHeader(http.StatusBadRequest)
      response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    user := models.User{
      Email:    request.Email,
      Password: request.Password,
    }

    // Check email
    user, err := h.AuthRepository.Login(user.Email)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
      json.NewEncoder(w).Encode(response)
      return
    }

    // Check password
    isValid := bcrypt.CheckPasswordHash(request.Password, user.Password)
    if !isValid {
      w.WriteHeader(http.StatusBadRequest)
      response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "wrong email or password"}
      json.NewEncoder(w).Encode(response)
      return
    }

    //generate token
    claims := jwt.MapClaims{}
    claims["id"] = user.ID
    claims["exp"] = time.Now().Add(time.Hour * 2).Unix() // 2 hours expired

    token, errGenerateToken := jwtToken.GenerateToken(&claims)
    if errGenerateToken != nil {
      log.Println(errGenerateToken)
      fmt.Println("Unauthorize")
      return
    }

    loginResponse := authdto.LoginResponse{
      Name:     user.Name,
      Email:    user.Email,
      Password: user.Password,
      Token:    token,
    }

    w.Header().Set("Content-Type", "application/json")
    response := dto.SuccessResult{Code: http.StatusOK, Data: loginResponse}
    json.NewEncoder(w).Encode(response)

  }
  ```

  - Inside `handlers` folder, On `product.go` file and write `CreateProduct` Function like this below code

  > File: `handlers/product.go`

  ```go

  ```

## Repository

- Inside `repositories` folder, create `auth.go` file and write this below code

  > File: `repositories/auth.go`

  ```go
  func (h *handlerProduct) CreateProduct(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // get data user token
    userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
    userId := int(userInfo["id"].(float64))

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

## Routes

- Inside `routes` folder, in `auth.go` file and write `Login` route this below code

  > File: `routes/auth.go`

  ```go
  package routes

  import (
    "dumbmerch/handlers"
    "dumbmerch/pkg/mysql"
    "dumbmerch/repositories"

    "github.com/gorilla/mux"
  )

  func AuthRoutes(r *mux.Router) {
    userRepository := repositories.RepositoryUser(mysql.DB)
    h := handlers.HandlerAuth(userRepository)

    r.HandleFunc("/register", h.Register).Methods("POST")
    r.HandleFunc("/login", h.Login).Methods("POST") // add this code
  }
  ```

- Inside `routes` folder, in `product.go` file and write `product` route with `middleware` like this below code

  > File: `routes/auth.go`

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

    r.HandleFunc("/products", middleware.Auth(h.FindProducts)).Methods("GET") // add this code
    r.HandleFunc("/product/{id}", h.GetProduct).Methods("GET")
    r.HandleFunc("/product", middleware.Auth(h.CreateProduct)).Methods("POST") // add this code
  }
  ```
