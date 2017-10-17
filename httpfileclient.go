// Copyright 2017 Blurt. All rights reserved.
// Use of this source code is governed by a Apache
// license 2.0 that can be found in the LICENSE file.
//
// a http file client for handling files

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func postFile(filename string, targetUrl string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		return "false", err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		return "false", err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return "false", err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return "false", err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "false", err
	}
	if resp.StatusCode != 200 {
		return "false", err
	} else {
		return string(resp_body), nil
	}
}

func queryFileTransferProgress(filename string, queryUrl string) (string, error) {
	resp, err := http.Get(queryUrl + "?filename=" + filename)
	if err != nil {
		return "false", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "false", err
	}

	if resp.StatusCode != 200 {
		fmt.Println("code:", resp.StatusCode)
		return "false", nil
	} else {
		return string(body), nil
	}
}

func downloadFile(filename string, targetUrl string) (string, error) {
	resp, err := http.Get(targetUrl + "?filename=" + filename)
	if err != nil {
		return "false", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "false", err
	} else {
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return "false", err
		}
		defer f.Close()
		io.Copy(f, resp.Body)
		return filename + "download success", nil
	}
}

// sample usage
func main() {
	cmd := `Usage:./httpfileclient -h 127.0.0.1 <pattern> <filename>
		-h server ip
		-u upload file
		-d download file
		-q result of file transfer`
	uploadFilename := flag.String("u", "", "upload filename")
	downloadFilename := flag.String("d", "", "download filename")
	queryFilename := flag.String("q", "", "query filename")
	serverIP := flag.String("h", "", "server ip")

	flag.Parse()

	if *serverIP == "" {
		fmt.Println("missing server ip")
		fmt.Println(cmd)
		return
	}
	if *uploadFilename == "" && *downloadFilename == "" && *queryFilename == "" {
		fmt.Println("choose one of upload or download or query")
		fmt.Println(cmd)
		return
	}

	upload_url := "http://" + *serverIP + ":12345/upload"
	query_url := "http://" + *serverIP + ":12345/query"
	download_url := "http://" + *serverIP + ":12345/download"

	if *uploadFilename != "" {
		result, err := postFile(*uploadFilename, upload_url)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
		return
	} else if *downloadFilename != "" {
		result, err := downloadFile(*downloadFilename, download_url)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
		return
	} else if *queryFilename != "" {
		result, err := queryFileTransferProgress(*queryFilename, query_url)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
		return
	}
}
