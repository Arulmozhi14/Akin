package config

import (
	"database/sql"
	. "entity-batch-script/utils/constants"
	"fmt"

	_ "github.com/go-sql-driver/mysql" /*imported a driver for sql.Open*/
)

// CreateSQLdatabase initializes our db base on the type of database
func CreateSQLdatabase(dbDetails map[string]string) (db *sql.DB, err error) {
	switch dbDetails[Db] {
	case "mysql":
		connectionString := fmt.Sprintf(
			"%v:%v@tcp(%v:%v)/%v",
			dbDetails[DbUser],
			dbDetails[DbPass],
			dbDetails[DbHost],
			dbDetails[DbPort],
			dbDetails[DbName],
		)

		db, err = sql.Open(dbDetails[Db], connectionString)
	}

	if err != nil {
		err = fmt.Errorf("Error in openning db: %v", err)
	}

	return
}
