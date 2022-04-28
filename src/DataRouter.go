package src

/*
統一管理路由
*/

import (
	"golangAPI/service"

	"github.com/gin-gonic/gin"
)

func FindDataRouter(r *gin.RouterGroup) {
	rawdata := r.Group("/mySQL")                     //搜尋mySQL資料庫
	rawdata.GET("/", service.FindAllData)            //搜尋全部
	rawdata.GET("/search", service.FindSelectedData) //可在GET中加參數做搜尋條件，不加也可以
	rawdata.GET("/download", service.DownloadCsv)    //下載帳號最新查詢的csv
	edata := r.Group("/elasticsearch")               //搜尋elasticsearch
	edata.GET("/", service.SearchForALL)             //搜尋全部
	edata.GET("/search", service.SearchByParm)       //可以在GET中加參數做搜尋條件，不加也可以
}
