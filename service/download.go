package service

import (
	"encoding/csv"
	"fmt"
	"golangAPI/pojo"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// GET查詢資料庫時，順便以帳號名稱將查詢結果存到tempscv中，以免下載時需要再查詢一次
func StructToCsv(filename string, Data []pojo.Rawdata) {
	filename = "./tempcsv/" + filename
	newFile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		newFile.Close()
	}()
	// 寫入UTF-8
	newFile.WriteString("\xEF\xBB\xBF") // 寫入UTF-8 BOM，防止中文亂碼
	// 寫資料到csv檔案
	w := csv.NewWriter(newFile)
	header := []string{"時間", "IP", "RawData"} //標題
	data := [][]string{
		header,
	}
	for _, v := range Data {
		context := []string{
			cast.ToString(v.Time),
			v.IP,
			v.Message,
		}
		data = append(data, context)
	}
	w.WriteAll(data) // WriteAll方法使用Write方法向w寫入多條記錄，並在最後呼叫Flush方法清空快取。
	w.Flush()
}

// 下載account在tempcsv中的csv
// @Summary      Download csv
// @Description  Download csv file which the account just searched
// @Tags         Download
// @Produce      text/csv
// @Param        account  query string true "account"
// @Success      200  {string} string "時間,IP,RawData"
// @Router /download/ [get]
func DownloadCsv(c *gin.Context) {
	// 讀取account以抓取文件
	account := c.Query("account")
	if len(account) == 0 {
		c.String(http.StatusBadRequest, "請提供帳號")
		return
	}
	filename := "./tempcsv/" + account + ".csv"
	if !checkFileIsExist(filename) {
		c.String(http.StatusBadRequest, "此帳號未有查詢文件可下載")
	}
	c.File(filename)
}

//判斷文件是否存在，存在則return True，不存在則return False
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
