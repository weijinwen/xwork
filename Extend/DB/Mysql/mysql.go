package Mysql

import (
	"xwork/Extend/DB"
	"xwork/Extend/Gophp"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Cache Memcache adapter.
type Db struct {
	conn      *gorm.DB
	server    string
	conninfo  string
	maxIdle   int
	maxLife   int64
	maxActive int
	showSql   bool
}

func NewMysqlDb() DB.DbInterface {
	return &Db{}
}

// 获取数据库连接
func (rc *Db) GetOrm() *gorm.DB {
	return rc.conn
}

// 初始化数据库
func (rc *Db) StartAndGC(server string) error {
	rc.server = server
	rc.conninfo = viper.GetString(fmt.Sprintf("db.%s.url", server))
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return nil
}

// 连接数据库
func (rc *Db) connectInit() error {
	var err error
	db, err := gorm.Open("mysql", rc.conninfo)
	if err != nil {
		logrus.Fatalf("opens database failed: %s", err.Error())
		return err
	}
	db.SingularTable(false)
	db.DB().SetMaxIdleConns(rc.getMaxIdle())
	db.DB().SetMaxOpenConns(rc.getMaxActive())
	db.DB().SetConnMaxLifetime(time.Duration(rc.getMaxLife()))
	db.LogMode(rc.getShowSql())
	rc.conn = db
	return nil
}

//最大连接数
func (rc *Db) getMaxIdle() int {
	maxIdle := viper.GetInt(fmt.Sprintf("db.%s.maxIdle", rc.server))
	return Gophp.Ternary(maxIdle <= 0, 10, maxIdle).(int)
}

//最大空闲连接数
func (rc *Db) getMaxActive() int {
	maxActive := viper.GetInt(fmt.Sprintf("db.%s.maxActive", rc.server))
	return Gophp.Ternary(maxActive <= 0, 20, maxActive).(int)
}

//连接存活时间
func (rc *Db) getMaxLife() int64 {
	maxLife := viper.GetInt64(fmt.Sprintf("db.%s.maxLife", rc.server))
	return int64(Gophp.Ternary(maxLife <= 0, 0, maxLife).(int))
}

//是否显示sql
func (rc *Db) getShowSql() bool {
	return viper.GetBool(fmt.Sprintf("db.%s.showSql", rc.server))
}

func init() {
	DB.Register("mysql", NewMysqlDb)
}
