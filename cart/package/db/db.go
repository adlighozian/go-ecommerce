package db

import (
	"api-gateway-go/helper/timeout"
	"database/sql"
	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDB struct {
	SqlDB *sql.DB
}

func NewGormDB(debug bool, driver, url string) (*GormDB, error) {
	if url == "" {
		return nil, errors.New("no database url")
	}

	gormDB := new(GormDB)
	err := gormDB.init(debug, driver, url)
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func (g *GormDB) init(debug bool, driver, url string) error {
	var gormLogger logger.Interface
	gormLogger = logger.Default.LogMode(logger.Silent)
	if debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	gormConf := new(gorm.Config)
	gormConf.Logger = gormLogger
	gormConf.PrepareStmt = true
	gormConf.SkipDefaultTransaction = true

	ctx, cancel := timeout.NewCtxTimeout()
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

	g.SqlDB = SqlDB
	return nil
}

func (g *GormDB) Close() error {
	return g.SqlDB.Close()
}
