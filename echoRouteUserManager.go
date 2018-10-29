// echoRouteUserManager
package main

import (
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"strconv"
	//"time"
	"encoding/hex"

	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
)

func echoRouteUserManager() {

	e.POST(projectName+"/user/login", userLogin)

}

func userLogin(c echo.Context) error {
	glogInfo("userLogin----------------------------------------------------------------------------------------------------------")
	c.Request().ParseForm()
	//1. 获取参数
	user := &User{}
	//user.UserType = getFormParamInt64(c, "userType") //用户类型
	user.UserName = getFormParam(c, "userName")    //用户名称
	user.Password = getFormParam(c, "password")    //用户密码
	user.StoreId = getFormParamInt64(c, "storeId") //门店编号（店员必传）

	userJsonBytes, _ := json.Marshal(user)
	glogInfo("Rqdata" + string(userJsonBytes))

	//2. 参数校验
	validateErrs := validate.Struct(user)
	if validateErrs != nil {
		if _, ok := validateErrs.(*validator.InvalidValidationError); !ok {
			glogInfo(validateErrs.Error())
			errResponse := getResPonse("4603", validateErrs.Error())
			return c.JSON(http.StatusOK, errResponse)
		}
	}

	// 3.逻辑处理
	isExit, dbuser := dbQueryUserRoleByName(user.UserName)
	if isExit {
		dbuserJsonBytes, _ := json.Marshal(dbuser)
		glogInfo("Dbdata" + string(dbuserJsonBytes))
	} else {
		glogInfo("用户不存在:" + user.UserName)
		errResponse := getResPonse("4609", "用户不存在:"+user.UserName)
		return c.JSON(http.StatusOK, errResponse)
	}

	// 密码校验
	h := sha256.New()
	h.Write([]byte(dbuser.Salt + user.Password))
	//fmt.Printf("%x", h.Sum(nil))
	shaWithSaltStr := hex.EncodeToString(h.Sum(nil))
	if dbuser.Password != shaWithSaltStr {
		glogInfo("登录密码错误")
		errResponse := getResPonse("4609", "登录密码错误")
		return c.JSON(http.StatusOK, errResponse)
	}

	token, uuidErr := uuid.NewV4()
	if nil != uuidErr {
		glogInfo(uuidErr.Error())
		errResponse := getResPonse("9999", uuidErr.Error())
		return c.JSON(http.StatusOK, errResponse)
	}
	// 用户不同类型处理
	if dbuser.UserType == 1 {
		glogInfo("i01")
		succResponse := getResPonse("00", "")
		succResponse.RetMsg = "请使用收银员或门店账号登录"
		succResponse.RetCode = "4621"
		return c.JSON(http.StatusOK, succResponse)
	} else if dbuser.UserType == 2 {
		glogInfo("i02")
		succResponse := getResPonse("00", "")
		succResponse.RetMsg = "请使用收银员或门店账号登录"
		succResponse.RetCode = "4621"
		return c.JSON(http.StatusOK, succResponse)
	} else if dbuser.UserType == 3 {
		// 店长   这台设备被哪个店长登录这台设备就属于哪家店
		// 查询店铺
		isShopExit, dbstore := dbQueryStoreByUserId(dbuser.UserId)
		succResponse := getResPonse("00", "")
		if isShopExit {
			dbshopJsonBytes, _ := json.Marshal(dbstore)
			dbshopJsonStr := string(dbshopJsonBytes)
			glogInfo("Dbdata:" + dbshopJsonStr)
			bodyvalue := make(map[string]string)
			bodyvalue["storeId"] = strconv.FormatInt(dbstore.StoreId, 10)
			bodyvalue["storeName"] = dbstore.StoreName
			bodyvalue["storeAddr"] = dbstore.StoreAddr
			bodyvalue["mrchId"] = strconv.FormatInt(dbstore.MrchId, 10)
			bodyvalue["token"] = token.String()
			setNoSqlStrExpire(REDIS_MODLE+"token_"+dbuser.UserName, 5*60, bodyvalue["token"])
			succResponse.Body = bodyvalue

			return c.JSON(http.StatusOK, succResponse)
		} else {
			succResponse.RetMsg = "门店未配置，请联系管理员"
			succResponse.RetCode = "4610"
			return c.JSON(http.StatusOK, succResponse)
		}
	} else if dbuser.UserType == 4 {
		// 店员登录   判断这个店员是否是这家店的
		// 查询店铺
		isShopExit, dbcashier := dbQueryCashierByUserId(dbuser.UserId)
		succResponse := getResPonse("00", "")
		if isShopExit {
			dbshopJsonBytes, _ := json.Marshal(dbcashier)
			dbshopJsonStr := string(dbshopJsonBytes)
			glogInfo("Dbdata:" + dbshopJsonStr)
			if dbcashier.StoreId == user.StoreId {
				bodyvalue := make(map[string]string)
				bodyvalue["storeId"] = strconv.FormatInt(dbcashier.StoreId, 10)
				bodyvalue["storeName"] = dbcashier.StoreName
				bodyvalue["storeAddr"] = dbcashier.StoreAddr
				bodyvalue["cashierId"] = strconv.FormatInt(dbcashier.CashierId, 10)
				bodyvalue["cashierName"] = dbcashier.CashierName
				bodyvalue["mrchId"] = strconv.FormatInt(dbcashier.MrchId, 10)
				bodyvalue["token"] = token.String()
				setNoSqlStrExpire(REDIS_MODLE+"token_"+dbuser.UserName, 7*24*60*60, bodyvalue["token"])
				succResponse.Body = bodyvalue
				return c.JSON(http.StatusOK, succResponse)
			} else {
				succResponse.RetMsg = "非门店店员"
				succResponse.RetCode = "4612"
				return c.JSON(http.StatusOK, succResponse)
			}
		} else {
			succResponse.RetMsg = "门店未配置，请联系管理员"
			succResponse.RetCode = "4611"
			return c.JSON(http.StatusOK, succResponse)
		}
	} else {
		glogInfo("用户类型不存在:" + string(user.UserType))
		errResponse := getResPonse("4608", "用户类型不存在:"+string(user.UserType))
		return c.JSON(http.StatusOK, errResponse)
	}

	return c.String(http.StatusOK, "wait。。。")
}

type UserRegistBase_User struct {
	ShopkeeperName string ` validate:"required,max=64"`
	UserName       string ` validate:"required,max=64"`
	Password       string ` validate:"required,min=8,max=64"`
	ShopId         int64  `validate:"-"`
	Token          string ` validate:"required,max=256"`
}

func validateToken(requestToken string, userName string) (validateTokenrtn bool) {
	token, _ := getNoSqlStr(REDIS_MODLE + "token_" + userName)
	if len(token) == 0 || token != requestToken {
		glogInfo("Token过期或未登录:[" + token + "],[" + requestToken + "]")
		return false
	} else {
		return true
	}
}
