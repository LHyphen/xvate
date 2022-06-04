package app

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

//RSA公钥加密
func RSAEncrypt(src []byte, filename string) ([]byte, error) {
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

	// // 6. 公钥加密
	// result, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, src)
	// return result, err

	// 6. 分段公钥加密
	keySize, srcSize := pubKey.Size(), len(src)
	// log.Println("密钥长度：", keySize, "\t明文长度：\t", srcSize)
	//单次加密的长度需要减掉padding的长度，PKCS1为11
	offSet, once := 0, keySize-11
	buffer := bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + once
		if endIndex > srcSize {
			endIndex = srcSize
		}
		// 加密一部分
		bytesOnce, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, src[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	bytesEncrypt := buffer.Bytes()
	return bytesEncrypt, err
}

//RSA私钥解密
func RSADecrypt(src []byte, filename string) ([]byte, error) {
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
	if err != nil {
		return nil, err
	}
	// // 6. 私钥解密
	// result, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
	// if err != nil {
	// 	return nil, err
	// }

	// 6. 分段公钥解密
	keySize := privateKey.Size()
	srcSize := len(src)
	// log.Println("密钥长度：", keySize, "\t密文长度：\t", srcSize)
	var offSet = 0
	var buffer = bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + keySize
		if endIndex > srcSize {
			endIndex = srcSize
		}
		bytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, src[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	bytesDecrypt := buffer.Bytes()

	return bytesDecrypt, nil
}

// func main() {
// 	//RsaGenKey(4096)
// 	src := []byte("我是小庄, 如果我死了, 肯定不是自杀...")
// 	cipherText := RSAEncrypt(src, "public.pem")
// 	fmt.Println(cipherText)
// 	plainText := RSADecrypt(cipherText, "private.pem")
// 	fmt.Println(string(plainText))
// }
