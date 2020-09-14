package tools

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	// postgres and file drivers for go migrate
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)


var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate - connects to db *hardcoded localhost*, creates databases (core) for base-project and runs `migrate up` for each",
	Run: func(cmd *cobra.Command, args []string) {
		err := InitalizeDB()
		if err != nil {
			panic(err)
		}
	},
}

var dbNames = []string{
	"core",
}

// TODO: currently localhost dev only consider moving to env or allowing args?
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
)

func InitalizeDB() error {
	connection := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	for _, dbName := range dbNames {
		source := fmt.Sprintf("file://db/migrations/%s/", dbName)
		destination := fmt.Sprintf("postgres://postgres:postgres@localhost/%s?sslmode=disable", dbName)
		err := CreateAndInitialize(connection, dbName, source, destination)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateDBIfNotExists(connection, dbName string) error {
	// piggybacking off of 'postgres' driver registered via migrate import
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return err
	}
	defer db.Close()

	res := db.QueryRow(fmt.Sprintf(`SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('%s'));`, dbName))
	dbExists := false
	if err := res.Scan(&dbExists); err != nil {
		return err
	}

	if dbExists {
		log.Println(fmt.Sprintf("Database %s already exists, skipping creation", dbName))
		return nil
	}
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Created database %s", dbName))

	return nil
}

func CreateAndInitialize(connection, dbName, source, destination string) error {
	err := CreateDBIfNotExists(connection, dbName)
	if err != nil {
		return err
	}

	m, err := migrate.New(source, destination)
	if err != nil {
		return err
	}

	ver, _, err := m.Version()
	versionLess := false
	if err == migrate.ErrNilVersion {
		versionLess = true
	} else if err != nil {
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Println(fmt.Sprintf("Database %s has already been initialized to v.%v", dbName, ver))
		return nil
	} else if err != nil {
		return err
	}

	ver2, _, err := m.Version()
	if err != nil {
		return err
	}
	if versionLess {
		log.Println(fmt.Sprintf("empty database %s has been initialized to  v.%v", dbName, ver2))
	} else {
		log.Println(fmt.Sprintf("migrated database %s from v.%v to v.%v", dbName, ver, ver2))
	}
	return nil
}
