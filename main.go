package main

import (
	"github.com/gin-gonic/gin"
	"summer-two/controller"
	"summer-two/model"
	"summer-two/service"
)

func main() {
	model.InitDB()
	service.InitService()
	r := gin.Default()
	r.GET("/getGoods", controller.SelectGoods)
	r.GET("/addGood", controller.AddGood)
	r.POST("/order", controller.MakeOrder)

	_ = r.Run(":8080")
}

 

