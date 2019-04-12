package db

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("vsummary")

// Backend interface.
type Backend struct {
	db *sqlx.DB
}

// new Backend
func InitBackend() (b *Backend, err error) {

	driver := viper.GetString("backend.db_driver")
	dsn := viper.GetString("backend.db_dsn")

	db, err := sqlx.Connect(driver, dsn)

	if err != nil {
		log.Error("failed to connect to database")
		err = fmt.Errorf("database error: %s", err)
	} else {
		b = &Backend{db: db}
		log.Info("connection to database successful")
	}

	return
}

func NewBackend() *Backend {
	return &Backend{}
}

// Test DB connection
func (b *Backend) checkDB() (err error) {
	if b.db == nil {
		return errors.New("Database connection is not set up")
	}
	err = b.db.Ping()
	return
}

// SetDB changes the underlying sql.DB object Accessor is manipulating.
func (b *Backend) SetDB(db *sqlx.DB) {
	b.db = db
	return
}

// return DB instance for datatables (for now)
func (b *Backend) GetDB() (db *sqlx.DB) {
	db = b.db
	return
}

// Apply database schemas
func (b *Backend) ApplySchemas() (err error) {

	// check if db connection is available
	err = b.checkDB()
	if err != nil {
		return
	}

	// apply all schemas
	for _, schema := range generateSqlSchemas() {

		log.Debugf("Applying schema: %s", schema.Name)
		_, err = b.db.MustExec(schema.SqlQuery).RowsAffected()
		if err == nil {
			log.Debugf("db schema is present: %s", schema.Name)
		} else {
			break
		}

	}

	return
}
