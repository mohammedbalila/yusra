package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/olekukonko/tablewriter"
)

func FilenameToTableName(filename string) string {
	// Remove the file extension
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		filename = filename[:idx]
	}

	filename = strings.ToLower(filename)
	re := regexp.MustCompile(`[^a-z0-9_]+`)
	tableName := re.ReplaceAllString(filename, "_")
	// Ensure the name starts with a letter or underscore
	if len(tableName) > 0 && !unicode.IsLetter(rune(tableName[0])) && tableName[0] != '_' {
		tableName = "_" + tableName
	}

	return tableName
}

// generic map function
func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PrintTable(header []string, rows []map[string]interface{}) {
	table := GetTableWriter(header)
	for _, row := range rows {
		var rowData []string
		for _, col := range header {
			rowData = append(rowData, fmt.Sprintf("%v", row[col]))
		}
		table.Append(rowData)
	}

	table.Render()
}

func GetTableWriter(header []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	// Set table options
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetRowSeparator("-")
	table.SetColumnSeparator("|")
	table.SetCenterSeparator("+")
	table.SetHeaderLine(true)

	return table
}
