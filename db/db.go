package db

import (
	"fmt"
	"math"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/xerrors"
)

func New(host string, port uint16, dbName string, user string, password string) (*sqlx.DB, error) {
	mysqlCfg := mysql.NewConfig()

	mysqlCfg.User = user
	mysqlCfg.Passwd = password
	mysqlCfg.Net = "tcp"
	mysqlCfg.Addr = fmt.Sprintf("%s:%d", host, port)
	mysqlCfg.DBName = dbName
	mysqlCfg.ParseTime = true
	mysqlCfg.ClientFoundRows = true
	mysqlCfg.Params = map[string]string{
		"charset": "utf8mb4",
	}
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, xerrors.Errorf("failed to get location: %w", err)
	}
	mysqlCfg.Loc = loc

	db, err := sqlx.Open("mysql", mysqlCfg.FormatDSN())
	if err != nil {
		return nil, xerrors.Errorf("failed to connect to database: %w", err)
	}

	// Check connection
	const retryMax = 15
	for i := 0; i < retryMax; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		if i == retryMax-1 {
			return nil, xerrors.Errorf("failed to connect to database: %w", err)
		}
		duration := time.Millisecond * time.Duration(math.Pow(1.5, float64(i))*1000)

		time.Sleep(duration)
	}
	return db, nil
}
