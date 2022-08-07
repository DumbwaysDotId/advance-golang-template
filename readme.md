### Table of Contents

- [GORM Relation belongs to](#gorm-relation-has-one)
  - [Repository](#repository)

---

# GORM Relation Has One

Reference: [Official GORM Website](https://gorm.io/docs/has_one.html)

## Relation

For this section, example Has One relation:

- `User` &rarr; `Profile`: to get User Profile

## Repository

- Inside `repositories` folder, in `users.go` file write this below code

  > File: `repositories/users.go`

  ```go
  func (r *repository) FindUsers() ([]models.User, error) {
    var users []models.User
    err := r.db.Preload("Profile").Find(&users).Error // add this code

    return users, err
  }

  func (r *repository) GetUser(ID int) (models.User, error) {
    var user models.User
    err := r.db.Preload("Profile").First(&user, ID).Error // add this code

    return user, err
  }
  ```

  \*In this case, just add `Preload` to make relation
