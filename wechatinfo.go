package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"golang.org/x/sys/windows/registry"
)

var (
	info     [][]byte
	filePath string
	wxid     []string
)

func main() {
	username, _ := user.Current()
	userdir := username.HomeDir

	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Tencent\WeChat`, registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	s, _, err := k.GetStringValue("FileSavePath")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(s)
	if s == "MyDocument:" {
		filePath = userdir + `\Documents\WeChat Files\`
	} else {
		filePath = s
	}
	fileList, err := ioutil.ReadDir(filePath)
	if err != nil {
		log.Fatal("read dir error")
	}

	for _, v := range fileList {
		if v.Name() != "All Users" && v.Name() != "Applet" {
			wxid = append(wxid, v.Name())
		}
	}

	for _, v := range wxid {
		fmt.Println("===========================================")
		path := filePath + v + `\config\AccInfo.dat`
		content, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		info = getInfo(content)
		for _, v := range info {
			fmt.Println(string(v))
		}
	}
	// fmt.Println(wxid)
	// content, err := os.ReadFile("AccInfo")
	// if err != nil {
	// 	panic(err)
	// }
	// info = getInfo(content)
	// for _, v := range info {
	// 	fmt.Println(string(v))
	// }
}

func getInfo(content []byte) (out [][]byte) {
	var raw []byte
	for i := 0; i < len(content); i++ {
		if content[i] == 8 && content[i+1] == 4 && content[i+2] == 18 && content[i+3] == 19 {
			raw = content[i:]
			break
		}
	}
	for h, t := 0, 0; h < len(raw); {
		if raw[h] == 8 && raw[h+2] == 18 {
			t = h + 4
			for ; raw[t] != 26; t++ {
			}
			out = append(out, raw[h+4:t])
			h = t + 1
		} else {
			h++
		}
	}
	return
}
