package main

import (
	"github.com/gin-gonic/gin"
	"what-unexpected-summer/summer-two/controller"
	"what-unexpected-summer/summer-two/model"
	"what-unexpected-summer/summer-two/service"
)

func main() {
	model.InitDB()
	service.InitService()
	r := gin.Default()
	r.GET("/getGoods", controller.SelectGoods)
	r.POST("/addGood", controller.AddGood)
	r.POST("/order", controller.MakeOrder)

	_ = r.Run(":8080")
}

 

