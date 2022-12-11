package statistic

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	CONFIG          *Config
	DB              *gorm.DB
	TrackerInfoChan chan *TrackerInfo
)

func InitViper(path string) {
	vi := viper.New()
	vi.SetConfigName("config")
	vi.SetConfigType("yaml")
	// vi.AddConfigPath("conf")
	vi.AddConfigPath(path)
	vi.SetDefault("outpath", defaultPath())
	vi.SetDefault("gap", 5)
	if err := vi.ReadInConfig(); err != nil {
		log.Println(err.Error())
	}
	if err := vi.Unmarshal(&CONFIG); err != nil {
		log.Println(err.Error())
		log.Fatal("解析配置文件出错！")
	}
	// log.Println(CONFIG)
	// global.VIPER = vi
}

func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		CONFIG.Mysql.Username,
		CONFIG.Mysql.Password,
		CONFIG.Mysql.Host,
		CONFIG.Mysql.Port,
		CONFIG.Mysql.DBname,
		CONFIG.Mysql.Charset,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("数据库连接错误！")
	}
	DB = db
}
