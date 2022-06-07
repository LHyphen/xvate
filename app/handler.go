package app

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Handler(filename string) error {
	// 计时
	t1 := time.Now() // get current time
	defer func() {
		// 如果 defer 语句调用的是一个匿名函数，
		// 那么匿名函数体内不管有什么复杂的逻辑，统统在 defer 外围函数即将退出时才执行。
		elapsed := time.Since(t1)
		fmt.Println("App elapsed: ", elapsed)
	}()

	// 获取程序所在目录
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	exPath := filepath.Dir(ex)

	//获取文件后缀
	filesuffix := filepath.Ext(filename)

	//确定输出文件名
	outfile := ""
	if filesuffix == ".xdat" {
		outfile = strings.TrimSuffix(filename, filesuffix)
	} else {
		outfile = filename + ".xdat"
	}

	// 打开输入文件
	rfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer rfile.Close()

	// 打开输出文件
	wfile, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer wfile.Close()

	// 打开privateKey
	privateKey, err := openPrivateKey(exPath + "/self/private.pem")
	if err != nil {
		return err
	}

	// 打开pubKey
	// 选择要加密公钥
	pubkeys := scanDir(exPath + "/other")
	username, err := selectUser(pubkeys)
	if err != nil {
		return err
	}
	pubKey, err := openPublicKey(exPath + "/other/" + username + ".pem")
	if err != nil {
		return err
	}

	// 确定单次操作的块大小
	var block int
	if filesuffix == ".xdat" {
		block = privateKey.Size()
	} else {
		block = pubKey.Size() - 11
	}

	// //sync pools to reuse the memory and decrease the preassure on //Garbage Collector
	// blockPool := sync.Pool{
	// 	New: func() interface{} {
	// 		lines := make([]byte, block)
	// 		return &lines
	// 	},
	// }

	rText := bufio.NewReader(rfile)
	for {
		// buf := blockPool.Get().(*[]byte) //the chunk size
		buf := make([]byte, block)
		n, err := rText.Read(buf) //loading chunk into buffer
		buf = buf[:n]
		if n == 0 {
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
				break
			}
			return err
		}

		// fmt.Println("=================input===============", "block: ", block, ",n: ", n)
		// fmt.Println(buf)

		var bytesOnce []byte
		// 私钥解密 公钥加密
		if filesuffix == ".xdat" {
			bytesOnce, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, buf)
		} else {
			bytesOnce, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, buf)
		}
		if err != nil {
			return err
		}
		wfile.Write(bytesOnce)
		// blockPool.Put(buf)
	}

	fmt.Println("\ninput file: " + filename)
	fmt.Println("output file: " + outfile)
	return nil
}

func scanDir(pubPath string) (files []string) {
	err := filepath.Walk(pubPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".pem" {
			return nil
		}
		files = append(files, strings.TrimSuffix(info.Name(), ".pem"))
		return nil
	})
	if err != nil {
		panic(err)
	}
	return
}

func selectUser(files []string) (username string, err error) {
	n := len(files)
	if n == 0 {
		return "", errors.New("not found public key, please check /other dir")
	} else if n == 1 {
		username = files[0]
		return
	}

	fmt.Println("\n==================ALL USER==================")
	for i, file := range files {
		fmt.Printf("[%2d]. %v\n", i, file)
	}

	num := -1
	fmt.Print("please select a user's pubkey(enter a number): ")
	fmt.Scanln(&num)
	for num == -1 || num < 0 || num > n-1 {
		fmt.Print("input is wrong, please input again: ")
		fmt.Scanln(&num)
	}

	username = files[num]

	return
}

func openPrivateKey(filename string) (*rsa.PrivateKey, error) {
	// 1. 根据文件名将文件内容从文件中读出
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// 2. 读文件
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	allText := make([]byte, info.Size())
	file.Read(allText)
	// 3. 关闭文件
	file.Close()
	// 4. 从数据中查找到下一个PEM格式的块
	block, _ := pem.Decode(allText)
	// 5. 解析一个pem格式的私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	return privateKey, err
}

func openPublicKey(filename string) (*rsa.PublicKey, error) {
	// 1. 根据文件名将文件内容从文件中读出
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// 2. 读文件
	info, err := file.Stat()
	allText := make([]byte, info.Size())
	file.Read(allText)
	// 3. 关闭文件
	file.Close()

	// 4. 从数据中查找到下一个PEM格式的块
	block, _ := pem.Decode(allText)
	if block == nil {
		return nil, err
	}
	// 5. 解析一个DER编码的公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubKey := pubInterface.(*rsa.PublicKey)
	return pubKey, nil
}
