// Logutils.go
package main

import (
	"flag"

	"github.com/golang/glog"
)

func initLog() {
	//  直接初始化，主要使服务器启动后自己直接加载，并不用命令行执行对应的参数
	flag.Set("alsologtostderr", "false") // 日志写入文件的同时，输出到stderr
	flag.Set("log_dir", "./golog")       // 日志文件保存目录
	flag.Set("v", "3")                   // 配置V输出的等级
	flag.Parse()                         // 1  解析命令行参数  go build ./hello_3_log.go   hello_3_log.exe -log_dir="./"

	glog.Flush() // 4

}

func glogInfo(str string) {
	glog.Info(str)
}

//func glogError(str string) {
//	glog.Error(str)
//}
//func glogWarning(str string) {
//	glog.Warning(str)
//}
//func glogVlog(v int, str string) {
//	glog.V(1).Infoln(str, v)
//}
func glogFlush() {
	glog.Flush()
}
