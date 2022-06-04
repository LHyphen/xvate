package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// 参数bits: 指定生成的秘钥的长度, 单位: bit
func RsaGenKey(bits int) error {
	// 获取程序所在目录
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	selfPath := exPath + "/self"
	otherPath := exPath + "/other"

	// 创建存放密钥的目录
	os.Mkdir(selfPath, os.ModePerm)
	os.Mkdir(otherPath, os.ModePerm)

	// 1. 生成私钥文件
	// GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	// 参数1: Reader是一个全局、共享的密码用强随机数生成器
	// 参数2: 秘钥的位数 - bit
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	// 2. MarshalPKCS1PrivateKey将rsa私钥序列化为ASN.1 PKCS#1 DER编码
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	// 3. Block代表PEM编码的结构, 对其进行设置
	block := pem.Block{
		Type:  "RSA PRIVATE KEY", //"RSA PRIVATE KEY",
		Bytes: derStream,
	}
	// 4. 创建文件
	privFile, err := os.Create(selfPath + "/private.pem")
	if err != nil {
		return err
	}
	// 5. 使用pem编码, 并将数据写入文件中
	err = pem.Encode(privFile, &block)
	if err != nil {
		return err
	}
	// 6. 最后的时候关闭文件
	defer privFile.Close()

	// 7. 生成公钥文件
	publicKey := privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}
	block = pem.Block{
		Type:  "RSA PUBLIC KEY", //"PUBLIC KEY",
		Bytes: derPkix,
	}
	pubFile, err := os.Create(selfPath + "/public_" + enterName() + ".pem")
	if err != nil {
		return err
	}
	// 8. 编码公钥, 写入文件
	err = pem.Encode(pubFile, &block)
	if err != nil {
		panic(err)
	}
	defer pubFile.Close()

	return nil

}

func enterName() (name string) {
	fmt.Print("please input a name: ")
	fmt.Scanln(&name)
	r, _ := regexp.Compile("^[a-zA-Z0-9_-]{3,16}$")
	for !r.MatchString(name) {
		fmt.Println("not legal name, only support letters, numbers, _- and 3~16 bits")
		fmt.Print("please enter again: ")
		fmt.Scanln(&name)
	}
	return
}

// func main() {
// 	RsaGenKey(4096)
// }
