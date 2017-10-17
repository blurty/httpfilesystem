// Copyright 2017 Blurt. All rights reserved.
// Use of this source code is governed by a Apache
// license 2.0 that can be found in the LICENSE file.
//
// a http file server for handling files

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var fileDir string = "/Users/steve/tmp/"

type MyMux struct {
}

func IsFileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/query" && r.Method == "GET" {
		query(w, r)
		return
	} else if r.URL.Path == "/upload" && r.Method == "POST" {
		upload(w, r)
		return
	} else if r.URL.Path == "/download" && r.Method == "GET" {
		download(w, r)
		return
	} else if r.URL.Path == "/list" && r.Method == "GET" {
		getFileList(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func getFileList(w http.ResponseWriter, r *http.Request) {
	type filer struct {
		Filename   string    `json:"filename"`
		FileSize   int64     `json:"file_size"`
		ModifyTime time.Time `json:"modify_time"`
	}
	dirList, err := ioutil.ReadDir(fileDir)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	fileList := make([]filer, len(dirList))
	for i, file := range dirList {
		fileList[i].Filename = file.Name()
		fileList[i].ModifyTime = file.ModTime()
		fileList[i].FileSize = file.Size()
	}
	ret, _ := json.Marshal(fileList)
	fmt.Fprint(w, string(ret))
}

func download(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
		return
	}

	filename := r.Form.Get("filename")
	realFile := fileDir + filename
	http.ServeFile(w, r, realFile)
	return
}

func query(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
		return
	}

	filename := r.Form.Get("filename")
	result := true
	if filename == "" {
		result = false
	} else {
		realFile := fileDir + filename
		result, _ = IsFileExists(realFile)
		result = !result
	}
	fmt.Fprintf(w, "%t", result)
}

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	realFile := fileDir + handler.Filename
	result, _ := IsFileExists(realFile)
	if result {
		fmt.Fprintf(w, "processing")
		return
	}

	f, err := os.OpenFile(fileDir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "false")
		return
	}
	defer f.Close()

	if err = f.Truncate(0); err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "false")
	}
	if _, err = io.Copy(f, file); err != nil {
		fmt.Fprintf(w, "false")
	}
	fmt.Fprintf(w, "true")
}

func main() {
	mux := &MyMux{}

	err := http.ListenAndServe(":12345", mux) // 设置监听端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
