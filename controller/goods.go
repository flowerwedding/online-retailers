package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"what-unexpected-summer/summer-two/service"
)

func SelectGoods(ctx *gin.Context) {
	goods := service.SelectGoods()
	ctx.JSON(http.StatusOK, gin.H{
		"status": 200,
		"info": "success",
		"data": struct {
			Goods []service.Goods `json:"goods"`
		}{goods},
	})
}

func AddGood(ctx *gin.Context){
	name := ctx.PostForm("name")
	price := ctx.PostForm("price")
	num := ctx.PostForm("num")
	prices,_ := strconv.Atoi(price)
	nums,_ := strconv.Atoi(num)

	err := service.AddGoods(name,nums,prices)
	if err != nil {
		ctx.JSON(200,gin.H{"status":10001,"message":"fail"})
		return
	}
	ctx.JSON(200,gin.H{"status":10001,"message":"success"})
}

 

