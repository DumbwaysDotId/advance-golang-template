### Table of Contents

- [GORM Relation Has Many](#gorm-relation-has-many)
  - [Repository](#repository)

---

# GORM Relation Has Many

Reference: [Official GORM Website](https://gorm.io/docs/has_many.html)

## Relation

For this section, example Has Many relation:

- `User` &rarr; `Product`: to get User Product

## Repository

- Inside `repositories` folder, in `users.go` file write this below code

  > File: `repositories/users.go`

  ```go
  func (r *repository) FindUsers() ([]models.User, error) {
    var users []models.User
    err := r.db.Preload("Profile").Preload("Products").Find(&users).Error // add this code

    return users, err
  }

  func (r *repository) GetUser(ID int) (models.User, error) {
    var user models.User
    err := r.db.Preload("Profile").Preload("Products").First(&user, ID).Error // add this code

    return user, err
  }
  ```

  \*In this case, just add `Preload` to make relation
