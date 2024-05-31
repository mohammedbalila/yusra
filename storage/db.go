package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mohammedbalila/yusra/internal"
	"github.com/mohammedbalila/yusra/utils"
)

var db *sql.DB
var DATASETS_TABLE_NAME = "yusra_datasets"

func GetDB() (*sql.DB, error) {
	var err error
	if db == nil {
		db, err = sql.Open("sqlite3", "./yusra.db")
		if err != nil {
			return nil, err
		}
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func createDBTable(cols []internal.Column, filename string) error {
	var err error
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (dataset_name text);", DATASETS_TABLE_NAME))
	if err != nil {
		return fmt.Errorf("couldn't setup database: %s", err)
	}

	tableName, createTable := generateTableCreateStmt(cols, filename)
	// drop if table exists
	_, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
	if err != nil {
		return err
	}
	_, err = db.Exec(createTable)
	if err != nil {
		return fmt.Errorf("couldn't create data table: %s", err)
	}
	err = logDatasetName(tableName)
	if err != nil {
		return err
	}
	return nil
}

func logDatasetName(tableName string) error {
	findQuery := fmt.Sprintf("select dataset_name from %s where dataset_name = (?)", DATASETS_TABLE_NAME)
	findStmt, err := db.Prepare(findQuery)
	if err != nil {
		return fmt.Errorf("couldn't update tables history: %s", err)
	}
	defer findStmt.Close()
	rows, err := findStmt.Query(tableName)
	if err != nil {
		return fmt.Errorf("couldn't update tables history: %s", err)
	}

	defer rows.Close()
	// the dataset was already logged, this happens in case of reloading a file
	if rows.Next() {
		return nil
	}

	insertQuery := fmt.Sprintf("insert into %s values (?)", DATASETS_TABLE_NAME)
	insertStmt, err := db.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("couldn't update tables history: %s", err)
	}

	defer insertStmt.Close()
	_, err = insertStmt.Exec(tableName)
	if err != nil {
		return fmt.Errorf("couldn't update tables history: %s", err)
	}
	return nil
}

func GetLoadedDatasets() error {
	query := fmt.Sprintf("select dataset_name from %s", DATASETS_TABLE_NAME)
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("couldn't update tables history: %s", err)
	}
	var datasetName string
	table := utils.GetTableWriter([]string{"Name"})
	for rows.Next() {
		rows.Scan(&datasetName)
		table.Append([]string{datasetName})
	}
	table.Render()
	return nil
}

func populateSQLTable(tableName string, records []map[string]interface{}, cols []internal.Column) error {
	datasetLength := len(records)
	columnNames := utils.Map(cols, func(col internal.Column) string { return col.Name })
	argumentsLength := len(columnNames)
	// sorting the column names will help us in writing it later
	sort.Strings(columnNames)

	values := make([]string, datasetLength)
	args := make([]any, datasetLength*argumentsLength)
	// start := 0
	pos := 0
	for i := 0; i < len(records); i++ {
		valueArgs := utils.Map(columnNames, func(_ string) string { return "?" })
		// build the values stmt (?, ?, ?, ...)
		values[i] = fmt.Sprintf("(%s)", strings.Join(valueArgs, ","))
		for j := 0; j < argumentsLength; j++ {
			args[pos+j] = interfaceToString(records[i][columnNames[j]])
		}
		pos += argumentsLength
	}

	query := fmt.Sprintf("insert into %s (%s) values %s", tableName, strings.Join(columnNames, ","), strings.Join(values, ","))
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}
	return nil
}

func LoadNewJsonFile(filename string) error {
	ok, err := utils.FileExists(filename)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("file does not exist, please check the path and try again")
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var results []map[string]interface{}
	err = json.Unmarshal(file, &results)
	if err != nil {
		return err
	}

	var cols []internal.Column
	entry := sampleEntry(&results)
	for key := range entry {
		col := internal.Column{DataType: reflect.TypeOf(entry[key]), Name: key}
		cols = append(cols, col)
	}

	tableName := utils.FilenameToTableName(filename)
	err = createDBTable(cols, filename)
	if err != nil {
		return err
	}

	err = populateSQLTable(tableName, results, cols)
	if err != nil {
		return err
	}
	return err
}

func sampleEntry(records *[]map[string]interface{}) map[string]interface{} {
	_records := *records
	entry := _records[0]
	for k, v := range entry {
		// TODO: fix later because this is not efficient!
		if v == nil {
			// in case all values were null, use a string as a placeholder
			entry[k] = ""
			for _, r := range _records {
				if r[k] != nil {
					entry[k] = r[k]
					break
				}
			}

		}
	}
	return entry
}

func GetTableStats(tableName string) error {
	rows, err := db.Query(fmt.Sprintf("select name AS column_name, type AS data_type from pragma_table_info('%s')", tableName))
	if err != nil {
		return fmt.Errorf("couldn't read tables stats: %s", err)
	}

	defer rows.Close()
	type TableStatRecord struct {
		ColumnName string
		DataType   string
	}

	var records []TableStatRecord
	var columnName string
	var dataType string

	for rows.Next() {
		rows.Scan(&columnName, &dataType)
		tsr := TableStatRecord{ColumnName: columnName, DataType: dataType}
		records = append(records, tsr)
	}

	table := utils.GetTableWriter([]string{"Column Name", "Data Type", "Non-Null Count", "Total Count"})
	totalCount := fmt.Sprintf("%d", getColumnCount(tableName))
	for _, record := range records {
		nonNullCount := fmt.Sprintf("%d", getColumnNonNullCount(tableName, record.ColumnName))
		table.Append([]string{record.ColumnName, record.DataType, nonNullCount, totalCount})
	}
	table.Render()
	return nil
}

func getColumnNonNullCount(tableName string, columnName string) int {
	var count int
	rows, err := db.Query(fmt.Sprintf("SELECT COUNT(%s) as count FROM %s WHERE %s IS NOT NULL", columnName, tableName, columnName))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("couldn't read tables stats: %s", err))
		return 0
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&count)
	}
	return count
}

func getColumnCount(tableName string) int {
	var count int
	rows, err := db.Query(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("couldn't read tables stats: %s", err))
		return 0
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&count)
	}
	return count
}
