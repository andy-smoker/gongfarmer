package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/andy-smoker/clerk"
)

// get issues list with status "In Testing" or "Ready for Testing"
func getWhiteList(login, pass, project string) ([]string, error) {
	req, err := http.NewRequest("GET", `https://`+project+`.atlassian.net/rest/api/latest/search?jql=project="HRL"%20AND%20status%20in%20("In%20Testing","Ready%20for%20Testing")`, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(login, pass)
	cl := http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// struct for parsing response
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
	return list, nil

}

// CFG - struct for parsing config file
type CFG struct {
	Logpath        string   `toml:"log"`
	WorkDir        string   `toml:"work_directory"`
	ServiceDirs    []string `toml:"service_dirs"`
	Ignore         []string `toml:"ignore_list"`
	Login          string   `toml:"login"`
	JiraToken      string   `toml:"jira_token"`
	HourOfCleaning int      `toml:"hour_of_cleaning"`
	Period         int      `toml:"period"`
	Project        string   `toml:"project"`
}

// getConfig - parse file config.toml
func getConfig() *CFG {
	cfg := &CFG{
		Logpath: "./rm.log",
		WorkDir: "./",
	}
	toml.DecodeFile("config.toml", cfg)
	return cfg
}

func main() {
	cfg := getConfig()
	// create new logs printer
	p := clerk.NewPrinter("INFO", "gongfarmer", cfg.Logpath)
	p.WriteLog(1, time.Now(), "start")
	// anonimus func for check file/directory name in ignore list
	check := func(name string, list []string) bool {
		for _, l := range list {
			for name == l {
				return false
			}
		}
		return true
	}

	for {
		if time.Now().Hour() == cfg.HourOfCleaning {
			cfg = getConfig()

			list, err := getWhiteList(cfg.Login, cfg.JiraToken, cfg.Project)
			if err != nil {
				p.WriteLog(2, time.Now(), err.Error())
			}
			list = append(list, cfg.Ignore...)

			dir, err := os.ReadDir(cfg.WorkDir)
			if err != nil {
				p.WriteLog(2, time.Now(), err.Error())
			}

			for _, f := range dir {
				if check(f.Name(), list) {
					err = os.RemoveAll(cfg.WorkDir + f.Name())
					if err != nil {
						p.WriteLog(2, time.Now(), err.Error())
					} else {
						p.WriteLog(1, time.Now(), f.Name()+" removed")
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

func CleanDir(p *clerk.Printer) {

}
