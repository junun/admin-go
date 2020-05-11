package util

import (
	"api/pkg/setting"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func GetIdRsaPath() string {
	return setting.AppSetting.IdRsaPath
}

func GenerateKey()  {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Private key cannot be created.", err.Error())
	}

	publickey := &key.PublicKey

	dir, _ := os.Getwd()
	path := dir + "/" + GetIdRsaPath()

	// 检查 保存公钥和私钥目录是否存在
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// dump private key to file
	err = DumpPrivateKeyFile(key, path + "id_rsa")
	if err != nil {
		log.Fatalf("Cannot dump private key file\n")
	}

	// dump public key to file
	err = DumpPublicKeyFile(publickey, path + "id_rsa_pub")
	if err != nil {
		log.Fatalf("Cannot dump public key file\n")
	}
}

// Dump private key into file
func DumpPrivateKeyFile(privatekey *rsa.PrivateKey, filename string) error {
	var keybytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	block := &pem.Block{
		Type  : "RSA PRIVATE KEY",
		Bytes :  keybytes,
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

// Dump public key into file
func DumpPublicKeyFile(publickey *rsa.PublicKey, filename string) error {
	keybytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type  : "PUBLIC KEY",
		Bytes :  keybytes,
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

// Dump private key to buffer.
func DumpPrivateKeyBuffer(privatekey *rsa.PrivateKey) (string, error) {
	var keybytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	block := &pem.Block{
		Type  : "RSA PRIVATE KEY",
		Bytes :  keybytes,
	}

	var keybuffer []byte = pem.EncodeToMemory(block)
	return string(keybuffer), nil
}

func DumpPublicKeyBuffer(publickey *rsa.PublicKey) (string, error) {
	keybytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return "", err
	}

	block := &pem.Block{
		Type  : "PUBLIC KEY",
		Bytes :  keybytes,
	}

	var keybuffer []byte = pem.EncodeToMemory(block)
	return string(keybuffer), nil
}

// Dump private key to base64 string
// Compared with DumpPrivateKeyBuffer this output:
//  1. Have no header/tailer line
//  2. Key content is merged into one-line format
func DumpPrivateKeyBase64(privatekey *rsa.PrivateKey) (string, error) {
	var keybytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)

	keybase64 := base64.StdEncoding.EncodeToString(keybytes)
	return keybase64, nil
}

func DumpPublicKeyBase64(publickey *rsa.PublicKey) (string, error) {
	keybytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return "", err
	}

	keybase64 := base64.StdEncoding.EncodeToString(keybytes)
	return keybase64, nil
}

// Load private key from private key file
func LoadPrivateKeyFile(keyfile string) (*rsa.PrivateKey, error) {
	keybuffer, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(keybuffer))
	if block == nil {
		return nil, errors.New("private key error!")
	}

	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("parse private key error!")
	}

	return privatekey, nil
}


func LoadPublicKeyFile(keyfile string) (*rsa.PublicKey, error) {
	keybuffer, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keybuffer)
	if block == nil {
		return nil, errors.New("public key error")
	}

	pubkeyinterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publickey := pubkeyinterface.(*rsa.PublicKey)
	return publickey, nil
}

func LoadPublicKeyFileToAuthorizedFormat(keyfile string) (string, error) {
	publickey, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(publickey)
	if block == nil {
		panic("invalid public key data")
	}
	if block.Type != "PUBLIC KEY" {
		panic(fmt.Sprintf("invalid public key type : %s", block.Type))
		return "", err
	}

	keyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("failed to convert bytes to public key.")
	}

	key, ok := keyInterface.(*rsa.PublicKey)
	if !ok {
		panic("not RSA public key")
	}

	skey, err := ssh.NewPublicKey(key)
	if (err != nil) {
		panic("failed to convert ras key to ssh key.")
	}

	b := string(ssh.MarshalAuthorizedKey(skey))

	b = strings.TrimSuffix(b, "\n")
	return b, nil
}

// Load private key from base64
func LoadPrivateKeyBase64(base64key string) (*rsa.PrivateKey, error) {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	privatekey, err := x509.ParsePKCS1PrivateKey(keybytes)
	if err != nil {
		return nil, errors.New("parse private key error!")
	}

	return privatekey, nil
}


func LoadPublicKeyBase64(base64key string) (*rsa.PublicKey, error) {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	pubkeyinterface, err := x509.ParsePKIXPublicKey(keybytes)
	if err != nil {
		return nil, err
	}

	publickey := pubkeyinterface.(*rsa.PublicKey)
	return publickey, nil
}

// encrypt
func Encrypt(plaintext string, publickey *rsa.PublicKey) (string, error) {
	label := []byte("")
	sha256hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(sha256hash, rand.Reader, publickey, []byte(plaintext), label)

	decodedtext := base64.StdEncoding.EncodeToString(ciphertext)
	return decodedtext, err
}

// decrypt
func Decrypt(ciphertext string, privatekey *rsa.PrivateKey) (string, error) {
	decodedtext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed, error=%s\n", err.Error())
	}

	sha256hash := sha256.New()
	decryptedtext, err := rsa.DecryptOAEP(sha256hash, rand.Reader, privatekey, decodedtext, nil)
	if err != nil {
		return "", fmt.Errorf("RSA decrypt failed, error=%s\n", err.Error())
	}

	return string(decryptedtext), nil
}

