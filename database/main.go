package database

import (
	"database/sql"
	"errors"
	"os"
	"strconv"

	"github.com/Improwised/golang-api/config"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql" // import mysql if it is used
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbURL string
var err error

const (
	postgres = "postgres"
	mysql    = "mysql"
	sqlite3  = "sqlite3"
)

// Connect with database
func Connect(cfg config.DBConfig) (*goqu.Database, error) {
	switch cfg.Dialect {
	case postgres:
		return postgresDBConnection(cfg)
	case mysql:
		return mysqlDBConnection(cfg)
	case sqlite3:
		return sqlite3DBConnection(cfg)
	default:
		return nil, errors.New("no suitable dialect found")
	}
}

func sqlite3DBConnection(cfg config.DBConfig) (*goqu.Database, error) {

	if _, err = os.Stat(cfg.SQLiteFilePath); err != nil {
		file, err := os.Create(cfg.SQLiteFilePath)
		if err != nil {
			panic(err)
		}
		err = file.Close()
		if err != nil {
			return nil, err
		}
	}
	db, err = sql.Open(sqlite3, "./"+cfg.SQLiteFilePath)
	if err != nil {
		return nil, err
	}
	return goqu.New(sqlite3, db), err
}

func mysqlDBConnection(cfg config.DBConfig) (*goqu.Database, error) {
	dbURL = cfg.Username + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + strconv.Itoa(cfg.Port) + ")/" + cfg.Db
	if db == nil {
		db, err = sql.Open(mysql, dbURL)
		if err != nil {
			return nil, err
		}
		return goqu.New(mysql, db), err
	}
	return goqu.New(mysql, db), err
}

func postgresDBConnection(cfg config.DBConfig) (*goqu.Database, error) {
	dbURL = "postgres://" + cfg.Username + ":" + cfg.Password + "@" + cfg.Host + ":" + strconv.Itoa(cfg.Port) + "/" + cfg.Db + "?" + cfg.QueryString
	if db == nil {
		db, err = sql.Open(postgres, dbURL)
		if err != nil {
			return nil, err
		}
		return goqu.New(postgres, db), err
	}
	return goqu.New(postgres, db), err
}
