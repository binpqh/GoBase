package query_builder

import (
	"reflect"
	"strings"
)

type QueryBuilder[T any] struct {
	selectFields []string
	whereClauses []string
	whereArgs    []interface{}
	joins        []string
	limit        int
	orderBy      []string
}

func NewQueryBuilder[T any]() *QueryBuilder[T] {
	return &QueryBuilder[T]{}
}

func (qb *QueryBuilder[T]) Select(fields ...string) *QueryBuilder[T] {
	qb.selectFields = append(qb.selectFields, fields...)
	return qb
}

func (qb *QueryBuilder[T]) WhereEqual(field string, value interface{}) *QueryBuilder[T] {
	qb.whereClauses = append(qb.whereClauses, field+" = ?")
	qb.whereArgs = append(qb.whereArgs, value)
	return qb
}

func (qb *QueryBuilder[T]) OrderByASC(field string) *QueryBuilder[T] {
	qb.orderBy = append(qb.orderBy, field+" ASC")
	return qb
}

func (qb *QueryBuilder[T]) OrderByDESC(field string) *QueryBuilder[T] {
	qb.orderBy = append(qb.orderBy, field+" DESC")
	return qb
}

func (qb *QueryBuilder[T]) Join(table, onCondition string) *QueryBuilder[T] {
	qb.joins = append(qb.joins, "JOIN "+table+" ON "+onCondition)
	return qb
}

func (qb *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
	qb.limit = limit
	return qb
}

func (qb *QueryBuilder[T]) Build() (string, []interface{}) {
	var sql strings.Builder
	sql.WriteString("SELECT ")
	if len(qb.selectFields) > 0 {
		sql.WriteString(strings.Join(qb.selectFields, ", "))
	} else {
		sql.WriteString("*")
	}

	typeName := reflect.TypeOf((*T)(nil)).Elem().Name()
	sql.WriteString(" FROM " + strings.ToLower(typeName))

	if len(qb.joins) > 0 {
		sql.WriteString(" " + strings.Join(qb.joins, " "))
	}

	if len(qb.whereClauses) > 0 {
		sql.WriteString(" WHERE " + strings.Join(qb.whereClauses, " AND "))
	}

	if len(qb.orderBy) > 0 {
		sql.WriteString(" ORDER BY " + strings.Join(qb.orderBy, ", "))
	}

	if qb.limit > 0 {
		sql.WriteString(" LIMIT ?")
		qb.whereArgs = append(qb.whereArgs, qb.limit)
	}

	return sql.String(), qb.whereArgs
}
