package service

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strconv"
	"strings"
	"what-unexpected-summer/summer-two/model"
)

func handleError(err error,msg string){
	if err != nil{
		log.Fatalf("%s:%s",msg,err)
	}
}

func OpenConsumer() {
	//链接RabbitMQ，此处链接已经成为我们抽象了socket的链接，同时为我们处理了协议版本号和身份证号等
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err,"Can't connect to MQ")
	defer conn.Close()

	//创建通道，在使用其他API完成任务的时候应该如此创建通道
	amqpChannel, err := conn.Channel()
	handleError(err,"Can't create a ampqChannel")
	defer amqpChannel.Close()

	//声明队列，参数：队列名字，是否持久化，是否自动删除队列，是否排外，是否等待服务器返回，相关参数（一般为nil）
	//是否持久化：一般情况下队列的默认声明是存放到内存中，如果rabbitmq重启就会丢失。因此如果想重启不丢失（具有持久化），则保存到erlang自带的mnesia数据库，重启后会读取该数据库
	//是否自动删除队列：当最后一个消费者断开开连接后队列是否自动删除，当consumers = 0 时自动删除
	//是否排外：一、当链接关闭时connection.close()该队列是否自动删除。二、该队列是否是私有的private，如果不排外的，可以使用两个消费者都同时访问一个队列；如果是排外，会对当前队列枷锁，其他通道channel不能访问，一个队列就一个消费者
	queue, err := amqpChannel.QueueDeclare("goodList",false,true,false,false,nil)
	handleError(err, "Could not declare 'add' queue")

	//qos服务质量保证功能，参数：消费端体现，是否应用于channel
	//消费端体现：一次最多处理多少条消息，基本为1
	//是否应用于channel，基本用false
	err = amqpChannel.Qos(1,0,false)
	handleError(err,"Could not configue QoS")

	//接受消息，参数：消费的队列名字，消费者，自动应答是否排外，是否本地，是否等待，相关参数
	//其中Auto ack可以设置为true。如果设为true则消费者一接收到就从queue中去除了，如果消费者处理消息中发生意外该消息就丢失了。
	//如果Auto ack设为false。consumer在处理完消息后，调用msg.Ack(false)后消息才从queue中去除。即便当前消费者处理该消息发生意外，只要没有执行msg.Ack(false)那该消息就仍然在queue中，不会丢失。
	//生成的Queue在生成是设定的参数，下次使用时不能更改设定参数，否则会报错
	messageChannel, err := amqpChannel.Consume(queue.Name,"",false,false,false,false,nil)
	handleError(err, "Could not register consumer")

	//主要用来防止主进程窗口退出
	stopChan := make(chan bool)

	go func(){
		log.Printf("Consumer ready,PID: %d",os.Getpid())
		for d := range messageChannel{
			log.Printf("Reeived a message: %s",string(d.Body))

			good := &model.Order{}
			err := json.Unmarshal(d.Body, good)
			if err != nil {
				log.Printf("Error decoding JSON: %s",err)
			}
			log.Printf("Good: %s",string(d.Body))


			var s[10] string
			messages := strings.Split(string(d.Body),",")
			for u,message := range messages{
				st := strings.Split(message,":")
				s[u] = st[1]
			}

			ss,_ := strconv.Atoi(s[5])
			sst,_ := strconv.Atoi(s[6])
			item := getItem(uint(ss))
			item.SecKilling(s[4],sst)

			if err := d.Ack(false); err != nil{
				log.Printf("Error acknowledging message : %s",err)
			}else{
				log.Printf("Acknowledeged message")
			}
		}
	}()

	<-stopChan
}