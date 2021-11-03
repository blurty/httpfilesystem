package main
import (
        "io"
	"fmt"
        "net/http"
		"log"
		"encoding/json"
)
type Data struct{
	Name string
	Age int
}
type Ret struct{
	ErrNo int `json:"status"`
	ErrMsg string `json:"msg"`
	//Data []Data `json:"data"`
	Data map[string]interface{} `json:"data"`
}
func HelloServer(w http.ResponseWriter, req *http.Request) {
	data := Data{Name: "why", Age: 18}

	ret := new(Ret)
	id := req.FormValue("id")
	fmt.Println("id:", id)
	oaid := req.FormValue("oaid")
	imei := req.FormValue("imei_md5")
	idfa := req.FormValue("idfa_md5")
	cuid := req.FormValue("cuid")
	device := req.FormValue("device_id")
	fmt.Println("oaid", oaid, "imei:", imei, "idfa:", idfa, "device_id:", device,  "cuid:", cuid)
	dateRange := req.FormValue("date_range")
	fmt.Println("date_range:", dateRange)
	//id := req.PostFormValue('id')

	ret.ErrNo = 0
	ret.ErrMsg = "success"
	retData := make([]Data, 0)
	retData = append(retData, data)
	retData = append(retData, data)
	retData = append(retData, data)
	ret.Data = make(map[string]interface{})
	ret.Data["rows"] = retData
	ret_json,_ := json.Marshal(ret)

	io.WriteString(w, string(ret_json))
}
func HelloServer1(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world1!\n")
}

func main() {
		http.HandleFunc("/hello", HelloServer)
		http.HandleFunc("/", HelloServer)
        err := http.ListenAndServe(":8888", nil)
        if err != nil {
                log.Fatal("ListenAndServe: ", err)
        }
}
