// Copyright 2017 Blurt. All rights reserved.
// Use of this source code is governed by a Apache
// license 2.0 that can be found in the LICENSE file.
//
// a http file client for handling files

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func postFile(filename, fileDir, targetUrl string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		return "false", err
	}

	//打开文件句柄操作
	realFile := path.Join(fileDir, filename)
	fh, err := os.Open(realFile)
	if err != nil {
		return "false", err
	}
	defer fh.Close()

	//ioCopy
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
		return fmt.Sprint("http status code:", resp.StatusCode), nil
	} else {
		return string(body), nil
	}
}

func downloadFile(filename, fileDir, targetUrl string) (string, error) {
	resp, err := http.Get(targetUrl + "?filename=" + filename)
	if err != nil {
		return "false", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "false", err
	} else {
		realFile := path.Join(fileDir, filename)
		f, err := os.OpenFile(realFile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return "false", err
		}
		defer f.Close()
		io.Copy(f, resp.Body)
		return "true", nil
	}
}

func getFileList(targetUrl string) (body []byte, err error) {
	resp, err := http.Get(targetUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("http status code not 200")
		return
	} else {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		} else {
			return
		}
	}
}

// sample usage
func main() {
	uploadFilename := flag.String("u", "", "upload file to server")
	downloadFilename := flag.String("d", "", "download file from server")
	queryFilename := flag.String("q", "", "result of file transfer")
	fileDir := flag.String("dir", "", "dir of file")
	serverIP := flag.String("h", "", "refer server ip")
	listFlag := flag.Bool("l", false, "list all files on server")
	flag.Usage = usage

	flag.Parse()

	if *serverIP == "" {
		fmt.Println("missing server ip")
		flag.Usage()
		return
	}

	upload_url := "http://" + *serverIP + ":12345/upload"
	query_url := "http://" + *serverIP + ":12345/query"
	download_url := "http://" + *serverIP + ":12345/download"
	list_url := "http://" + *serverIP + ":12345/list"

	if *uploadFilename != "" {
		result, err := postFile(*uploadFilename, *fileDir, upload_url)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
		return
	} else if *downloadFilename != "" {
		result, err := downloadFile(*downloadFilename, *fileDir, download_url)
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
	} else if *listFlag {
		fileList, err := getFileList(list_url)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(fileList))
		}
	} else {
		flag.Usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `httpfileclient version: httpfileclient/3.0
Usage: ./httpfileclient [-h server] [-u filename] [-d filename] [-q filename] [-dir filepath] [-l]

Options:
`)
	flag.PrintDefaults()
}
