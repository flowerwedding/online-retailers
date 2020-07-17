package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"what-unexpected-summer/summer-two/service"
)


func MakeOrder(ctx *gin.Context) {
	userId := ctx.PostForm("userId")
	goodsId := ctx.PostForm("goodsId")
	nums := ctx.PostForm("num")
	itemId,_ := strconv.Atoi(goodsId)
	num,_ := strconv.Atoi(nums)
/*    service.OrderChan <- service.User{
		UserId:  userId,
		GoodsId: uint(itemId),
		Num : num,
	}
*/    service.Order(service.User{UserId:  userId, GoodsId: uint(itemId),Num :num})
    ctx.JSON(200, gin.H{
    	"status": 200,
    	"info": "success",
	})
}

 

