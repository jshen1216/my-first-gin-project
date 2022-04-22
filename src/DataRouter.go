package src

/*
統一管理路由
*/

import (
	"golangAPI/service"

	"github.com/gin-gonic/gin"
)

func FindDataRouter(r *gin.RouterGroup) {
	rawdata := r.Group("/accesslog")                 //搜尋table:accesslog
	rawdata.GET("/", service.FindAllData)            //搜尋全部
	rawdata.GET("/search", service.FindSelectedData) //可在GET中加參數做搜尋條件，不加也可以
	rawdata.GET("/download", service.DownloadCsv)    //下載帳號最新查詢的csv
}
