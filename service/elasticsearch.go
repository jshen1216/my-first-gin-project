package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golangAPI/pojo"
	"net/http"
	"net/url"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
)

var es *elasticsearch.Client
var host string = "http://140.137.219.56:9200"

// 查詢全部資料（最多回傳 1,000筆）
// @Summary      Find All Data
// @Description  get Time, IP and RawData in ElasticSearch
// @Tags         ES
// @Produce      json
// @Success      200  {array}   pojo.Rawdata
// @Router /elasticsearch/ [get]
func SearchForALL(c *gin.Context) {
	// 連線Elasticsearch (without username & password)
	var err error
	cfg := elasticsearch.Config{
		Addresses: []string{host},
	}
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		c.JSON(http.StatusBadGateway, "Elasticsearch連線失敗，原因："+err.Error()) // 連線失敗
	}
	var r map[string]interface{}
	// DSL （只有抓1,000筆資料）
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"size":             1000,
		"from":             0,
		"track_total_hits": true,
	}
	jsonBody, _ := json.Marshal(query)
	req := esapi.SearchRequest{
		Index: []string{"accesslog"}, // 索引名稱
		Body:  bytes.NewReader(jsonBody),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	var dataList []pojo.Rawdata
	// 迴圈抓出hits中的資料
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		data := new(pojo.Rawdata)
		source := hit.(map[string]interface{})["_source"]
		data.Time = fmt.Sprintf("%v", source.(map[string]interface{})["@timestamp"])
		data.IP = fmt.Sprintf("%v", source.(map[string]interface{})["IP"])
		data.Message = fmt.Sprintf("%v", source.(map[string]interface{})["RawData"])
		dataList = append(dataList, *data)
	}
	StructToCsv("account.csv", dataList)
	c.JSON(http.StatusOK, dataList)
	//c.JSON(http.StatusOK, len(dataList))
}

// 條件搜尋（最多回傳 1,000筆）
// @Summary      Find Data under conditional search
// @Description  get Time, IP and RawData in ElasticSearch
// @Tags         ES
// @Produce      json
// @Param        StartTime query string false "起始時間"
// @Param        DueTime   query string false "結束時間"
// @Param        IP        query string false "IP"
// @Param        KeyWord   query string false "關鍵字搜尋"
// @Success      200  {array}   pojo.Rawdata
// @Router /elasticsearch/search [get]
func SearchByParm(c *gin.Context) {
	// 連線Elasticsearch (without username & password)
	var err error
	cfg := elasticsearch.Config{
		Addresses: []string{host},
	}
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		c.JSON(http.StatusBadGateway, "Elasticsearch連線失敗，原因："+err.Error()) // 連線失敗
	}
	var r map[string]interface{}

	// 抓取時間條件
	maxDateTime := c.Query("DueTime")
	maxDateTime, _ = url.QueryUnescape(maxDateTime) //如果是swagger ui送的request需要解碼
	minDateTime := c.Query("StartTime")
	minDateTime, _ = url.QueryUnescape(minDateTime) //如果是swagger ui送的request需要解碼
	OriginalForm := "02/Jan/2006:15:04:05"
	dslForm := "2006-01-02 15:04:05"
	var t1, t2 string
	// 結束時間若沒設定則為Now
	if len(maxDateTime) != 0 {
		dueTime, _ := time.Parse(OriginalForm, maxDateTime)
		t1 = dueTime.Format(dslForm)
	} else if len(maxDateTime) == 0 {
		t1 = time.Now().Format(dslForm)
	}
	// 起始時間若沒設定則為 Now - 10years
	if len(minDateTime) != 0 {
		startTime, _ := time.Parse(OriginalForm, minDateTime)
		t2 = startTime.UTC().Format(dslForm)
	} else if len(minDateTime) == 0 {
		t2 = time.Now().AddDate(-10, 0, 0).Format(dslForm)
	}

	//抓取ip
	ip_from := c.Query("IP")
	ip_to := ip_from
	if len(ip_from) == 0 {
		ip_from = "0.0.0.0"
		ip_to = "255.255.255.255"
	}

	// DSL 查詢條件
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"range": map[string]interface{}{
						"DateTime": map[string]interface{}{
							"gte": t2,
							"lte": t1,
						},
					}},
					{"range": map[string]interface{}{
						"IP": map[string]interface{}{
							"gte": ip_from,
							"lte": ip_to,
						},
					}},
				},
			},
		},
		"size":             1000,
		"from":             0,
		"track_total_hits": true,
	}

	if keyword := c.Query("KeyWord"); len(keyword) != 0 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{"range": map[string]interface{}{
							"DateTime": map[string]interface{}{
								"gte": t2,
								"lte": t1,
							},
						}},
						{"range": map[string]interface{}{
							"IP": map[string]interface{}{
								"gte": ip_from,
								"lte": ip_to,
							},
						}},
						{"term": map[string]interface{}{
							"RawData": keyword,
						},
						},
					},
				},
			},
			"size":             10000,
			"from":             0,
			"track_total_hits": true,
		}
	}

	jsonBody, _ := json.Marshal(query)
	req := esapi.SearchRequest{
		Index: []string{"accesslog"}, // 索引名稱
		Body:  bytes.NewReader(jsonBody),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		c.JSON(http.StatusBadRequest, "查詢失敗，失敗原因："+err.Error())
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		c.JSON(http.StatusBadRequest, "解析失敗，失敗原因："+err.Error())
	}

	var dataList []pojo.Rawdata

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		data := new(pojo.Rawdata)
		source := hit.(map[string]interface{})["_source"]
		data.Time = fmt.Sprintf("%v", source.(map[string]interface{})["@timestamp"])
		data.IP = fmt.Sprintf("%v", source.(map[string]interface{})["IP"])
		data.Message = fmt.Sprintf("%v", source.(map[string]interface{})["RawData"])
		dataList = append(dataList, *data)
	}
	if len(dataList) == 0 {
		c.JSON(http.StatusOK, "Nothing found")
		return
	}
	StructToCsv("account.csv", dataList)
	c.JSON(http.StatusOK, dataList)
}
