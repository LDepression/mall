/**
 * @Author: lenovo
 * @Description:
 * @File:  good_test.go
 * @Version: 1.0.0
 * @Date: 2023/03/03 22:39
 */

package query

import (
	"fmt"
	"log"
	"mall/internal/config"
	"mall/internal/dao"
	"os"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/stretchr/testify/assert"
)

//var DB *gorm.DB
//
func init() {
	m := &config.Mysql{
		User:     "root",
		Host:     "127.0.0.1",
		Port:     3306,
		Password: "zxz123456",
		DbName:   "mall",
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.Password, m.Host, m.Port, m.DbName)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢阙值
			Colorful:      true,        //禁用彩色
			LogLevel:      logger.Info,
		})
	//全局模式
	var err error
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})
	if err != nil {
		panic(err)
	}
	dao.Group.DB = DB
}

func TestCheckGoodByID(t *testing.T) {

	good := NewGood()
	assert.True(t, good.CheckGoodByID(421))
	//fmt.Println(good)
}
