/***
	日志类

	path:c:/sdfdsf
	filename:test.text
	wireteString :"sdfsf"

	logers := zloger.NewLog(path)
	logers.DebugLog(filename,wireteString)

 */
package Zloger

import (
	"os"
	"fmt"
	"io"
	"time"
	"sync"
)

type Loger struct{
	path string
}

// 同步锁
var chInfoLock sync.Mutex
var chCompleteLock sync.Mutex
var chErrorLock sync.Mutex


func NewLog(path string) *Loger {

	//创建日志目录
	err := os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Printf("crete fail %s", err)
		os.Exit(1)
	} else {
		// fmt.Print("Create Directory OK!")
	}

	//返回日志对象
	Log := &Loger{
		path : path+"/",
	}
	return Log
}

func (a *Loger) DebugLog(wireteString string){
	var f    *os.File
	var err   error

	var filename = "debug"+getCurDate()+".log"

	if checkFileIsExist(a.path+filename) {  //如果文件存在
		f, err = os.OpenFile(a.path+filename, os.O_APPEND,0)  //打开文件
		//fmt.Println("文件存在");
	}else {
		f, err = os.Create(a.path+filename)  //创建文件
		//fmt.Println("文件不存在");
	}
	check(err)
	defer f.Close()

	wireteString = wireteString+"\r\n"
	_,err = io.WriteString(f, wireteString) //写入文件(字符串)
	check(err)
	//格式化用的日期是特定的，123（15）45 -.-lll
	now_time := time.Now().Format("2006-01-02 15:04:05")
	io.WriteString(f, now_time+"\r\n") //写入文件(字符串)
	//fmt.Printf("写入 %d 个字节", n);
	return
}

func (a *Loger) CompleteLog(wireteString string){
	chCompleteLock.Lock()
	defer chCompleteLock.Unlock()
	var f    *os.File
	var err   error

	var filename = "complete_"+getCurDate()+".log"

	if checkFileIsExist(a.path+filename) {  //如果文件存在
		f, err = os.OpenFile(a.path+filename,  os.O_APPEND|os.O_WRONLY, os.ModeAppend)  //打开文件
		//fmt.Println("文件存在");
	}else {
		f, err = os.Create(a.path+filename)  //创建文件
		//fmt.Println("文件不存在");
	}


	check(err)

	wireteString = wireteString+"\r\n"
	_,err = io.WriteString(f, wireteString) //写入文件(字符串)
	check(err)
	//格式化用的日期是特定的，123（15）45 -.-lll
	now_time := time.Now().Format("2006-01-02 15:04:05")
	io.WriteString(f, now_time+"\r\n") //写入文件(字符串)
	//fmt.Printf("写入 %d 个字节", n);
	return
}

func (a *Loger) InfoLog(wireteString string,extfilename string){


	var f    *os.File
	var err   error

	var filename = "info_"+extfilename + getCurDate()+".log"

	if checkFileIsExist(a.path+filename) {  //如果文件存在
		f, err = os.OpenFile(a.path+filename,  os.O_APPEND|os.O_WRONLY, os.ModeAppend)  //打开文件
		//fmt.Println("文件存在");
	}else {
		f, err = os.Create(a.path+filename)  //创建文件
		//fmt.Println("文件不存在");
	}


	check(err)

	wireteString = wireteString+"\r\n"
	_,err = io.WriteString(f, wireteString) //写入文件(字符串)
	check(err)
	//格式化用的日期是特定的，123（15）45 -.-lll
	now_time := time.Now().Format("2006-01-02 15:04:05")
	io.WriteString(f, now_time+"\r\n") //写入文件(字符串)
	//fmt.Printf("写入 %d 个字节", n);


	return
}


func (a *Loger) ErrorLog(wireteString string){
	chErrorLock.Lock()
	defer chErrorLock.Unlock()
	var f    *os.File
	var err   error

	var filename = "error_"+getCurDate()+".log"

	if checkFileIsExist(a.path+filename) {  //如果文件存在
		f, err = os.OpenFile(a.path+filename,  os.O_APPEND|os.O_WRONLY, os.ModeAppend)  //打开文件
		//fmt.Println("文件存在");
	}else {
		f, err = os.Create(a.path+filename)  //创建文件
		//fmt.Println("文件不存在");
	}


	check(err)

	wireteString = wireteString+"\r\n"
	_,err = io.WriteString(f, wireteString) //写入文件(字符串)
	check(err)
	//格式化用的日期是特定的，123（15）45 -.-lll
	now_time := time.Now().Format("2006-01-02 15:04:05")
	io.WriteString(f, now_time+"\r\n") //写入文件(字符串)
	//fmt.Printf("写入 %d 个字节", n);
	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
/**
* 判断文件是否存在  存在返回 true 不存在返回false
*/
func checkFileIsExist(path string) (bool) {
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




//2006-1-2
func getCurDate() string  {
	return time.Now().Format("2006-1-2")
}