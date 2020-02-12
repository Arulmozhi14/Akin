package identity

import (
	"database/sql"
	"entity-batch-script/schema"
	"fmt"
	"strings"
)

// Model will have a map call db
type Model struct {
	Db        *sql.DB
	DbName    string
	TableName string
}

// TableExist is used for checking if the entity table exist [should be used for this one]
func (model *Model) TableExist() (bool, error) {
	query := fmt.Sprintf(
		"SELECT * FROM information_schema.tables WHERE table_schema='%v' AND table_name='%v' limit 1",
		model.DbName,
		model.TableName,
	)
	rows, err := model.Db.Query(query)

	if err != nil {
		return false, fmt.Errorf("Error in checking database table: %v", err)
	}

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

// CreateTable creates a new table given a column name
func (model *Model) CreateTable(columnNames []string) error {
	var columnTypes []string

	for _, columnName := range columnNames {
		columnType := fmt.Sprintf("%v VARCHAR(64)", columnName)
		columnTypes = append(columnTypes, columnType)
	}

	query := fmt.Sprintf(
		"CREATE TABLE %v(%v)",
		model.TableName,
		strings.Join(columnTypes, ", "),
	)

	// "CREATE TABLE tableName (columnName VARCHAR(64), ...)"
	// varchar64 for a 64 length hash
	_, err := model.Db.Exec(query)

	if err != nil {
		return fmt.Errorf("Error in creating table: %v", err)
	}

	return nil
}

// DropTable drops a table
func (model *Model) DropTable() error {
	query := fmt.Sprintf(
		"DROP TABLE %v",
		model.TableName,
	)

	_, err := model.Db.Exec(query)

	if err != nil {
		return fmt.Errorf("Error in dropping table: %v", err)
	}

	return nil
}

func createValuesTemplate(columns []schema.TopicData) (string, string) {
	var wildcards []string
	var fields []string

	for _, col := range columns {
		wildcards = append(wildcards, "?")
		fields = append(fields, col.Topic)
	}

	return fmt.Sprintf("(%v)", strings.Join(fields, ", ")), fmt.Sprintf("(%v)", strings.Join(wildcards, ", "))
}

// BulkInsert inserts bulk given a two-dimentional array of TopicData
func (model *Model) BulkInsert(topicDataMatrix [][]schema.TopicData) error {
	firstElement := topicDataMatrix[0]
	valueStrings := make([]string, 0, len(topicDataMatrix))
	valueArgs := make([]interface{}, 0, len(topicDataMatrix)*len(firstElement))
	valueTemplate, valueWildcards := createValuesTemplate(firstElement)

	for _, row := range topicDataMatrix {
		valueStrings = append(valueStrings, valueWildcards)
		for _, col := range row {
			valueArgs = append(valueArgs, col.Raw)
		}
	}

	stmt := fmt.Sprintf("INSERT INTO %v %v VALUES %v", model.TableName, valueTemplate, strings.Join(valueStrings, ", "))
	_, err := model.Db.Exec(stmt, valueArgs...)

	if err != nil {
		return fmt.Errorf("Error in inserting: %v", err)
	}

	return nil
}
