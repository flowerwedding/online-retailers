package service

import (
	"fmt"
	"log"
	"sync"
	"time"
	"what-unexpected-summer/summer-two/model"
)

//定义的一个用户结构体，在controller层放入管道
type User struct {
	UserId string
	GoodsId  uint
	Num      int
}

//就是那个用户结构体被放入的管道，传输用户类型的，最大缓存1024个
//换成消息队列就不用channel了，但是下面那个map不变，消息队列参数也改查User
var OrderChan = make(chan User, 1024)

//一个以商品id为key，商品信息的结构体item的指针为value的map
var ItemMap = make(map[uint]*Item)

//传说中的商品信息的结构体，里面有把锁，有给struct管道等等奇怪的东西
type Item struct {
	ID        uint   // 商品id
	Name      string // 名字
	Total     int    // 商品总量
	Left      int    // 商品剩余数量
	IsSoldOut bool   // 是否售罄
	leftCh    chan int
	sellCh    chan int
	done      chan struct{}
	Lock      sync.Mutex
}

// TODO 写一个定时任务，每天定时从数据库加载数据到Map
func build(){
	ticker := time.NewTicker(6 * time.Hour)

	ItemMap = make(map[uint]*Item)
	goods, _ := model.SelectGoods()

	for _,good := range goods{
		<-ticker.C
		item := &Item{
			ID:        uint(good.Model.ID),
			Name:      good.Name,
			Total:     good.Num,
			Left:      good.Num,
			IsSoldOut: false,
			leftCh:    make(chan int),
			sellCh:    make(chan int),
		}
		ItemMap[item.ID] = item

		item.Monitor()
		go item.OffShelve()
	}

}

//第一个Map也就是被初始化的那个map，为什么会有这个的存在？？？
func initMap() {
	goods, _ := model.SelectGoods()

	for _,good := range goods{
		item := &Item{
			ID:        uint(good.Model.ID),
			Name:      good.Name,
			Total:     good.Num,
			Left:      good.Num,
			IsSoldOut: false,
			leftCh:    make(chan int),
			sellCh:    make(chan int),
		}
		ItemMap[item.ID] = item

		fmt.Println(ItemMap)
		item.Monitor()
		go item.OffShelve()
	}

}

//通过参数商品id，来返回它在map里面对应的value，也就是商品信息的结构体指针
func getItem(itemId uint) *Item{
	fmt.Println(ItemMap)
	fmt.Println(ItemMap[itemId])

	return ItemMap[itemId]
}

//它在循环，还是个死循环，还是开了十个协程在一起循环的循环
//第一句话也就是它在等管道里面的信息，等信息就是它的使命，又是一个以为自己是百妖谱里的庆忌的家伙
//第二句话是用接收到的用户结构体里面的商品id为引，来找它对应的商品信息
//第三句话是也就是下面那个函数，就是买买买
func order() {
	for {
		user := <- OrderChan

		item := getItem(user.GoodsId)
		item.SecKilling(user.UserId,user.Num)
	    }
}

//这个函数带锁，一进去就锁最后才开锁，可能是怕太多人同时买这个商品造成信息错误
//判断该商品是否卖完，卖完直接退出，否正就买（一个），与此同时增加订单
func (item *Item) SecKilling(userId string,num int) {

	item.Lock.Lock()
	defer item.Lock.Unlock()
	// 等价
	// var lock = make(chan struct{}, 1}
	// lock <- struct{}{}
	// defer func() {
	// 		<- lock
	// }
	if item.IsSoldOut {
		return
	}
	item.BuyGoods(num)

	MakeOrder(userId, item.ID,num)


}

//定时下架
//timer定时器，到固定时间后执行一次，timer只触发一次，ticker隔一段时间就触发，五分钟后商品下架，并从商品序列的map中删除，且关闭它的管道。
func (item *Item) OffShelve() {
	beginTime := time.Now()
	// 获取第二天时间
	//nextTime := beginTime.Add(time.Hour * 24)
	// 计算次日零点，即商品下架的时间
	//offShelveTime := time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), 0, 0, 0, 0, nextTime.Location())
	offShelveTime := beginTime.Add(time.Second*5)
	timer := time.NewTimer(offShelveTime.Sub(beginTime))

	//现在的时间加五分钟是下架时间，然后设置定时器的间隔时间为下架时间减当前时间（五分钟）？？？
	//定时器的实质是单向通道，在最少过去时间段 d 后到期，向其自身的 C 字段发送当时的时间
	<-timer.C
	delete(ItemMap, item.ID)//delete(map,key)，使用delete内建函数从map删除一组键值对
	close(item.done)//用完channel及时关闭，不论哪种类型，关闭对应协程的channel

}
// 出售商品
func (item *Item) SalesGoods() {
	for {
		select {
		case num := <-item.sellCh://有人来买
		    if item.Left -= num; item.Left <= 0 {
				item.IsSoldOut = true
			}

			if err := model.UpdateGoodsByUserId(int(item.ID),item.Left);err != nil{
				log.Println(1)
				return
			}

		case item.leftCh <- item.Left://这个有什么意义，和下面获取剩余库存的有关系吗？？？
		case <-item.Done()://这里返回的是ok吗，这个语句执行后应该关闭管道吧？？？
			log.Println("我自闭了")
			return
		}
	}
}

//struct类型的channel，struct其实就一个普通数据类型，就是没具体的值
//struct零值就是本身，读取close的channel返回零值
func (item *Item) Done() <-chan struct{} {
	if item.done == nil {
		item.done = make(chan struct{})//这个对象一定要make出来才能使用!
	}
	d := item.done
	return d
}

//去卖商品
func (item *Item) Monitor() {
	go item.SalesGoods()//一句go 函数，开启一个协程
}

// 获取剩余库存
//这个函数好像才是那个没人调用的函数
func (item *Item) GetLeft() int {
	var left int
	left = <-item.leftCh
	return left
}

// 购买商品
//这个函数在order里面，有人来买，然后开始买就是这个函数，这个函数管道里面增加了会触动上面的出售商品的函数吧
func (item *Item) BuyGoods(num int) {
	item.sellCh <- num
}

//这个函数是初始化，只执行一次吧
func InitService() {
	initMap()
	go build()
//	for _,item := range ItemMap { //这下面两个就是一个打开协程的作用，但是执不执行等命令，这里的ItemMap只有一对键值，为什么要循环？？？
//		item.Monitor()      //为什么要用一个方法，不能直接开协程吗？？？
//		go item.OffShelve() //这个函数里面有判断是否卖完
//	}

	for i := 0; i < 10; i++ {
//		go order()//然后开了十个管道去监听是否有人买
		go OpenConsumer()
	}
}