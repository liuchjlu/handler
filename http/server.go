package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/liuchjlu/handler/cli"
	"github.com/liuchjlu/handler/datatype"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var LinkedQueue datatype.LinkedQueue

type Submit string
type GetLog string

// func (m *Myhttp) SubmitTask(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	r.ParseForm()
// 	if err != nil {
// 		log.Errorf("http.SubmitTask() r.ParaseForm:%+v", err)
// 	}
// 	//LinkedQueue.Add(r.Form)
// 	log.Infof("*********Size of the task queue:%+v", LinkedQueue.SizeOf())
// 	w.Write([]byte("Your task has been submit. Please wait for the schduling"))
// }
func (submit *Submit) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Errorf("http.SubmitTask() r.ParaseForm:%+v", err)
		return
	}
	var task datatype.Task
	task.Task_id = r.Form.Get("task_id")
	task.User_name = r.Form.Get("user_name")
	task.Script_path = r.Form.Get("script_path")
	task.Gpu_num, err = strconv.Atoi(r.Form.Get("gpu_num"))
	task.Image = r.Form.Get("image_name")
	if err != nil {
		log.Errorf("http.ServerHttp strconv.Atoi:%+v", err)
	}
	log.Infof("task info:%+v", task)
	LinkedQueue.Add(task)

	w.Write([]byte("Your task has been submit. Please wait for the schduling"))
}
func (getlog *GetLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Errorf("http.SubmitTask() r.ParaseForm:%+v", err)
		return
	}
	podName := r.Form.Get("pod_name")
	podLog, err1 := cli.GetPodLog(podName)
	if err1 != nil {
		log.Errorf("http.ServerHTTP() r.Form.Get", err1)
	}

	w.Write([]byte(podLog))
}

func main() {

	logFilename := "./handler.log"
	logFile, _ := os.OpenFile(logFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer logFile.Close()

	writers := []io.Writer{
		logFile,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)

	log.SetOutput(fileAndStdoutWriter)
	log.SetLevel(log.DebugLevel)

	log.Infoln("main.main():Start  Main")

	cli.InitConfig("../cfg/sysconfig.yaml")
	LinkedQueue.LockQ = new(sync.Mutex)
	go cli.Run(&LinkedQueue)
	fmt.Println("Starting server.")

	mux := http.NewServeMux()

	submit := new(Submit)
	getlog := new(GetLog)
	//mux.Handle("/submit", m)
	mux.Handle("/submit", submit)
	mux.Handle("/getlog", getlog)

	log.Println("Listening...")
	http.ListenAndServe("localhost:3000", mux)

}
