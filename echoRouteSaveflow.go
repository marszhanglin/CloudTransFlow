// echoRouteSaveflow
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

func echoRouteSaveflow() {
	e.POST(projectName+"/uploadflow", uploadflow)
	e.POST(projectName+"/uploadflowbatch", uploadflowbatch)
}

func uploadflow(c echo.Context) error {
	glogInfo("uploadflow----------------------------------------------------------------------------------------------------------")
	c.Request().ParseForm()
	// 获取参数
	token := getHeaderParam(c, "token")
	userName := getFormParam(c, "userName")
	// 校验token
	isTokenValidateSucc := validateToken(token, userName)
	if !isTokenValidateSucc {
		errResponse := getResPonse("4613", "Token过期或未登录")
		return c.JSON(http.StatusOK, errResponse)
	}

	//获取用户信息   为了校验店铺及运营商
	isUserExit, user := dbQueryUserCashierByName(userName)
	if !isUserExit {
		errResponse := getResPonse("4609", "用户不存在("+userName+")")
		return c.JSON(http.StatusOK, errResponse)
	}
	flows := &TransFlowVO{}
	flows.DeviceSn = getFormParam(c, "deviceSn")
	flows.StoreId = getFormParam(c, "storeId")
	flows.CashierId = getFormParam(c, "cashierId")
	flows.FlowNo = getFormParam(c, "flowNo")
	flows.UploadTime = getFormParam(c, "uploadTime")
	flows.TransTime = getFormParam(c, "transTime")
	flows.TransType = getFormParam(c, "transType")
	flows.ChannelId = getFormParam(c, "channelId")
	flows.MerchantId = getFormParam(c, "merchantId")
	flows.TerminalId = getFormParam(c, "terminalId")
	flows.MerchantName = getFormParam(c, "merchantName")
	flows.Amount = getFormParam(c, "amount")
	flows.TransAmount = getFormParam(c, "transAmount")
	flows.CurrencyCode = getFormParam(c, "currencyCode")
	flows.OutOrderNo = getFormParam(c, "outOrderNo")
	flows.VoucherNo = getFormParam(c, "voucherNo")
	flows.ReferenceNo = getFormParam(c, "referenceNo")
	flows.AuthCode = getFormParam(c, "authCode")
	flows.OriVoucherNo = getFormParam(c, "oriVoucherNo")
	flows.OriOutOrderNo = getFormParam(c, "oriOutOrderNo")
	flows.OriReferenceNo = getFormParam(c, "oriReferenceNo")
	flows.OriAuthCode = getFormParam(c, "oriAuthCode")
	flows.CardNo = getFormParam(c, "cardNo")
	flows.OperatorNo = getFormParam(c, "operatorNo")
	flows.CombinationNo = getFormParam(c, "combinationNo")
	flows.CardType = getFormParam(c, "cardType")
	flows.Remark = getFormParam(c, "remark")
	flows.ExtendParams = getFormParam(c, "extendParams")

	flowsJsonBytes, _ := json.Marshal(flows)
	glogInfo("requestData:" + string(flowsJsonBytes))

	//1.校验参数,对象转换
	validateRst, validateMsg := validateTransFlow(flows)
	if !validateRst {
		errResponse := getResPonse("4603", validateMsg)
		return c.JSON(http.StatusOK, errResponse)
	}

	//2.对象转换
	vo2poTransFlowRst, flow, vo2poTransFlowMsg := vo2poTransFlow(flows)
	if !vo2poTransFlowRst {
		errResponse := getResPonse("4603", vo2poTransFlowMsg)
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3 判断运营商与店铺是否匹配
	if user.StoreId != flow.StoreId {
		errResponse := getResPonse("4615", "非该门店收银员上送")
		return c.JSON(http.StatusOK, errResponse)
	}
	//flow.CashierId = user.CashierId
	flow.MrchId = user.MrchId
	//4.业务操作  DeviceSn&flowNo 重复更新
	isExit, isExitErr, flowid := dbisFlowExit(flow)
	if isExit {
		flow.TransflowId = flowid
		// 更新
		updateErr := dbupdateFlow(flow)
		if nil == updateErr {
			succResponse := getResPonse("00", "succ")
			return c.JSON(http.StatusOK, succResponse)
		} else {
			succResponse := getResPonse("4606", updateErr.Error())
			return c.JSON(http.StatusOK, succResponse)
		}
	} else {
		// 查询失败
		if nil != isExitErr {
			errResponse := getResPonse("4606", isExitErr.Error())
			return c.JSON(http.StatusOK, errResponse)
		}
		// 保存流水
		saveErr := dbsaveFlow(flow)
		if nil == saveErr {
			succResponse := getResPonse("00", "succ")
			return c.JSON(http.StatusOK, succResponse)
		} else {
			succResponse := getResPonse("4606", saveErr.Error())
			return c.JSON(http.StatusOK, succResponse)
		}
	}
}

// 流水批上送
func uploadflowbatch(c echo.Context) error {
	glogInfo("uploadflowbatch----------------------------------------------------------------------------------------------------------")
	c.Request().ParseForm()
	// 获取参数
	token := getHeaderParam(c, "token")
	userName := getFormParam(c, "userName")
	// 校验token
	isTokenValidateSucc := validateToken(token, userName)
	if !isTokenValidateSucc {
		errResponse := getResPonse("4613", "Token过期或未登录")
		return c.JSON(http.StatusOK, errResponse)
	}

	//获取用户信息   为了校验店铺及运营商
	isUserExit, user := dbQueryUserCashierByName(userName)
	if !isUserExit {
		errResponse := getResPonse("4609", "用户不存在("+userName+")")
		return c.JSON(http.StatusOK, errResponse)
	}
	// 获取流水集合  返序列化
	flowbatch := getFormParam(c, "flowbatch")
	var flows []TransFlowVO
	flowsUnmarshalErr := json.Unmarshal([]byte(flowbatch), &flows)
	if nil != flowsUnmarshalErr {
		errResponse := getResPonse("4603", flowsUnmarshalErr.Error())
		return c.JSON(http.StatusOK, errResponse)
	}

	flowsJsonBytes, _ := json.Marshal(flows)
	glogInfo("requestData:" + string(flowsJsonBytes))

	//1.校验参数,对象转换,保存
	for _, item := range flows {
		validateRst, validateMsg := validateTransFlow(&item)
		if !validateRst {
			errResponse := getResPonse("4603", "FlowNo:"+item.FlowNo+","+validateMsg)
			return c.JSON(http.StatusOK, errResponse)
		}

		//2.对象转换
		vo2poTransFlowRst, flow, vo2poTransFlowMsg := vo2poTransFlow(&item)
		if !vo2poTransFlowRst {
			errResponse := getResPonse("4603", vo2poTransFlowMsg)
			return c.JSON(http.StatusOK, errResponse)
		}

		// 3判断运营商与店铺是否匹配
		if user.StoreId != flow.StoreId {
			errResponse := getResPonse("4615", "非该门店收银员上送")
			return c.JSON(http.StatusOK, errResponse)
		}

		//flow.CashierId = user.CashierId
		flow.MrchId = user.MrchId

		//4.业务操作  DeviceSn&flowNo 重复更新
		isExit, isExitErr, flowid := dbisFlowExit(flow)
		if isExit {
			flow.TransflowId = flowid
			// 更新
			updateErr := dbupdateFlow(flow)
			if nil != updateErr {
				succResponse := getResPonse("4606", updateErr.Error())
				return c.JSON(http.StatusOK, succResponse)
			}
		} else {
			// 查询失败
			if nil != isExitErr {
				errResponse := getResPonse("4606", isExitErr.Error())
				return c.JSON(http.StatusOK, errResponse)
			}
			// 保存流水
			saveErr := dbsaveFlow(flow)
			if nil != saveErr {
				succResponse := getResPonse("4606", saveErr.Error())
				return c.JSON(http.StatusOK, succResponse)
			}
		}
	}
	succResponse := getResPonse("00", "succ")
	return c.JSON(http.StatusOK, succResponse)

}

func validateTransFlow(flow *TransFlowVO) (validateRst bool, errMsg string) {
	validateErrs := validate.Struct(flow)
	if validateErrs != nil {
		if _, ok := validateErrs.(*validator.InvalidValidationError); !ok {
			glogInfo(validateErrs.Error())
			return false, validateErrs.Error()
		}
	}
	return true, ""
}

func vo2poTransFlow(flow *TransFlowVO) (isSucc bool, rtnflow *TransFlow, errMsg string) {
	poTransFlow := &TransFlow{}
	var shopErr, cashierErr, uploadTimeErr, transTimeErr, amountErr, transAmountErr error
	poTransFlow.StoreId, shopErr = strconv.ParseInt(flow.StoreId, 10, 64)
	poTransFlow.CashierId, cashierErr = strconv.ParseInt(flow.CashierId, 10, 64)
	poTransFlow.FlowNo = flow.FlowNo
	poTransFlow.DeviceSn = flow.DeviceSn
	poTransFlow.UploadTime, uploadTimeErr = time.ParseInLocation("20060102150405", flow.UploadTime, cstZone)
	poTransFlow.TransTime, transTimeErr = time.ParseInLocation("20060102150405", flow.TransTime, cstZone)
	poTransFlow.TransType = flow.TransType
	poTransFlow.ChannelId = flow.ChannelId
	poTransFlow.MerchantId = flow.MerchantId
	poTransFlow.TerminalId = flow.TerminalId
	poTransFlow.MerchantName = flow.MerchantName
	fmt.Print(strconv.ParseInt(flow.Amount, 10, 64))
	poTransFlow.Amount, amountErr = strconv.ParseInt(flow.Amount, 10, 64)
	poTransFlow.TransAmount, transAmountErr = strconv.ParseInt(flow.TransAmount, 10, 64)
	poTransFlow.CurrencyCode = flow.CurrencyCode
	poTransFlow.OutOrderNo = flow.OutOrderNo
	poTransFlow.VoucherNo = flow.VoucherNo
	poTransFlow.ReferenceNo = flow.ReferenceNo
	poTransFlow.AuthCode = flow.AuthCode
	poTransFlow.OriVoucherNo = flow.OriVoucherNo
	poTransFlow.OriOutOrderNo = flow.OriOutOrderNo
	poTransFlow.OriReferenceNo = flow.OriReferenceNo
	poTransFlow.OriAuthCode = flow.OriAuthCode
	poTransFlow.CardNo = flow.CardNo
	poTransFlow.OperatorNo = flow.OperatorNo
	poTransFlow.CombinationNo = flow.CombinationNo
	poTransFlow.CardType = flow.CardType
	poTransFlow.Remark = flow.Remark
	poTransFlow.ExtendParams = flow.ExtendParams
	if nil != shopErr {
		return false, poTransFlow, shopErr.Error()
	} else if nil != cashierErr {
		return false, poTransFlow, cashierErr.Error()
	} else if nil != uploadTimeErr {
		return false, poTransFlow, uploadTimeErr.Error()
	} else if nil != transTimeErr {
		return false, poTransFlow, transTimeErr.Error()
	} else if nil != amountErr {
		return false, poTransFlow, amountErr.Error()
	} else if nil != transAmountErr {
		return false, poTransFlow, transAmountErr.Error()
	} else {
		return true, poTransFlow, ""
	}
}

type TransFlowVO struct {
	Id             string `json:"id" validate:"-"`
	StoreId        string `json:"storeId" validate:"required,gte=1"`
	CashierId      string `json:"cashierId" validate:"required,gte=1"`
	FlowNo         string `json:"flowNo" validate:"required,max=64"`
	DeviceSn       string `json:"deviceSn" validate:"required,max=64"`
	UploadTime     string `json:"uploadTime" validate:"required,max=64"`
	TransTime      string `json:"transTime" validate:"required,max=64"`
	TransType      string `json:"transType" validate:"required,oneof= 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 60 61 62 63 64 65 ,max=64"`
	ChannelId      string `json:"channelId" validate:"required,oneof=acquire wxpay alipay jdpay unionscan unionpaycode boc ccb abc icbc cmbc cmb citic,max=64"`
	MerchantId     string `json:"merchantId" validate:"max=64"`
	TerminalId     string `json:"terminalId" validate:"max=64"`
	MerchantName   string `json:"merchantName" validate:"max=64"`
	Amount         string `json:"amount" validate:"required,gte=1"`
	TransAmount    string `json:"transAmount" validate:"required,gte=1"`
	CurrencyCode   string `json:"currencyCode" validate:"required,max=64"`
	OutOrderNo     string `json:"outOrderNo" validate:"required,max=64"`
	VoucherNo      string `json:"voucherNo" validate:"required,max=64"`
	ReferenceNo    string `json:"referenceNo" validate:"required,max=64"`
	AuthCode       string `json:"authCode" validate:"max=64"`
	OriOutOrderNo  string `json:"oriOutOrderNo" validate:"max=64"`
	OriVoucherNo   string `json:"oriVoucherNo" validate:"max=64"`
	OriReferenceNo string `json:"oriReferenceNo" validate:"max=64"`
	OriAuthCode    string `json:"oriAuthCode" validate:"max=64"`
	CardNo         string `json:"cardNo" validate:"max=64"`
	OperatorNo     string `json:"operatorNo" validate:"max=64"`
	CombinationNo  string `json:"combinationNo" validate:"max=64"`
	CardType       string `json:"cardType" validate:"max=64"`
	Remark         string `json:"remark" validate:"max=256"`
	ExtendParams   string `json:"extendParams" validate:"max=512"`
}
