# Gobase - Generic Repository & Query Builder for Go

## 📌 1. Overview  
**Gobase** is a Go module that provides a **generic repository pattern** with a **query builder**.  
It allows you to **build SQL queries dynamically** and **handle database operations** easily.  

✅ **Automatically map struct fields → column names**  
✅ **Flexible queries using `WhereEqual`, `Join`, `OrderBy`, etc.**  
✅ **Generic Repository for CRUD operations**  
✅ **High performance (field mapping cache, minimal reflection)**  

---

## 📌 2. Installation  

```sh
go get github.com/binpqh/gobase
```

## 📌 3. Usage
### 🚀 3.1. Register Models
Before using the Query Builder, register models to cache field mappings.
(Run this once at application startup)
```go
import "github.com/binpqh/gobase/repository"

func main() {
    repository.RegisterModels(models.User{}, models.Order{}, models.Product{})
}
```

### 🚀 3.2. Using Query Builder
Once models are registered, you can build queries dynamically.

🛠 Example: Simple Query
```go
User := models.User{}

qb := repository.NewQueryBuilder[models.User]().
    Select(repository.GetField(User, "ID"), repository.GetField(User, "Name")).
    WhereEqual(repository.GetField(User, "Age"), 25).
    OrderByASC(repository.GetField(User, "Name")).
    Limit(10)

sql, args := qb.Build()
fmt.Println("SQL:", sql)
fmt.Println("Args:", args)
```
```sql
🔹 Output:
SQL: SELECT id, name FROM user WHERE age = ? ORDER BY name ASC LIMIT ?
Args: [25 10]
```

### 🚀 3.3. Using Generic Repository
Instead of writing SQL manually, use Generic Repository for database operations.

🛠 Define Repository Interface
```go
type Repository[T any, TKey any] interface {
    GetByID(id TKey) (*T, error)
    GetAll() ([]T, error)
    GetByExpression(expression func(T) bool) ([]T, error)
    Create(entity *T) error
    Update(entity *T) error
    Delete(id TKey) error
}
```

🛠 Implement Repository (UserRepository.go)
```go
package repositories

import (
    "database/sql"
    "github.com/your-repo/gobase/repository"
    "your-app/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
    User := models.User{}

    qb := repository.NewQueryBuilder[models.User]().
        Select(repository.GetField(User, "ID"), repository.GetField(User, "Name")).
        WhereEqual(repository.GetField(User, "ID"), id).
        Limit(1)

    sql, args := qb.Build()
    
    row := r.db.QueryRow(sql, args...)
    var user models.User
    if err := row.Scan(&user.ID, &user.Name); err != nil {
        return nil, err
    }
    return &user, nil
}
```
### 📌 4. Example: Usage in a Service
```go
package services

import (
    "database/sql"
    "fmt"
    "your-app/models"
    "your-app/repositories"
)

func ExampleUsage(db *sql.DB) {
    userRepo := repositories.NewUserRepository(db)

    // Fetch a user by ID
    user, err := userRepo.GetByID(1)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("User:", user)
}
```

### 📌 5. Summary
🚀 Gobase makes database queries simple and efficient:
✅ Dynamic Query Builder
✅ Generic Repository for CRUD operations
✅ Optimized Performance (Caching & Minimal Reflection)

🔥 Happy Coding