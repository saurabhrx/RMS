package database

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
)

var RMS *sqlx.DB

func ConnectToDB(host, port, user, password, dbname string) error {
	plsqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	DB, err := sqlx.Open("postgres", plsqlInfo)
	if err != nil {
		return err
	}
	err = DB.Ping()
	if err != nil {
		return err
	}
	RMS = DB
	return migrateUp(RMS)

}

func migrateUp(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://database/migrations", "postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	fmt.Println("migration completed")
	return nil
}

func CloseDBConnection() error {
	return RMS.Close()
}

func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := RMS.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()
	err = fn(tx)
	return err
}

func SetUpBindVars(stmt, bindVar string, length int) string {
	bindVar += ","
	stmt = fmt.Sprintf(stmt, strings.Repeat(bindVar, length))
	return replaceSQL(strings.TrimSuffix(stmt, ","), "?")
}
func replaceSQL(stmt, pattern string) string {
	count := strings.Count(stmt, "?")
	for i := 1; i <= count; i++ {
		stmt = strings.Replace(stmt, pattern, "$"+strconv.Itoa(i), 1)
	}
	return stmt
}
