package service

import (
	"log"
	"what-unexpected-summer/summer-two/model"
)

type Goods struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Num   int    `json:"num"`
}

// 添加商品
func AddGoods(name string,num int,price int) error{
	// TODO
	newgood := model.Goods{
		Name:  name,
		Price: price,
		Num:   num,
	}

    err :=  newgood.AddGoods()
    if err != nil {
    	return err
    }

	return nil
}

func SelectGoods() (goods []Goods) {
	_goods, err := model.SelectGoods()
	if err != nil {
		log.Printf("Error get goods info. Error: %s", err)
	}
	for _, v := range _goods {
		good := Goods{
			ID:    uint(v.ID),
			Name:  v.Name,
			Price: v.Price,
			Num:   v.Num,
		}
		goods = append(goods, good)
	}
	return goods
}
