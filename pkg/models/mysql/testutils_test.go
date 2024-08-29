package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	// Establish a sql.DB connection pool for our test database. Because our
	// setup and teardown scripts contains multiple SQL statements, we need
	// to use the `multiStatements=true` parameter in our DSN. This instructs
	// our MySQL database driver to support executing multiple SQL statements
	// in one db.Exec()` call.
	db, err := sql.Open("mysql", "root:root@/test_snippetbox?parseTime=true")
	if err != nil {
		t.Fatal(err)
	}
	// Read the setup SQL script from file and execute the statements.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	scriptsSlice := strings.Split(string(script), ";")
	newSlice := make([]string, len(scriptsSlice)-1)
	copy(newSlice, scriptsSlice[:len(scriptsSlice)-1])
	scriptsSlice = newSlice
	for _, s := range scriptsSlice {
		_, err = db.Exec(fmt.Sprintf("%s;", s))
		if err != nil {
			t.Fatal(err)
		}

	}
	// Return the connection pool and an anonymous function which reads and
	// executes the teardown script, and closes the connection pool. We can
	// assign this anonymous function and call it later once our test has
	// completed.
	return db, func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		scriptsSlice := strings.Split(string(script), ";")
		newSlice := make([]string, len(scriptsSlice)-1)
		copy(newSlice, scriptsSlice[:len(scriptsSlice)-1])
		scriptsSlice = newSlice
		for _, s := range scriptsSlice {
			_, err = db.Exec(fmt.Sprintf("%s;", s))
			if err != nil {
				t.Fatal(err)
			}
		}
		db.Close()
	}
}
