package main

import (
	"fmt"
	"github.com/gertd/go-pluralize"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"strings"
)

type BasicCrud struct {
	tableName            string
	methodName           string
	idField              string
	fields               []string
	defaultOrderingField string
	createQuery          string
	getByIdQuery         string
	getAllQuery          string
	getAllPaginated      string
	updateByIdQuery      string
	deleteByIdQuery      string
}

func (bc *BasicCrud) Construct(tableName, modelFilePath string) {
	bc.tableName = tableName
	bc.methodName = getSingularCamelCasedName(tableName)
	bc.fields = parseModel(bc.methodName, modelFilePath)
	bc.idField = bc.fields[0]
	if slices.Contains(bc.fields, "created_at") {
		bc.defaultOrderingField = "created_at"
	} else {
		bc.defaultOrderingField = bc.idField
	}
	bc.makeCreateQuery()
	bc.makeUpdateByIdQuery()
	bc.makeDeleteByIdQuery()
	bc.makeGetByIdQuery()
	bc.makeGetAllQuery()
	bc.makeGetAllPaginatedQuery()
	fmt.Println(bc.createQuery)
	fmt.Println(bc.updateByIdQuery)
	fmt.Println(bc.deleteByIdQuery)
	fmt.Println(bc.getAllQuery)
	fmt.Println(bc.getAllPaginated)
	fmt.Println(bc.getByIdQuery)
}

func (bc *BasicCrud) makeCreateQuery() {
	var fields []string
	var values []string
	placement := 1
	for _, column := range bc.fields {
		if column != bc.idField && column != bc.defaultOrderingField && column != "updated_at" {
			fields = append(fields, column)
			values = append(values, fmt.Sprintf("$%d", placement))
			placement++
		} else if column == "updated_at" {
			fields = append(fields, column)
			values = append(values, "now()")
		}
	}
	fieldsStr := strings.Join(fields, ", ")
	valuesStr := strings.Join(values, ", ")
	bc.createQuery = fmt.Sprintf("-- name: Create%s :one\nINSERT INTO %s (%s) VALUES (%s) RETURNING *;",
		bc.methodName, bc.tableName, fieldsStr, valuesStr)
}

func (bc *BasicCrud) makeUpdateByIdQuery() {
	var setStatements []string
	for _, column := range bc.fields {
		if column != bc.idField && column != bc.defaultOrderingField && column != "updated_at" {
			setStatements = append(setStatements, fmt.Sprintf("%s = coalesce(sqlc.narg('%s'), %s)",
				column, column, column))
		} else if column == "updated_at" {
			setStatements = append(setStatements, fmt.Sprintf("%s = now()", column))
		}
	}
	setStatementsStr := strings.Join(setStatements, ", ")
	bc.updateByIdQuery = fmt.Sprintf("-- name: Update%sById :one\nUPDATE %s SET %s WHERE %s = $1 RETURNING *;",
		bc.methodName, bc.tableName, setStatementsStr, bc.idField)
}

func (bc *BasicCrud) makeDeleteByIdQuery() {
	bc.deleteByIdQuery = fmt.Sprintf("-- name: Delete%sById :execrows\nDELETE FROM %s WHERE %s = $1;",
		bc.methodName, bc.tableName, bc.idField)
}

func (bc *BasicCrud) makeGetAllQuery() {
	pluralClient := pluralize.NewClient()
	plural := pluralClient.Plural(bc.methodName)
	bc.getAllQuery = fmt.Sprintf("-- name: GetAll%s :many\nSELECT * FROM %s ORDER BY %s;",
		plural, bc.tableName, bc.defaultOrderingField)
}

func (bc *BasicCrud) makeGetAllPaginatedQuery() {
	pluralClient := pluralize.NewClient()
	plural := pluralClient.Plural(bc.methodName)
	bc.getAllPaginated = fmt.Sprintf("-- name: GetAll%s :many\nSELECT * FROM %s ORDER BY %s OFFSET $1 LIMIT $2;",
		plural, bc.tableName, bc.defaultOrderingField)
}

func (bc *BasicCrud) makeGetByIdQuery() {
	bc.getByIdQuery = fmt.Sprintf("-- name: Get%sById :one\nSELECT * FROM %s WHERE %s = $1 LIMIT 1;",
		bc.methodName, bc.tableName, bc.idField)
}

func (bc *BasicCrud) CreateQueriesFile(path string) {
	var queries []string
	queries = append(queries, bc.createQuery)
	queries = append(queries, bc.getByIdQuery)
	queries = append(queries, bc.getAllQuery)
	queries = append(queries, bc.getAllPaginated)
	queries = append(queries, bc.updateByIdQuery)
	queries = append(queries, bc.deleteByIdQuery)
	queriesString := strings.Join(queries, "\n\n")
	os.MkdirAll(path, os.ModeDir|0755)
	filePath := fmt.Sprintf("./%s/%s.sql", path, bc.tableName)
	err := os.WriteFile(filePath, []byte(queriesString), 0644)
	if err != nil {
		log.Panicf("file creation failed with error: %v", err)
		return
	}
}
