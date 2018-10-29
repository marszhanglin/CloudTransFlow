// echoRoute
package main

import (
	"CloudMis_TransFlow/fconf"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func echoRoute() {
	e.POST("/mtms", test)
	e.GET("/flush", flushLog)
	echoRouteSaveflow()
	echoRouteUserManager()
}

func flushLog(c echo.Context) error {
	glogFlush()
	return c.String(http.StatusOK, "ok")
}

func test(c echo.Context) error {
	c.Request().ParseForm()

	form := c.Request().PostForm["form"]
	if len(form) == 0 {
		return c.String(http.StatusOK, "form:null")
	} else {
		fmt.Println("form:" + form[0])
	}

	errs := validate.Var(form, "required")
	if errs != nil {
		fmt.Println(errs.Error()) // output: Key: "" Error:Field validation for "" failed on the "email" tag
		return c.String(http.StatusOK, errs.Error())
	}

	header := c.Request().Header.Get("header")
	fmt.Println("header:" + header)
	rsp := &ResponseBody{}
	rsp.RetCode = "00"
	rsp.RetMsg = "succ"
	bodyvalue := make(map[string]string)
	bodyvalue["a"] = "a"
	bodyvalue["b"] = "b"
	bodyvalue["c"] = "c"
	bodyvalue["c"] = "1"
	rsp.Body = bodyvalue
	//textIni()
	return c.JSON(http.StatusOK, rsp)
}

func textIni() {
	c, err := fconf.NewFileConf("./conf.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(c.String("mysql.db1.Host"))
	fmt.Println(c.String("mysql.db1.Name"))
	fmt.Println(c.String("mysql.db1.User"))
	fmt.Println(c.String("mysql.db1.Pwd"))
	fmt.Println(c.String("mysql.db1.Pwd"))
	iniDataJsonBytes, _ := json.Marshal(c)
	glogInfo("iniData:" + string(iniDataJsonBytes))
	// 取得配置时指定类型
	port, err := c.Int("mysql.db1.Port")
	if err != nil {
		panic(err)
	}
	fmt.Println(port)
}

func getFormParam(c echo.Context, key string) string {

	value := c.Request().PostForm[key]
	if value == nil {
		return ""
	} else {
		return value[0]
	}
}

func getHeaderParam(c echo.Context, key string) string {
	value := c.Request().Header.Get(key)
	return value
}

func getFormParamInt64(c echo.Context, key string) int64 {

	value := c.Request().PostForm[key]
	if value == nil {
		return 0
	} else {
		valueint64, _ := strconv.ParseInt(value[0], 10, 64)
		return valueint64
	}
}

func getResPonse(retCode string, retMsg string) *ResponseBody {
	response := &ResponseBody{}
	response.RetCode = retCode
	response.RetMsg = retMsg
	return response
}

type ResponseBody struct {
	RetCode string            `json:"retCode" xml:"retCode_"`
	Body    map[string]string `json:"body" xml:"body_"`
	RetMsg  string            `json:"retMsg" xml:"retMsg_"`
}

func getBaseResPonse(retCode string, retMsg string) *Response {
	response := &Response{}
	response.RetCode = retCode
	response.RetMsg = retMsg
	glogInfo(retCode + ":" + retMsg)
	return response
}

type Response struct {
	RetCode string                 `json:"retCode" xml:"retCode_"`
	Body    map[string]interface{} `json:"body" xml:"body_"`
	RetMsg  string                 `json:"retMsg" xml:"retMsg_"`
}

//#string到int
//int,err:=strconv.Atoi(string)
//#string到int64
//int64, err := strconv.ParseInt(string, 10, 64)
//#int到string
//string:=strconv.Itoa(int)
//#int64到string
//string:=strconv.FormatInt(int64,10)
