package service

import (
	"fmt"
	"os"
)

func ParseMP3(path string) {
	//fmt.Println("path:", path)
	//files, err := os.Open(path)
	//if err != nil {
	//	fmt.Println("error opening directory:", err)
	//	return
	//}
	//defer files.Close()
	//fileInfos, err := files.ReadDir(-1)
	//if err != nil {
	//	fmt.Println("error reading directory:", err)
	//	return
	//}
	//for _, fileInfo := range fileInfos {
	//	fmt.Println(fileInfo.Name())
	//}

	fileInfos, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("error reading directory:", err)
		return
	}
	for _, fileInfo := range fileInfos {
		fmt.Println(fileInfo.Name())
	}
}
