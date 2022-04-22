package pojo

/*
定義抓出來資料的struct
*/

type Rawdata struct {
	Time    string `json:"時間"`
	IP      string `json:"IP"`
	Message string `json:"RawData"`
}
