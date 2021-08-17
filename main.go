package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/andy-smoker/clerk"
)

func getWhiteList(login, pass, project string) []string {
	req, err := http.NewRequest("GET", `https://`+project+`.atlassian.net/rest/api/latest/search?jql=project="HRL"%20AND%20status%20in%20("In%20Testing","Ready%20for%20Testing")`, nil)
	if err != nil {
		return nil
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(login, pass)
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

type CFG struct {
	Logpath string   `toml:"log"`
	WorkDir string   `toml:"work"`
	Ignore  []string `toml:"ignore"`
	Login   string   `toml:"login"`
	Pass    string   `toml:"pass"`
	Hour    int      `toml:"hour"`
	Project string   `toml:"project"`
}

func getConfig() *CFG {
	cfg := CFG{
		Logpath: "./rm.log",
		WorkDir: "./",
		Login:   "",
		Pass:    "",
		Project: "",
	}
	toml.DecodeFile("config.toml", &cfg)
	return &cfg
}

func main() {
	cfg := getConfig()
	p := clerk.NewPrinter("INFO", "gongfarmer", cfg.Logpath)
	p.WriteLog(1, time.Now(), "start")
	check := func(name string, list []string) bool {
		list = append(list, cfg.Ignore...)
		for _, l := range list {
			for name == l {
				return false
			}
		}
		return true
	}

	for {
		if time.Now().Hour() == cfg.Hour {
			list := getWhiteList(cfg.Login, cfg.Pass, cfg.Project)
			dir, err := os.ReadDir(cfg.WorkDir)
			if err != nil {
				p.WriteLog(2, time.Now(), err.Error())
				return
			}

			for _, f := range dir {
				if check(f.Name(), list) {
					err = os.RemoveAll(cfg.WorkDir + f.Name())
					if err != nil {
						p.WriteLog(2, time.Now(), err.Error())
					} else {
						p.WriteLog(1, time.Now(), "remove "+f.Name())
					}

				}
			}
			p.WriteLog(1, time.Now(), "Wait 24h")
			time.Sleep(time.Hour * 24)
		} else {
			p.WriteLog(1, time.Now(), "Wait 1h")
			time.Sleep(time.Hour * 1)
		}
	}

}
