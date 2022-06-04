package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Handler(filename string) error {
	// 获取程序所在目录
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	exPath := filepath.Dir(ex)

	filesuffix := filepath.Ext(filename) //获取文件后缀
	outfile := ""                        //输出文件名

	// 1. 根据文件名将文件内容从文件中读出
	rfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	// 2. 读文件
	rinfo, err := rfile.Stat()
	if err != nil {
		return err
	}
	rText := make([]byte, rinfo.Size())
	rfile.Read(rText)
	// fmt.Println("=================input===============")
	// fmt.Println(rText)
	// 3. 关闭文件
	rfile.Close()

	var wText []byte
	if filesuffix == ".xdat" {
		outfile = strings.TrimSuffix(filename, filesuffix)
		wText, err = RSADecrypt(rText, exPath+"/self/private.pem")
	} else {
		outfile = filename + ".xdat"
		wText, err = RSAEncrypt(rText, exPath+"/other/public.pem")
	}
	if err != nil {
		return err
	}
	// fmt.Println("=================output===============")
	// fmt.Println(wText)

	// 1. 根据文件名创建文件并打开
	wfile, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	// 2. 写文件
	wfile.Write(wText)
	// 3. 关闭文件
	wfile.Close()

	fmt.Println("input file: " + filename)
	fmt.Println("output file: " + outfile)
	return nil
}
