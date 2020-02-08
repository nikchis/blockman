// Copyright (c) 2018 Nikita Chisnikov
// Distributed under the MIT/X11 software license

package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
)

const C_DEF_STORAGE_SALT = "s0me.go0d.s@lt"

func checksum(path string) (string, error) {
	var buf bytes.Buffer
	var result string
	file, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer file.Close()
	if _, err := io.Copy(&buf, file); err != nil {
		return result, err
	}
	h := make([]byte, 64)
	sha3.ShakeSum256(h, buf.Bytes())
	buf.Reset()
	result = hex.EncodeToString(h)
	h = nil
	return result, nil
}

func genPwd(phrase string) string {
	info, err := cpu.Info()
	if err != nil {
		log.Fatal("cpu.Info():", err.Error())
	}
	for i := range info {
		js, err := json.Marshal(info[i])
		if err != nil {
			log.Fatal("json.Marshal()", err.Error())
		}
		phrase += string(js)
		js = nil
	}
	info = nil
	phash := make([]byte, 32)
	sha3.ShakeSum128(phash, []byte(phrase))
	result := hex.EncodeToString(phash)
	phash = nil
	phrase = ""
	return result
}

func genSalt() []byte {
	res := strings.Join(getMacs(), ":") + C_DEF_STORAGE_SALT
	if res == "" {
		log.Panic("generated salt is null!")
	}
	salt := make([]byte, 32)
	sha3.ShakeSum128(salt, []byte(res))
	salt = nil
	res = ""
	return salt
}

func getKey(pass string) []byte {
	key, err := scrypt.Key([]byte(pass), genSalt(), 32768, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

func encrypt(data []byte, pswd string) []byte {
	block, _ := aes.NewCipher(getKey(pswd))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal("cipher.NewGCM():", err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal("io.ReadFull():", err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	nonce = nil
	block = nil
	gcm = nil
	return ciphertext
}

func decrypt(data []byte, pswd string) ([]byte, error) {
	block, err := aes.NewCipher(getKey(pswd))
	if err != nil {
		log.Fatal("aes.NewCipher():", err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal("cipher.NewGCM():", err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	nonce = nil
	block = nil
	gcm = nil
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return plaintext, nil
}
