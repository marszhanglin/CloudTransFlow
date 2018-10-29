package main

import (
	"time"
)

func dbQueryStoreByUserId(userId int64) (isExit bool, shop Store) {

	querydatas := []Store{}
	db.Raw("SELECT t1.store_id ,t1.store_name,t1.store_addr, t1.user_id , t1.mrch_id FROM `t_store` t1 where t1.user_id=?", userId).Scan(&querydatas)
	if len(querydatas) > 0 {
		return true, querydatas[0]
	} else {
		return false, querydatas[0]
	}
}

type Store struct {
	StoreId   int64     `gorm:"column:store_id;primary_key" validate:"-"`
	StoreName string    `gorm:"column:store_name" validate:"required,max=64"`
	StoreAddr string    `gorm:"column:store_addr" validate:"required,max=64"`
	MrchId    int64     `gorm:"column:mrch_id" validate:"-"`
	CreTime   time.Time `gorm:"column:cre_time" validate:"-"`
	UpdTime   time.Time `gorm:"column:upd_time" validate:"-"`
	Remark    string    `gorm:"column:remark" validate:"max=256"`
}
