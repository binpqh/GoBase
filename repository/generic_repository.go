package repository

import (
	"database/sql"
)

type Repository[TEntity any, TKey any] interface {
	GetByID(id TKey) (*TEntity, error)
	GetAll() ([]TEntity, error)
	GetByExpression(expression func(TEntity) bool) ([]TEntity, error)
	Create(entity *TEntity) error
	Update(entity *TEntity) error
	Delete(id TKey) error
}

type GenericRepository[T any, TKey any] struct {
	db *sql.DB
}

func NewGenericRepository[T any, TKey any](db *sql.DB) *GenericRepository[T, TKey] {
	return &GenericRepository[T, TKey]{db: db}
}

// GetByID
func (r *GenericRepository[T, TKey]) GetByID(id TKey) (*T, error) {
	entity := new(T)
	fields := GetField[*T]().(map[string]string)

	qb := NewQueryBuilder[T]().
		Select(fields["ID"]).
		WhereEqual(fields["ID"], id).
		Limit(1)

	sql, args := qb.Build()

	row := r.db.QueryRow(sql, args...)
	err := row.Scan(entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// GetAll
func (r *GenericRepository[T, TKey]) GetAll() ([]T, error) {
	var results []T
	fields := GetModelFields[*T]().(map[string]string)

	qb := NewQueryBuilder[T]().
		Select(fields["ID"])

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

// Create
func (r *GenericRepository[T, TKey]) Create(entity *T) error {
	// Build insert query using reflection
	return nil // Implement insert logic
}

// Update
func (r *GenericRepository[T, TKey]) Update(entity *T) error {
	// Build update query using reflection
	return nil // Implement update logic
}

// Delete
func (r *GenericRepository[T, TKey]) Delete(id TKey) error {
	entity := new(T)
	fields := GetModelFields[*T]().(map[string]string)

	qb := NewQueryBuilder[T]().
		WhereEqual(fields["ID"], id)

	sql, args := qb.Build()

	_, err := r.db.Exec(sql, args...)
	return err
}
