package service

import (
	"golangAPI/pojo"

	"fmt"
	"net/http"
	"time"

	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// parameters of mysql at localhost
// which username and passwd are default setting
const (
	USERNAME = "root"
	PASSWORD = "root"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 8889
	DATABASE = "db_log"
)

//find all data
func FindAllData(c *gin.Context) {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)

	db, err := sql.Open("mysql", conn)
	if err != nil {
		c.JSON(http.StatusBadGateway, "Connect DB failed")
		return
	}
	if err := db.Ping(); err != nil {
		c.JSON(http.StatusBadGateway, "Connect DB failed")
		return
	}

	rows, err := db.Query("SELECT Time, IP, Message FROM accesslog")

	if err != nil {
		c.JSON(http.StatusBadRequest, "Connot Find Data")
		return
	}

	defer rows.Close()
	defer db.Close()

	var dataList []pojo.Rawdata

	for rows.Next() {
		data := new(pojo.Rawdata)
		err = rows.Scan(&data.Time, &data.IP, &data.Message)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Cannot Find Data")
			return
		}
		dataList = append(dataList, *data)
	}
	c.JSON(http.StatusOK, dataList)
}

//find data (可以取代find all data)
func FindSelectedData(c *gin.Context) {

	//抓日期（日期格式：dd/MMM/yyy:HH:mm:ss
	maxDateTime := c.Query("DueTime")
	minDateTime := c.Query("StartTime")
	//抓IP
	ip := c.Query("IP")
	//抓關鍵字
	keyword := c.Query("KeyWord")

	//設定連接資料，連接不到回傳報錯訊息
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		c.JSON(http.StatusBadGateway, "Connect DB failed")
		return
	}
	if err := db.Ping(); err != nil {
		c.JSON(http.StatusBadGateway, "Connect DB failed")
		return
	}

	// 依Get的參數條件設定SQL
	sql := "SELECT TIME, IP, Message FROM accesslog "
	OriginalForm := "02/Jan/2006:15:04:05"
	SqlForm := "2006-01-02 15:04:05"
	if len(maxDateTime) != 0 {
		t1, _ := time.Parse(OriginalForm, maxDateTime)
		sql = sql + "WHERE TIME <= '" + t1.Format(SqlForm) + "' "
	} else if len(maxDateTime) == 0 {
		t1 := time.Now().Format(SqlForm)
		sql = sql + "WHERE TIME <= '" + t1 + "' "
	}
	if len(minDateTime) != 0 {
		t2, _ := time.Parse(OriginalForm, minDateTime)
		sql = sql + "AND TIME >= '" + t2.Format(SqlForm) + "' "
	}
	if len(ip) != 0 {
		sql = sql + "AND IP = '" + ip + "' "
	}
	if len(keyword) != 0 {
		sql = sql + "AND Message LIKE '%" + keyword + "%'"
	}
	sql = sql + ";"

	//查詢符合欄位，並解析
	rows, err := db.Query(sql)

	if err != nil {
		c.JSON(http.StatusBadRequest, "Cannot Find Data, SQL: "+sql)
		return
	}

	defer rows.Close()
	defer db.Close()

	var dataList []pojo.Rawdata

	for rows.Next() {
		data := new(pojo.Rawdata)
		err = rows.Scan(&data.Time, &data.IP, &data.Message)
		if err != nil {
			c.JSON(http.StatusBadRequest, "Cannot Find Data, SQL: "+sql)
			return
		}
		dataList = append(dataList, *data)
	}
	StructToCsv("account.csv", dataList)
	c.JSON(http.StatusOK, dataList)
}
