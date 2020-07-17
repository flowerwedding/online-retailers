package service

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
	"what-unexpected-summer/summer-two/model"
)

func failError(err error,msg string){
	if err != nil{
		log.Fatalf("%s: %s", msg, err)
	}
}

func Order(user User){
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failError(err, "Can't connect to MQ")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	failError(err, "Can't create a Channel")
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("goodList",false,true,false,false,nil)
	failError(err, "Could not declare queue")

	rand.Seed(time.Now().UnixNano())
//	good := model.Order{UserID:"1", Num:2, GoodsID:3}
	good := model.Order{UserID:user.UserId, Num:user.Num, GoodsID:user.GoodsId}
	body, err:= json.Marshal(good)
	if err != nil{
		failError(err, "Error encoding JSON")
	}

	//发布消息，其中amqp.Publishing的DeliveryMode如果设为amqp.Persistent则消息会持久化。
	//需要注意的是如果需要消息持久化Queue也是需要设定为持久化才有效
	err = amqpChannel.Publish("",queue.Name,false,false,amqp.Publishing{
		DeliveryMode : amqp.Persistent,
		ContentType : "text/plain",
		Body : body,
	})
	if err != nil{
		log.Fatalf("Error publishing message: %s",err)
	}

	log.Printf("AddGood: %s",string(body))
}