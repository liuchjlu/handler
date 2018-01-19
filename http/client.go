package main

import (
	"fmt"
	//"github.com/liuchjlu/handler/cli"
	//"github.com/liuchjlu/handler/datatype"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	//log "github.com/Sirupsen/logrus"
)

func submit_test() {
	postParam := url.Values{
		"user_name":   {"liuchjlu"},
		"task_id":     {"liuchjlu-233333"},
		"gpu_num":     {"4"},
		"script_path": {"/home/liuchjlu/tensorflow_mnist.py"},
		"image_name":  {"192.168.12.41:5000/test:latest"},
	}

	resp, err := http.PostForm("http://localhost:3000/submit", postParam)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}
func getLog_test() {
	postParam := url.Values{
		"pod_name": {"liuchjlu-233333"},
	}

	resp, err := http.PostForm("http://localhost:3000/getlog", postParam)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}
func main() {
	if os.Args[1] == "submit" {
		submit_test()
	} else if os.Args[1] == "getlog" {
		getLog_test()
	}
	//submit_test()
	//getLog_test()
}
