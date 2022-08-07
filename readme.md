### Table of Contents

- [GORM Relation Many to Many](#gorm-relation-has-many)
  - [Repository](#repository)

---

# GORM Relation Many to Many

Reference: [Official GORM Website](https://gorm.io/docs/many_to_many.html)

## Relation

For this section, example Many to Many relation:

- `Product` &rarr; `Category`: to get Product Category

## Repository

- Inside `repositories` folder, in `product.go` file write this below code

  > File: `repositories/product.go`

  ```go
  func (r *repository) FindProducts() ([]models.Product, error) {
    var products []models.Product
    err := r.db.Preload("User").Preload("Category").Find(&products).Error // add this code

    return products, err
  }

  func (r *repository) GetProduct(ID int) (models.Product, error) {
    var product models.Product
    // not yet using category relation, cause this step doesnt Belong to Many
    err := r.db.Preload("User").Preload("Category").First(&product, ID).Error // add this code

    return product, err
  }
  ```

  \*In this case, just add `Preload` to make relation
