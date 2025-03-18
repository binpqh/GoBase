package repository

import (
	"database/sql"
	"reflect"

	"github.com/binpqh/GoBase/entity"
	"github.com/binpqh/GoBase/query"
)

// Repository defines generic CRUD operations.
type Repository[TEntity any, TKey any] interface {
	GetByID(id TKey) (*TEntity, error)                                // Fetch entity by ID
	GetAll() ([]TEntity, error)                                       // Fetch all entities
	GetByExpression(expression func(TEntity) bool) ([]TEntity, error) // Fetch entities by condition
	Create(entity *TEntity) error                                     // Insert new entity
	Update(entity *TEntity) error                                     // Update entity by ID
	Delete(id TKey) error                                             // Delete entity by ID
}

// GenericRepository provides default implementations of Repository.
type GenericRepository[T any, TKey any] struct {
	db *sql.DB
}

// NewGenericRepository creates a new instance of GenericRepository.
func NewGenericRepository[T any, TKey any](db *sql.DB) *GenericRepository[T, TKey] {
	return &GenericRepository[T, TKey]{db: db}
}

// GetByID retrieves an entity by its primary key.
//
// Params:
//   - id: Primary key value.
//
// Returns:
//   - *T: The entity if found.
//   - error: Error if the query fails or no entity is found.
func (r *GenericRepository[T, TKey]) GetByID(id TKey) (*T, error) {
	entityInstance := new(T)

	qb := query.NewQueryBuilder[T]().
		WhereEqual(entity.Field[T]("ID"), id).
		Limit(1)

	sqlQuery, args := qb.Build()

	row := r.db.QueryRow(sqlQuery, args...)
	err := row.Scan(entityInstance)
	if err != nil {
		return nil, err
	}
	return entityInstance, nil
}

// GetAll retrieves all records of type T.
//
// Returns:
//   - []T: A slice of entities.
//   - error: Error if the query fails.
func (r *GenericRepository[T, TKey]) GetAll() ([]T, error) {
	var results []T

	qb := query.NewQueryBuilder[T]() // SELECT * FROM table
	sql, args := qb.Build()

	rows, err := r.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entity T
		if err := rows.Scan(&entity); err != nil {
			return nil, err
		}
		results = append(results, entity)
	}
	return results, nil
}

// Create inserts a new entity into the database.
//
// Params:
//   - entity: The entity to be inserted.
//
// Returns:
//   - error: Error if insertion fails.
func (r *GenericRepository[T, TKey]) Create(entity *T) error {
	qb := query.NewQueryBuilder[T]().Insert(*entity)
	sql, args := qb.Build()

	_, err := r.db.Exec(sql, args...)
	return err
}

// Update modifies an existing entity based on its ID.
//
// Params:
//   - entity: The entity with updated values.
//
// Returns:
//   - error: Error if the update fails.
func (r *GenericRepository[T, TKey]) Update(ent *T) error {
	idValue := reflect.ValueOf(ent).Elem().FieldByName("ID").Interface()

	qb := query.NewQueryBuilder[T]().
		Update(*ent).
		WhereEqual(entity.Field[T]("ID"), idValue)

	sql, args := qb.Build()
	_, err := r.db.Exec(sql, args...)
	return err
}

// Delete removes an entity from the database by its ID.
//
// Params:
//   - id: Primary key value of the entity to delete.
//
// Returns:
//   - error: Error if deletion fails.
func (r *GenericRepository[T, TKey]) Delete(id TKey) error {
	qb := query.NewQueryBuilder[T]().
		Delete().
		WhereEqual(entity.Field[T]("ID"), id)

	sql, args := qb.Build()
	_, err := r.db.Exec(sql, args...)
	return err
}
