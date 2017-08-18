package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var testGroups = [2]string{"killpod", "nodeexec"}

type KillPodData struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Interval  int    `json:"interval"`
	Slack     bool   `json:"slack"`
	Fatal     bool   `json:"fatal"`
}

type RunnerExecData struct {
	Command    string   `json:"command"`
	Parameters []string `json:"parameters"`
}

type ExecData struct {
	Target   string           `json:"target"`
	Commands []RunnerExecData `json:"commands"`
}

func configPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	l := filepath.Join(cwd, "cookies.d")
	log.Println("checking directory: ", l)
	if _, err := os.Stat(l); err == nil {
		return l
	}
	l = filepath.Join("/", "etc", "cookies.d")
	log.Println("checking directory: ", l)
	if _, err := os.Stat(l); err == nil {
		return l
	}
	panic("Can not file configuration directory")
}

/*
// parse JSON from request body
func readJSONData(r *http.Request, data interface{}) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, data); err != nil {
		panic(err)
	}

	log.Println("request data: ", data)
}
*/

func listAvailableTests(group string) []string {
	var testList []string
	if group == "" {
		for _, g := range testGroups {
			testList = append(testList, listAvailableTests(g)...)
		}
	} else {
		groupPath := filepath.Join(configPath(), group)
		if _, err := os.Stat(groupPath); os.IsNotExist(err) {
			panic(err)
		}
		fileList, err := ioutil.ReadDir(groupPath)
		if err != nil {
			panic(err)
		}
		for _, fl := range fileList {
			testList = append(testList, fl.Name()[:len(fl.Name())-5])
		}
	}
	return testList
}

func loadJSON(group, name string, data interface{}) {
	testFile := filepath.Join(configPath(), group, name+".json")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		panic(err)
	}
	log.Println("loading ", testFile)
	b, err := ioutil.ReadFile(testFile)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, data); err != nil {
		panic(err)
	}
}

func loadKillPod(name string) KillPodData {
	var data KillPodData
	loadJSON("killpod", name, &data)
	return data
}

func loadNodeExec(name string) ExecData {
	var data ExecData
	loadJSON("nodeexec", name, &data)
	return data
}