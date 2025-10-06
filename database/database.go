package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// DB represents a database connection
type DB struct {
	conn *sql.DB
}

// Model represents a base model with common fields
type Model struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// QueryBuilder provides a fluent query builder interface
type QueryBuilder struct {
	table      string
	wheres     []string
	whereArgs  []interface{}
	selects    []string
	orderBy    string
	limit      int
	offset     int
	db         *DB
}

// New creates a new database connection
func New(driver, dsn string) (*DB, error) {
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	
	return &DB{conn: conn}, nil
}

// Table creates a new query builder for the specified table
func (db *DB) Table(name string) *QueryBuilder {
	return &QueryBuilder{
		table:     name,
		selects:   []string{"*"},
		wheres:    []string{},
		whereArgs: []interface{}{},
		db:        db,
	}
}

// Select specifies which columns to select
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.selects = columns
	return qb
}

// Where adds a WHERE clause
func (qb *QueryBuilder) Where(column string, operator string, value interface{}) *QueryBuilder {
	qb.wheres = append(qb.wheres, fmt.Sprintf("%s %s ?", column, operator))
	qb.whereArgs = append(qb.whereArgs, value)
	return qb
}

// OrderBy adds an ORDER BY clause
func (qb *QueryBuilder) OrderBy(column string, direction string) *QueryBuilder {
	qb.orderBy = fmt.Sprintf("%s %s", column, direction)
	return qb
}

// Limit adds a LIMIT clause
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset adds an OFFSET clause
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// Get executes the query and returns all results
func (qb *QueryBuilder) Get() (*sql.Rows, error) {
	query := qb.buildQuery()
	return qb.db.conn.Query(query, qb.whereArgs...)
}

// First executes the query and returns the first result
func (qb *QueryBuilder) First() *sql.Row {
	qb.Limit(1)
	query := qb.buildQuery()
	return qb.db.conn.QueryRow(query, qb.whereArgs...)
}

// buildQuery builds the SQL query string
func (qb *QueryBuilder) buildQuery() string {
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(qb.selects, ", "), qb.table)
	
	if len(qb.wheres) > 0 {
		query += " WHERE " + strings.Join(qb.wheres, " AND ")
	}
	
	if qb.orderBy != "" {
		query += " ORDER BY " + qb.orderBy
	}
	
	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}
	
	if qb.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offset)
	}
	
	return query
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Connection returns the underlying sql.DB connection
func (db *DB) Connection() *sql.DB {
	return db.conn
}
