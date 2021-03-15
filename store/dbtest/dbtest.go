package dbtest

import (
	"os"

	"github.com/fox-one/pkg/store/db"
)

func Connect() (*db.DB, error) {
	var (
		dialect = "sqlite3"
		uri     = ":memory:"
	)

	if os.Getenv("PANDO_DATABASE_DIALECT") != "" {
		dialect = os.Getenv("PANDO_DATABASE_DIALECT")
		uri = os.Getenv("PANDO_DATABASE_DATASOURCE")

	}

	c, err := db.Connect(dialect, uri)
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(c); err != nil {
		return nil, err
	}

	return c, nil
}

func Disconnect(db *db.DB) error {
	return db.Close()
}
