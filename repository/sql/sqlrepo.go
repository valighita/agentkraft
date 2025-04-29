package sql

import (
	"errors"
	"log"
	"os"
	"time"
	"valighita/agentkraft/repository"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	sqliteDbFile  = "./agentkraft.sqlite"
	defaultDbName = "agentkraft"
)

func initDb() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	mysqlHost := os.Getenv("MYSQL_HOST")

	if mysqlHost != "" {
		log.Println("Initializing MySQL database")

		var tries = 0
		for {
			tries++
			mysqlUser := os.Getenv("MYSQL_USER")
			mysqlPassword := os.Getenv("MYSQL_PASSWORD")
			dbName := os.Getenv("MYSQL_DB")
			if dbName == "" {
				dbName = defaultDbName
			}
			dsn := mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlHost + ":3306)/" + dbName + "?parseTime=true"
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				if tries > 5 {
					return nil, err
				}

				log.Println("Error connecting to MySQL, retrying in 200ms")
				time.Sleep(200 * time.Millisecond)
			} else {
				break
			}
		}
	} else {
		log.Println("Initializing sqlite database")
		db, err = gorm.Open(sqlite.Open(sqliteDbFile), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	}

	db.AutoMigrate(&repository.Agent{})

	return db, nil
}

func GetSqlRepositories() (repository.AgentsRepository, error) {
	db, err := initDb()
	if err != nil {
		log.Println("Error initializing database", err)
		return nil, err
	}
	agentsRepository := NewSQLAgentsRepository(db)
	if agentsRepository == nil {
		return nil, errors.New("error initializing agent repository")
	}

	return agentsRepository, nil
}
