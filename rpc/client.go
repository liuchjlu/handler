package main

import (
	"fmt"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/liuchjlu/handler/cli"
	"github.com/liuchjlu/handler/datatype"

	log "github.com/Sirupsen/logrus"
)

var S string

var Team struct {
	User_name   string `json:"user_name"`
	Gpu_num     int    `json:"gpu_num"`
	Script_path string `json:"script_path"`
}

func Run() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "team.yaml")
		os.Exit(1)
	}
	var task datatype.Task
	var err error
	task, err = datatype.UnmarshalTask(os.Args[2], task)
	if err != nil {
		log.Fatalln("rpc.Run() UnmarshalTask:", err)
	}

}

func main() {

	cli.InitConfig("../cfg/sysconfig.yaml")
	/*	if len(os.Args) != 3 {
			fmt.Println("usage:",os.Args[0], "ip:port")
			os.Exit(1)
		}

		addr:= os.Args[1]
		client,err:=rpc.DialHTTP("tcp",addr)
		if err!= nil {
			log.Fatal("dialhttp:",err)
		}
		var reply *string
		S = "This is Client. Hello Server RPC."
		err = client.Call("MyRPC.HelloRPC",S,&reply)
		if err != nil{
			log.Fatal("call hellorpc:",err)
		}
		fmt.Println(*reply)


		if len(os.Args) != 3 {
			fmt.Println("Usage:",os.Args[0], "ip:port  team.yaml")
			os.Exit(1)
		}*/
	var task datatype.Task
	var err error
	task, err = datatype.UnmarshalTask(os.Args[2], task)
	if err != nil {
		log.Fatal("rpc.main() UnmarshalTask:", err)
	}
	task.Task_id = task.User_name + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	log.Infoln("Your Task info:", task)

	addr := os.Args[1]
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Fatal("dialhttp:", err)
	}
	var reply *string
	//err = client.Call("MyRPC.SubmitTask", task, &reply)
	err = client.Call("MyRPC.SubmitTask", task, &reply)
	if err != nil {
		log.Fatal("call SubmitTask:", err)
	}
	fmt.Println(*reply)
}
