### Table of Contents

- [GORM Relation belongs to](#gorm-relation-has-one)
  - [Repository](#repository)

---

# GORM Relation Has One

Reference: [Official GORM Website](https://gorm.io/docs/has_one.html)

## Relation

For this section, example Has One relation:

- User &rarr; Profile: to get User Profile

## Repository

- Inside `handlers` folder, create `profile.go` file, and write this below code

  > File: `handlers/profile.go`

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

- Inside `handlers` folder, create `product.go` file, and write this below code

  > File: `handlers/product.go`

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
