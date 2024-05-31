package storage

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mohammedbalila/yusra/internal"
	"github.com/mohammedbalila/yusra/utils"
)

func generateTableCreateStmt(cols []internal.Column, filename string) (string, string) {
	tableName := utils.FilenameToTableName(filename)

	colsToSQL := utils.Map(cols, func(item internal.Column) string {

		return fmt.Sprintf("%s %s", item.Name, convertGoTypeToSQLite(item.DataType))
	})
	createTableSQL := fmt.Sprintf("create table %s(%s)", tableName, strings.Join(colsToSQL, ",\n"))
	return tableName, createTableSQL
}

func convertGoTypeToSQLite(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "INTEGER"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.Array, reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return "BLOB" // Byte slices (e.g., []byte) should map to BLOB
		}
		return "TEXT" // Other slices/arrays can be stored as TEXT (JSON, CSV, etc.)
	case reflect.String, reflect.Struct, reflect.Map:
		return "TEXT" // For now Strings, Structs, and maps will be stored as TEXT
	default:
		return "TEXT" // Default to TEXT for unsupported types
	}
}

func interfaceToString(value interface{}) interface{} {
	if value == nil {
		return "null"
	}

	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return value
	case reflect.Array, reflect.Slice, reflect.String, reflect.Struct, reflect.Map:
		return fmt.Sprintf("%s", value)
	default:
		return fmt.Sprintf("%s", value)
	}
}
