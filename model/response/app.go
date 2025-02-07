package response

import (
	"fmt"
	"github.com/yunhanshu-net/runcher/pkg/stringsx"
	"github.com/yunhanshu-net/runcher/pkg/timex"
	"strconv"
	"time"
)

type UnInstallInfo struct {
}

type RollbackVersion struct {
}

type InstallInfo struct {
	TempPath     string `json:"temp_path"`     //软件安装时候临时目录，下载到该目录，然后copy到所属目录
	RootPath     string `json:"root_path"`     //存储根路径
	StoreRoot    string `json:"store_root"`    //云存储根路径
	Pc           string `json:"pc"`            //软件平台
	Name         string `json:"name"`          //软件名称
	FullName     string `json:"full_name"`     //软件名称,带后缀
	User         string `json:"user"`          //所属用户
	DownloadPath string `json:"download_path"` //软件的云端地址
	//InstallPath  string //安装后的所属目录
	Version string `json:"version"` //安装的软件版本

	Other map[string]interface{} `json:"other"`
}

type UpdateVersion struct {
}

//type Response struct {
//	Code    int         `json:"code"`
//	Message string      `json:"message"`
//	Data    interface{} `json:"data"`
//}

type RunTime struct {
	MEM int
}

type Run struct {
	//StatusCode  int         `json:"status_code"`
	//Msg         string      `json:"msg"`
	//ContentType string      `json:"content_type"`
	//HasFile     bool        `json:"has_file"`
	//FilePath    string      `json:"path"`
	//DeleteFile  bool        `json:"delete_file"`
	//Data        interface{} `json:"data"`
	//
	//CallCostTime     time.Duration `json:"-"`
	//ResponseMetaData string        `json:"-"`

	StatusCode int    `json:"status_code"`
	Msg        string `json:"msg"`
	//ContentType    string      `json:"content_type"`
	HasFile        bool                   `json:"has_file"`
	FilePath       string                 `json:"path"`
	DeleteFile     bool                   `json:"delete_file"`
	DeleteFileTime int                    `json:"delete_file_time"` //-1 不删除文件，0响应成功后立刻删除文件，>0是时间戳给出具体时间戳，达到该时间戳时刻系统会自动清理该文件
	Body           map[string]interface{} `json:"data"`

	Header map[string]string `json:"header"` // response header

	//meta data
	CallCostTime     time.Duration `json:"-"` // 执行引擎发起调用到程序执行结束的总耗时
	ResponseMetaData string        `json:"-"`

	RunTime RunTime
}

func (r *Run) GetContentType() string {
	if r.Header != nil {
		return r.Header["Content-Type"]
	}
	return ""
}

type RuntimeMetaData struct {
	FuncRunTime  int
	CallCostTime time.Duration
}

func (r *Run) GetResponseMetaData() *RuntimeMetaData {
	rm := &RuntimeMetaData{}
	funcRunTimeList := stringsx.ParserHtmlTagContent(r.ResponseMetaData, "UserCost")
	if len(funcRunTimeList) > 0 {
		funcRunTime := funcRunTimeList[0]
		i, err := strconv.ParseInt(funcRunTime, 10, 64)
		if err == nil {
			rm.FuncRunTime = int(i)
		}
	}
	rm.CallCostTime = r.CallCostTime

	return rm
}
func (r *RuntimeMetaData) Print(rsp *Run) {
	timex.Println(time.Duration(r.FuncRunTime)*time.Nanosecond, "函数执行")
	timex.Println(r.CallCostTime, "程序调用")
	if rsp.RunTime.MEM > 1024*1024 {
		fmt.Printf("内存占用 %.2f MB", float64(rsp.RunTime.MEM)/1024/1024)
	} else {
		fmt.Printf("内存占用 %.2f KB", float64(rsp.RunTime.MEM)/1024)
	}

}
