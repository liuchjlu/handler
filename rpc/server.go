package main

import (
	"fmt"
	"net/http"
	"net/rpc"
	//"time"
	"github.com/liuchjlu/handler/cli"
	"github.com/liuchjlu/handler/datatype"
	"io"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type MyRPC int

var (
	LinkedQueue datatype.LinkedQueue
	m           *sync.Mutex = new(sync.Mutex)
)

func (r *MyRPC) HelloRPC(S string, reply *string) error {
	fmt.Println(S)
	*reply = "This is Server. Hello Client RPC."
	return nil
}

/*func (r *MyRPC) CreateJob(task datatype.Task,response *string)error{
	if err:=cli.SchduleOneTask(task);err!=nil{
		log.Errorf("rpc.CreateJob() SchduleOneTask:%+v",err)
		*response = "Schdule Task Error."
		return err
	}
	*response = "Your Task has been schduled. Please wait for the stdout."
	return nil
}*/

func (r *MyRPC) CreateJob(task datatype.Task, response *string) error {

	log.Infof("Waiting for schduling the task: %+v", task)
	m.Lock()
	err, b := cli.SchduleOneTask(task)
	m.Unlock()
	if err != nil {
		log.Errorf("rpc.CreateJob() SchduleOneTask:%+v", err)
		*response = "Schdule Task Error."
		return err
	}
	log.Infof("Status is %+v of SchduleOneTask of the task:%+v", b, task)
	log.Infof("@@@ The task has been schduled: %+v \n", task)
	*response = "Your Task has been schduled. Please wait for the stdout."
	return nil
}

func (r *MyRPC) SubmitTask(task datatype.Task, response *string) error {

	LinkedQueue.Add(task)
	log.Infof("*********Size of the task queue:%+v", LinkedQueue.SizeOf())
	*response = "Your Task has been submit. Please wait for the schduling."
	return nil
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
	r := new(MyRPC)

	rpc.Register(r)
	rpc.HandleHTTP()
	err := http.ListenAndServe("localhost:1111", nil)
	if err != nil {
		fmt.Println("in main", err.Error())
	}

}
