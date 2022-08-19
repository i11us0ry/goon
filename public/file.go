package public

import (
	"bufio"
	"github.com/axgle/mahonia"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//文件读取，返回文件流
func FileRead(fileName string) *bufio.Scanner{
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		Error.Printf("Read file:%s failed!\n",fileName)
		os.Exit(0)
	}
	datas := bufio.NewScanner(f)
	return datas
}

// 文件读取，将读取内容按行保存到数组返回
func FileReadByline(fileName string) []string{
	datas := []string{}
	fileData := FileRead(fileName)
	for fileData.Scan(){
		datas = append(datas, fileData.Text())
	}
	return datas
}

func FileWrite(fileName string, writeInfo string) {
	enc:=mahonia.NewEncoder("utf-8")
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777) //打开文件
	if err != nil {
		Error.Printf("%v\n",err)
		os.Exit(0)
	}
	if _, err = io.WriteString(f, enc.ConvertString(writeInfo)); err != nil {
		Error.Printf("%v\n",err)
		os.Exit(0)
	}
	defer f.Close()
}

/* 获取程序执行路径 */
func GetCurrentDir() string{
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error.Printf("%v\n",err)
		os.Exit(0)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

/*
文件是否存在
*/
func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return false
	}
	return true
}

// 生成16位token
func GetRandom() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
