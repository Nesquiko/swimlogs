package integration

import (
	"database/sql"
	"log"
	"reflect"

	"github.com/Nesquiko/swimlogs/pkg/data"
)

var (
	CleanSession = "truncate table session"
)

// cleanDB clears all tables in database.
func cleanDB(db data.DBConn) {
	field := reflect.ValueOf(db).Elem().FieldByName("DB")

	if !field.IsValid() {
		log.Fatal("couldn't get the field")
	}

	sqlDB, ok := field.Interface().(*sql.DB)
	if !ok {
		log.Fatal(ok)
	}

	_, err := sqlDB.Exec(CleanSession)
	if err != nil {
		log.Fatal(err)
	}
}
