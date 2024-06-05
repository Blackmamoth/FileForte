package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Blackmamoth/fileforte/config"
	"github.com/go-sql-driver/mysql"
)

var DB = getDBInstance()

func initDB() (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 config.MySQLConfig.MYSQL_USER,
		Passwd:               config.MySQLConfig.MYSQL_PASS,
		Addr:                 fmt.Sprintf("%s:%s", config.MySQLConfig.MYSQL_HOST, config.MySQLConfig.MYSQL_PORT),
		DBName:               config.MySQLConfig.MYSQL_DB_NAME,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	config.Logger.INFO("Application connected to MySQL server.")
	return db, nil
}

func getDBInstance() *sql.DB {
	db, err := initDB()
	if err != nil {
		config.Logger.ERROR(fmt.Sprintf("Application could connect to MySQL server: %v", err))
		os.Exit(1)
	}
	return db
}
