package query

import (
	"reflect"
	"strings"

	"github.com/binpqh/GoBase/entity"
)

// QueryBuilder constructs SQL queries dynamically for entities implementing EntityBase.
type QueryBuilder[T any] struct {
	selectFields []string
	whereClauses []string
	whereArgs    []interface{}
	joins        []string
	limit        int
	orderBy      []string
	insertFields []string
	insertValues []interface{}
	updateFields []string
}

// NewQueryBuilder creates a new QueryBuilder instance for the specified entity type.
//
// Returns:
//   - *QueryBuilder[T]: A new instance of QueryBuilder.
func NewQueryBuilder[T any]() *QueryBuilder[T] {
	return &QueryBuilder[T]{}
}

// Insert generates an INSERT INTO statement.
//
// Params:
//   - data: The entity data to insert (must be a struct).
//
// Returns:
//   - *QueryBuilder[T]: The updated QueryBuilder instance.
func (qb *QueryBuilder[T]) Insert(data T) *QueryBuilder[T] {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		qb.insertFields = append(qb.insertFields, field.Name)
		qb.insertValues = append(qb.insertValues, value)
	}
	return qb
}

// Update generates an UPDATE statement.
//
// Params:
//   - data: The entity data to update.
//
// Returns:
//   - *QueryBuilder[T]: The updated QueryBuilder instance.
func (qb *QueryBuilder[T]) Update(data T) *QueryBuilder[T] {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		qb.updateFields = append(qb.updateFields, field.Name+" = ?")
		qb.whereArgs = append(qb.whereArgs, value)
	}
	return qb
}

// Delete generates a DELETE statement.
//
// Returns:
//   - *QueryBuilder[T]: The updated QueryBuilder instance.
func (qb *QueryBuilder[T]) Delete() *QueryBuilder[T] {
	return qb
}

// Select adds fields to the SELECT clause of the query.
//
// Params:
//   - fields: One or more entity fields to be selected.
//
// Returns:
//   - *QueryBuilder[T]: The updated QueryBuilder instance.
func (qb *QueryBuilder[T]) Select(fields ...entity.Field[T]) *QueryBuilder[T] {
	for _, field := range fields {
		qb.selectFields = append(qb.selectFields, string(field))
	}
	return qb
}

// WhereEqual adds a WHERE condition filtering results where a field equals a given value.
//
// Params:
//   - field: The entity field to filter by.
//   - value: The value to compare against.
//
// Returns:
//   - *QueryBuilder[T]: The updated QueryBuilder instance.
func (qb *QueryBuilder[T]) WhereEqual(field entity.Field[T], value interface{}) *QueryBuilder[T] {
	qb.whereClauses = append(qb.whereClauses, string(field)+" = ?")
	qb.whereArgs = append(qb.whereArgs, value)
	return qb
}

// OrderByASC adds an ORDER BY clause sorting results in ascending order.
//
// Params:
//   - field: The entity field to sort by.
//
// Returns:
//   - *QueryBuilder[T]: The updated QueryBuilder instance.
func (qb *QueryBuilder[T]) OrderByASC(field entity.Field[T]) *QueryBuilder[T] {
	qb.orderBy = append(qb.orderBy, string(field)+" ASC")
	return qb
}

// OrderByDESC adds an ORDER BY clause sorting results in descending order.
//
// Params:
//   - field: The entity field to sort by.
//
// Returns:
//   - *QueryBuilder[T]: The updated QueryBuilder instance.
func (qb *QueryBuilder[T]) OrderByDESC(field entity.Field[T]) *QueryBuilder[T] {
	qb.orderBy = append(qb.orderBy, string(field)+" DESC")
	return qb
}

// Join adds a JOIN clause to the query.
//
// Params:
//   - table: The table name to join with.
//   - onCondition: The join condition (e.g., "users.id = orders.user_id").
//
// Returns:
//   - *QueryBuilder[T]: The updated QueryBuilder instance.
func (qb *QueryBuilder[T]) Join(table string, onCondition string) *QueryBuilder[T] {
	qb.joins = append(qb.joins, "JOIN "+table+" ON "+onCondition)
	return qb
}

// Limit sets the maximum number of rows to be returned.
//
// Params:
//   - limit: The maximum number of rows.
//
// Returns:
//   - *QueryBuilder[T]: The updated QueryBuilder instance.
func (qb *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
	qb.limit = limit
	return qb
}

// Build generates the final SQL query string and its arguments.
//
// Returns:
//   - string: The generated SQL query.
//   - []interface{}: The slice of arguments to be used with the query.
func (qb *QueryBuilder[T]) Build() (string, []interface{}) {
	var sqlBuilder strings.Builder

	var entityInstance T
	typeName := reflect.TypeOf(entityInstance).Name()
	tableName := strings.ToLower(typeName)

	const whereClause = " WHERE "

	if len(qb.insertFields) > 0 {
		// INSERT INTO table (col1, col2) VALUES (?, ?)
		sqlBuilder.WriteString("INSERT INTO " + tableName + " (")
		sqlBuilder.WriteString(strings.Join(qb.insertFields, ", "))
		sqlBuilder.WriteString(") VALUES (")
		placeholders := make([]string, len(qb.insertFields))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		sqlBuilder.WriteString(strings.Join(placeholders, ", "))
		sqlBuilder.WriteString(")")
		return sqlBuilder.String(), qb.insertValues
	}

	if len(qb.updateFields) > 0 {
		// UPDATE table SET col1 = ?, col2 = ? WHERE ...
		sqlBuilder.WriteString("UPDATE " + tableName + " SET ")
		sqlBuilder.WriteString(strings.Join(qb.updateFields, ", "))

		if len(qb.whereClauses) > 0 {
			sqlBuilder.WriteString(whereClause + strings.Join(qb.whereClauses, " AND "))
		}
		return sqlBuilder.String(), qb.whereArgs
	}

	if len(qb.whereClauses) > 0 && len(qb.insertFields) == 0 {
		// DELETE FROM table WHERE ...
		sqlBuilder.WriteString("DELETE FROM " + tableName)
		sqlBuilder.WriteString(whereClause + strings.Join(qb.whereClauses, " AND "))
		return sqlBuilder.String(), qb.whereArgs
	}

	// SELECT * FROM table
	sqlBuilder.WriteString("SELECT ")
	if len(qb.selectFields) > 0 {
		sqlBuilder.WriteString(strings.Join(qb.selectFields, ", "))
	} else {
		sqlBuilder.WriteString("*")
	}
	sqlBuilder.WriteString(" FROM " + tableName)

	if len(qb.joins) > 0 {
		sqlBuilder.WriteString(" " + strings.Join(qb.joins, " "))
	}

	if len(qb.whereClauses) > 0 {
		sqlBuilder.WriteString(whereClause + strings.Join(qb.whereClauses, " AND "))
	}

	if len(qb.orderBy) > 0 {
		sqlBuilder.WriteString(" ORDER BY " + strings.Join(qb.orderBy, ", "))
	}

	if qb.limit > 0 {
		sqlBuilder.WriteString(" LIMIT ?")
		qb.whereArgs = append(qb.whereArgs, qb.limit)
	}

	return sqlBuilder.String(), qb.whereArgs
}
