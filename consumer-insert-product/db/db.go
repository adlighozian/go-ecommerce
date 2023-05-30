package db

import (
	"consumer-product-go/helpers"
	"database/sql"
	"errors"
	"log"

	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDB struct {
	SQLDB *sql.DB
}

func NewGormDB(debug bool, driver string, url string) (*GormDB, error) {
	if url == "" {
		return nil, errors.New("no database url")
	}

	gormDB := new(GormDB)
	err := gormDB.init(debug, driver, url)
	if err != nil {
		log.Println("error init")
		return nil, err
	}

	return gormDB, nil
}

func (g *GormDB) init(debug bool, driver string, url string) error {
	var gormLogger logger.Interface
	gormLogger = logger.Default.LogMode(logger.Silent)
	if debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	gormConf := new(gorm.Config)
	gormConf.Logger = gormLogger
	gormConf.PrepareStmt = true
	gormConf.SkipDefaultTransaction = true

	ctx, cancel := helpers.NewCtxTimeout()
	defer cancel()

	var dialector gorm.Dialector
	if driver == "postgres" {
		dialector = postgres.Open(url)
	}

	db, err := gorm.Open(dialector, gormConf)
	if err != nil {
		return err
	}

	SqlDB, err := db.DB()
	if err != nil {
		return err
	}

	SqlDB.SetMaxIdleConns(10)
	SqlDB.SetMaxOpenConns(100)
	SqlDB.SetConnMaxIdleTime(5 * time.Minute)
	SqlDB.SetConnMaxLifetime(60 * time.Minute)

	if err = SqlDB.PingContext(ctx); err != nil {
		return err
	}

	g.SQLDB = SqlDB
	return nil
}

func (g *GormDB) Close() error {
	return g.SQLDB.Close()
}
