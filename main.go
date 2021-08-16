package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func getWhiteList() []string {
	req, err := http.NewRequest("GET", `https://hr-link.atlassian.net/rest/api/latest/search?jql=project="HRL"%20AND%20status%20in%20("In%20Testing","Ready%20for%20Testing")`, nil)
	if err != nil {
		return nil
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("abalan@hr-link.ru", "5kpsOrriEo3YhBXjRQ5b736C")
	cl := http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	t := struct {
		Issues []struct {
			Key string `json:"key"`
		} `json:"issues"`
	}{}
	json.NewDecoder(resp.Body).Decode(&t)
	list := []string{}
	for _, i := range t.Issues {
		list = append(list, i.Key)
	}
	return list

}

func main() {

	logPath := "/var/hr-link/rm.log"
	logPath = "./rm.log"
	workPath := ""
	workPath = "/home/andew/test_folder/"

	file, _ := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	file.WriteString(fmt.Sprintf("%s %s", time.Now().Format("\n2006-01-02 15:04:05"), "start"))
	file.Close()

	check := func(name string, list []string) bool {
		list = append(list, "hotfix", "release", "master", "test-merge")
		for _, l := range list {
			for name == l {
				return false
			}
		}
		return true
	}

	if _, err := os.Stat(logPath); err == os.ErrNotExist {
		f, _ := os.Create(logPath)
		f.Chmod(0777)
		f.Close()
	}
	for {
		if time.Now().Hour() == 12 {
			fmt.Println(1)
			list := getWhiteList()
			dir, _ := os.ReadDir(workPath)
			file, _ := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)

			for _, f := range dir {
				if check(f.Name(), list) {
					os.RemoveAll(workPath + f.Name())
					file.WriteString(fmt.Sprintf("%s %s", time.Now().Format("\n2006-01-02 15:04:05"), f.Name()))
				}
			}
			file.Close()
		} else {
			fmt.Println(2)
			time.Sleep(time.Hour * 1)
		}
	}
}
