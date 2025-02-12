package pool

import (
	"database/sql"
	"math/rand"
	"sync"
	"time"

	"github.com/wswz/go_commons/log"
)

var (
	connNum = 1
	lock    sync.Mutex
)

type Connection struct {
	*sql.DB
	Dsn string
}

func (c *Connection) ReConnect() error {
	sqlDB, err := sql.Open("clickhouse", c.Dsn)
	if err != nil {
		log.Info("reconnect to ", c.Dsn, err.Error())
		return err
	}
	log.Info("reconnect success to  ", c.Dsn)
	c.DB = sqlDB
	return nil
}

var poolMaps = map[string][]*Connection{}

func SetDsn(name string, dsn string, maxLifeTime time.Duration) {
	lock.Lock()
	defer lock.Unlock()

	sqlDB, err := sql.Open("clickhouse", dsn)
	if err != nil {
		panic(err)
	}
	sqlDB.SetConnMaxLifetime(maxLifeTime * time.Second)
	setConnectionParams(sqlDB)
	if ps, ok := poolMaps[name]; ok {
		//达到最大限制了，不需要新建conn
		if len(ps) >= connNum {
			return
		}
		ps = append(ps, &Connection{sqlDB, dsn})
		poolMaps[name] = ps
	} else {
		poolMaps[name] = []*Connection{&Connection{sqlDB, dsn}}
	}
	log.Info(">>>. Set clickhouse dsn success.")
}

func GetConn(name string) *Connection {
	lock.Lock()
	defer lock.Unlock()

	ps := poolMaps[name]
	return ps[rand.Intn(len(ps))]
}

func setConnectionParams(db *sql.DB) {
	db.SetMaxOpenConns(36)
	db.SetMaxIdleConns(12)
}

func CloseAll() {
	for _, ps := range poolMaps {
		for _, c := range ps {
			c.Close()
		}
	}
}