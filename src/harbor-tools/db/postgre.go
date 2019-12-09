package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"harbor-tools/harbor-tools/config"
	"log"
	"sync"
)

var db *gorm.DB
var harborDB *gorm.DB
var once sync.Once

func DB() (*gorm.DB, *gorm.DB, error) {
	cfg := config.NewConfig()
	var err error
	once.Do(func() {
		db, err = gorm.Open("postgres", cfg.Postgres.DSN())
		if err != nil {
			log.Panicf("连接数据失败:", err)
		}
		err = db.DB().Ping()
		harborDB, err = gorm.Open("")
		harborDB, err = gorm.Open("postgres", cfg.HarborPostgres.DSN())
		if err != nil {
			log.Panicf("harbor DB 连接失败:", err)
		}
		err = harborDB.DB().Ping()

	})
	return db, harborDB, err
}
