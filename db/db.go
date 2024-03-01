package db

import (
	"os"

	"github.com/abhishek-bhangalia-busy/detect-cycle/models"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var DB *pg.DB

func ConnectToDB() *pg.DB {
	DB = pg.Connect(&pg.Options{
		Addr:     os.Getenv("DB_ADDR"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	})
	return DB
}

func CreateSchema(db *pg.DB) error {
	models := []interface{}{
		(*models.EmployeeManager)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
