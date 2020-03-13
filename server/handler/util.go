package handler

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
)

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func MkDirAll(path string){
	if !Exists(path) {
		err := os.MkdirAll(path, 0766) // default will be effected by umask
		if err != nil {
			log.Fatal(err)
			return
		}
		if err := os.Chmod(path, 0766); err != nil {
			log.Fatal(err)
		}
	}
}