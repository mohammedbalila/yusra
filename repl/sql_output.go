package repl

import (
	"fmt"

	"github.com/mohammedbalila/yusra/storage"
	"github.com/mohammedbalila/yusra/utils"
)

func processSQLStmt(sql string) string {

	db, err := storage.GetDB()
	if err != nil {
		return err.Error()
	}

	rows, err := db.Query(sql)
	if err != nil {
		return fmt.Sprintf("Error reading input: %s", err)
	}

	columns, err := rows.Columns()
	if err != nil {
		return err.Error()
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	var results []map[string]interface{}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return err.Error()
		}
		// Create a map to hold the column data
		columnsMap := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			columnsMap[col] = v
		}
		results = append(results, columnsMap)
	}
	// print the data
	utils.PrintTable(columns, results)
	return ""
}
