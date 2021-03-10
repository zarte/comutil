// Copyright 2018 ouqiang authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package comutil

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unsafe"
)


// MD5 生成MD5摘要
func MD5(s string) string {
	/*
	m := md5.New()
	m.Write([]byte(s))

	return hex.EncodeToString(m.Sum(nil))
	*/

	data := []byte(s)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return  md5str1
}

// MD5 byte生成MD5摘要
func MD5bt(s []byte) string {
	m := md5.New()
	m.Write(s)
	return hex.EncodeToString(m.Sum(nil))
}

// RandNumber 生成min - max之间的随机数
func RandNumber(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return min + r.Intn(max-min)
}

// PanicToError Panic转换为error
func PanicToError(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic: %s", e)
		}
	}()
	f()
	return
}

// PrintAppVersion 打印应用版本
func PrintAppVersion(appVersion, GitCommit, BuildDate string) {
	versionInfo, err := FormatAppVersion(appVersion, GitCommit, BuildDate)
	if err != nil {
		panic(err)
	}
	fmt.Println(versionInfo)
}

// FormatAppVersion 格式化应用版本信息
func FormatAppVersion(appVersion, GitCommit, BuildDate string) (string, error) {
	content := `
   Version: {{.Version}}
Go Version: {{.GoVersion}}
Git Commit: {{.GitCommit}}
     Built: {{.BuildDate}}
   OS/ARCH: {{.GOOS}}/{{.GOARCH}}
`
	tpl, err := template.New("version").Parse(content)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, map[string]string{
		"Version":   appVersion,
		"GoVersion": runtime.Version(),
		"GitCommit": GitCommit,
		"BuildDate": BuildDate,
		"GOOS":      runtime.GOOS,
		"GOARCH":    runtime.GOARCH,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), err
}

// DownloadFile 文件下载
func DownloadFile(filePath string, rw http.ResponseWriter) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	filename := path.Base(filePath)
	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	_, err = io.Copy(rw, file)

	return err
}

// WorkDir 获取程序运行时根目录
func WorkDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	wd := filepath.Dir(execPath)
	if filepath.Base(wd) == "bin" {
		wd = filepath.Dir(wd)
	}

	return wd, nil
}


func Hmac256(content string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(content))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func Checkexist(path string) bool{
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	fmt.Println(err)
	return false
}


func Curlv1(ourl string,data map[string]string, header map[string]string,dtype string) (string,error) {

	reader := strings.NewReader("")

	var mentype string
	if dtype =="JSON" {
		mentype = "POST"
	}else{
		mentype = dtype
	}
	if data!=nil {
		if mentype == "POST"{
			str, err := json.Marshal(data)
			if err != nil {
				fmt.Println("json.Marshal failed:", err)
				return "",err
			}
			reader = strings.NewReader(string(str))
		} else if mentype== "GET"{
			params := url.Values{}
			parseURL, err := url.Parse(ourl)
			if err != nil {
				log.Println("err")
			}
			for key,val := range data {
				params.Set(key, val)
			}
			//如果参数中有中文参数,这个方法会进行URLEncode
			parseURL.RawQuery = params.Encode()
			ourl = parseURL.String()
		}

	}
	//fmt.Println(ourl)
	// request, err := http.Get(ourl)
	//fmt.Println(reader)
	request, err := http.NewRequest(mentype, ourl, reader)
	if err != nil {
		return "",err
	}
	if dtype =="JSON" {
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
	}


	for key,item := range header{
		request.Header.Set(key,item)
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "",err
	}
	//utf8Reader := transform.NewReader(resp.Body,simplifiedchinese.GBK.NewDecoder())
	respBytes, err := ioutil.ReadAll(resp.Body)
	//respBytes, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		return "",err
	}
	//byte数组直接转成string，优化内存

	//utf8 := mahonia.NewDecoder("utf8").ConvertString(string(respBytes))
	//ioutil.WriteFile("./output2.txt", respBytes, 0666) //写入文件(字节数组)
	res := (*string)(unsafe.Pointer(&respBytes))
	return *res,nil
}