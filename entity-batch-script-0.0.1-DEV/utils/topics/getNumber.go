package topics

import "fmt"

// GetNumber gets the corresponding number of the column name on the map
func GetNumber(topicMap map[string]interface{}, columnName string) string {
	var topicNumber string

	for key, value := range topicMap {
		if columnName == value.(string) {
			topicNumber = string(key)
			break
		}
	}

	return fmt.Sprintf("topic%v", topicNumber)
}
