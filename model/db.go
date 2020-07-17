package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var DB *gorm.DB

func InitDB()  {
	var err error
	db, err := gorm.Open("mysql","mysql","root:@tcp(127.0.0.1:3306)/dome7?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Panicf("Panic while connecting the gorm. Error: %s",err)
	}

	DB = db
	DB.SingularTable(true)
	if !DB.HasTable(&Goods{}) {
		if err := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(&Goods{}).Error; err != nil {
			panic(err)
		}
	}

	if !DB.HasTable(&Order{}) {
		if err := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(&Order{}).Error; err != nil {
			panic(err)
		}
	}

}

 

