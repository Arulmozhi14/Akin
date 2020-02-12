package entity

import (
	"database/sql"
	"entity-batch-script/schema"
	"entity-batch-script/utils/keccak256"
	"entity-batch-script/utils/topics"
	"fmt"
	"strings"

	"github.com/thoas/go-funk"
)

// Model will have a map call db
type Model struct {
	Db        *sql.DB
	TopicMap  map[string]interface{}
	TableName string
}

// GetRows return an array of []TopicData, error
func (model *Model) GetRows() ([][]schema.TopicData, error) {
	rows, queryErr := model.Db.Query(fmt.Sprintf("SELECT * FROM %v", model.TableName))

	if queryErr != nil {
		return nil, fmt.Errorf("Error in checking database table: %v", queryErr)
	}

	columns, columnsErr := rows.Columns()

	if columnsErr != nil {
		return nil, fmt.Errorf("Error in getting the columns: %v", columnsErr)
	}

	colCount := len(columns)
	values := make([]interface{}, colCount) // initialize array of interface with len colCount
	valuePtrs := make([]interface{}, colCount)
	TopicData := []schema.TopicData{} // store all TopicData here, can't make a dynamic two-dimentional array
	var rowCount int

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i] // populate valuePtrs with value pointer
		}

		scanErr := rows.Scan(valuePtrs...) // accepts array of pointers

		if scanErr != nil {
			return nil, fmt.Errorf("Error in rows scan: %v", scanErr)
		}

		for i, col := range columns {
			var data interface{} // interface{} - code that handles values of unknown type

			val := values[i]
			b, ok := val.([]byte) //get the underlying []byte value from val

			newColumnName := topics.GetNumber(model.TopicMap, col)

			// may not be usefull because we assume that every data is a string and to be hashed
			if ok { // returns true if b is []byte
				allCaps := strings.ToUpper(string(b))

				if funk.Contains(topics.GetTopicNames, newColumnName) {
					data = strings.ReplaceAll(allCaps, ".", "")
				} else {
					data = allCaps
				}
			} else {
				data = val
			}

			TopicData = append(
				TopicData,
				schema.TopicData{
					Topic: newColumnName,
					Raw:   keccak256.Hash(fmt.Sprintf("%v", data)),
				},
			)
		}

		rowCount++
	}

	TopicDataMatrix := make([][]schema.TopicData, rowCount)

	// make one-dimentional array into two-dimentional array
	for i := range TopicDataMatrix {
		TopicDataMatrix[i] = TopicData[i*colCount : (i+1)*colCount]
	}

	return TopicDataMatrix, nil
}
