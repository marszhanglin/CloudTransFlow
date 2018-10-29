package main

import (
	"time"
)

func dbQueryCashierByUserId(userId int64) (isExit bool, cashier Cashier) {

	querydatas := []Cashier{}
	db.Raw("SELECT t1.cashier_id , t1.cashier_name , t1.store_id , t2.store_name,t2.store_addr, t1.user_id FROM `t_cashier` t1 left join t_store t2 on t2.store_id=t1.store_id   where t1.user_id=?", userId).Scan(&querydatas)
	if len(querydatas) > 0 {
		return true, querydatas[0]
	} else {
		return false, querydatas[0]
	}
}

type Cashier struct {
	CashierId   int64     `gorm:"column:cashier_id;primary_key" validate:"-"`
	CashierName string    `gorm:"column:cashier_name" validate:"required,max=64"`
	StoreName   string    `gorm:"column:store_name" validate:"required,max=64"`
	StoreAddr   string    `gorm:"column:store_addr" validate:"required,max=64"`
	StoreId     int64     `gorm:"column:store_id" validate:"-"`
	MrchId      int64     `gorm:"column:mrch_id" validate:"-"`
	CreTime     time.Time `gorm:"column:cre_time" validate:"-"`
	UpdTime     time.Time `gorm:"column:upd_time" validate:"-"`
	Remark      string    `gorm:"column:remark" validate:"max=256"`
}
