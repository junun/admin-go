package util

import (
	"api/pkg/setting"
	"os"
	"strconv"
)

func GetSyncPath() string {
	return setting.AppSetting.SyncPath
}

func ReturnSyncRunDir(env string, appId int) string {
	dir, _ 		:= os.Getwd()
	path 		:= dir + "/" + GetSyncPath() + env + "/" + strconv.Itoa(appId)
	_, err 		:= os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	return  path
}