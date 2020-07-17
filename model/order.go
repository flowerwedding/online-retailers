package model

import "github.com/jinzhu/gorm"

type Order struct {
	gorm.Model
	UserID  string
	GoodsID uint
	Num     int
}

// 下单
func (order *Order)MakeOrder() error{
	return DB.Model(&Order{}).Create(&order).Error
}


// 查询订单
func GetOrderByUserID(userId string) (orders []Order, err error){
	err = DB.Table("orders").Where("user_id = ?",userId).Find(&orders).Error
 	if err != nil {
 		return nil, err
	}
	return orders,nil
}

func UpdateGoodsByUserId(id int,num int) ( err error){
	err = DB.Model(Goods{}).Where("id = ?",id).Update("num",num).Error
	if err != nil {
		return err
	}
	return nil
}