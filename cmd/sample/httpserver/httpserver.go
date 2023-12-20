package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "net/http/pprof"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	omap := make(map[string]interface{})
	omap["header"] = r.Header
	omap["body"] = string(reqBody)

	data, _ := json.Marshal(omap)
	fmt.Println(string(data))
	fmt.Fprintln(w, string(data))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please input port")
		return
	}

	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Args[1]), nil)
}
