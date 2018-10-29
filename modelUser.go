// echoRouteSaveflow
package main

import (
	"time"
)

func dbsaveUser(user *User) error {
	//tx := db.Begin()
	user.UserId = getTimeUUID()
	if err := db.Table("t_user").Create(&user).Error; err != nil {
		glogInfo(err.Error())
		//tx.Rollback()
		return err
	}
	//tx.Commit()
	return nil
}

type UserRole struct {
	UserId   int64  `gorm:"column:user_id;primary_key" validate:"-"`
	UserName string `gorm:"column:username" validate:"required,max=64"`
	Password string `gorm:"column:password" validate:"required,max=64"`
	Salt     string `gorm:"column:salt" validate:"-"`
	UserType int64  `gorm:"column:role_id" validate:"required"`
}

func dbQueryUserRoleByName(userName string) (isExit bool, userrtn UserRole) {

	userRole := []UserRole{}
	db.Raw("SELECT t1.user_id ,t1.username,t1.password, t1.salt,t2.role_id FROM `sys_user` t1 LEFT JOIN sys_user_role t2 on  t2.user_id=t1.user_id where t1.username=?", userName).Scan(&userRole)
	if len(userRole) > 1 {
		glogInfo("存在多个同名用户")
		return false, userRole[0]
	} else if len(userRole) == 1 {
		return true, userRole[0]
	} else {
		user := &UserRole{}
		return false, *user
	}
}

type UserCashier struct {
	UserId    int64  `gorm:"column:user_id;primary_key" validate:"-"`
	UserName  string `gorm:"column:username" validate:"required,max=64"`
	StoreId   int64  `gorm:"column:store_id" validate:"-"`
	CashierId int64  `gorm:"column:cashier_id" validate:"-"`
	MrchId    int64  `gorm:"column:mrch_id" validate:"-"`
}

func dbQueryUserCashierByName(userName string) (isExit bool, data UserCashier) {

	userCashier := []UserCashier{}
	db.Raw("SELECT t1.user_id, t1.username,  t2.store_id,  t2.cashier_id,  t2.mrch_id FROM `sys_user` t1 LEFT JOIN t_cashier t2 ON t2.user_id = t1.user_id WHERE t1.username =?", userName).Scan(&userCashier)
	if len(userCashier) > 1 {
		glogInfo("存在多个同名用户")
		return false, userCashier[0]
	} else if len(userCashier) == 1 {
		return true, userCashier[0]
	} else {
		user := &UserCashier{}
		return false, *user
	}
}

func dbQueryUserByShopId(shopId int64, userType int64) ([]QueryUser, error) {

	queryUsers := []QueryUser{}
	err := db.Table("t_user").Select("user_id,user_name,user_type,shop_id ").Where(" shop_id = ? and user_type = ? ", shopId, userType).Find(&queryUsers).Error
	return queryUsers, err
}

func dbDeleteUserByName(userName string) error {
	tx := db.Begin()
	if err := tx.Table("t_user").Where(" user_name = ? ", userName).Delete(User{}).Error; err != nil {
		glogInfo(err.Error())
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

type User struct {
	UserId     int64     `gorm:"column:user_id;primary_key" validate:"-"`
	UserName   string    `gorm:"column:user_name" validate:"required,max=64"`
	Password   string    `gorm:"column:password" validate:"required,max=64"`
	UserStatus int64     `gorm:"column:user_status" validate:"-"`
	UserType   int64     `gorm:"column:user_type" validate:"-"`
	OperatorId int64     `gorm:"column:operator_id" validate:"max=64"`
	StoreId    int64     `gorm:"column:store_id" validate:"-"`
	CreTime    time.Time `gorm:"column:cre_time" validate:"-"`
	UpdTime    time.Time `gorm:"column:upd_time" validate:"-"`
	Remark     string    `gorm:"column:remark" validate:"max=256"`
}

type QueryUser struct {
	UserId     int64  `gorm:"column:user_id;primary_key" validate:"-"`
	UserName   string `gorm:"column:user_name" validate:"required,max=64"`
	UserStatus int64  `gorm:"column:user_status" validate:"-"`
	UserType   int64  `gorm:"column:user_type" validate:"required"`
	OperatorId int64  `gorm:"column:operator_id" validate:"max=64"`
	ShopId     int64  `gorm:"column:shop_id" validate:"-"`
}
