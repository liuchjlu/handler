package main

import (
	"fmt"
	"io"
	"os"

	//"encoding/json"
	//"strings"
	//"reflect"
	//"context"

	"github.com/liuchjlu/handler/cli"
	"github.com/liuchjlu/handler/client"
	//"github.com/liuchjlu/handler/datatype"

	log "github.com/Sirupsen/logrus"
)

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

	log.Infoln("main.main():Start handler Main")
	//cli.Run()

	etcdpath := "http://192.168.12.48:2379"
	etcd, err := client.NewEtcdClient(etcdpath)
	if err != nil {
		fmt.Println("NewEtcdClient err:", err)
	}
	mp, err := cli.UpdateGmap(etcd)
	for node := range mp {
		fmt.Println("node:", node, " gpu_num:", mp[node])
	}
	gmap_bind := cli.Binding(4, mp)
	for g := range gmap_bind {
		fmt.Println("bind node:", g, "bind num:", gmap_bind[g])
	}
	cli.InitConfig("./cfg/sysconfig.yaml")
	log.Infoln("cli.SysConfig:", cli.SysConfig)
	cli.Task.Task_id = "liuchjlu2017121223242912"
	cli.Task.User_name = "liuchjlu"
	cli.Task.Status = "not ready"
	cli.Task.Gpu_num = 4
	cli.Task.Script_path = "/home/liuchjlu/handler/handler.sh"
	jobpath, err := cli.LoadTask(cli.SysConfig, cli.Task, gmap_bind)
	if err != nil {
		fmt.Println("cli.LoadTask:", err)
	}
	fmt.Println("jobpath:", jobpath)

	cli.JobStart(jobpath)

	/*    resp,err:=etcd.GetAbsoluteDir("/registry/minions/")
	      if err!=nil{
	          fmt.Println("GetAbsoluteDir err:",err)
	      }


	      fmt.Println("GetAbsoluteDir /registry/minions  resp.Node :",resp.Node)
	      fmt.Println("GetAbsoluteDir /registry/minions  resp.Node.Nodes :",resp.Node.Nodes)
	      for _,node := range resp.Node.Nodes{
	          fmt.Println("node.key:",node.Key)
	          //fmt.Println("node.value:",node.Value)
	          err=json.Unmarshal([]byte(node.Value),&nodeinfo)
	          if err!=nil {
	              fmt.Println("json.Unmarshal err:",err)
	          }
	          //fmt.Println(nodeinfo.Spec.Allocatable.Cpu )
	          //fmt.Println(nodeinfo.Spec.Allocatable)
	          fmt.Println("node.value.type:",reflect.TypeOf(node.Value))
	          fmt.Println("#####liuchjlu: ",strings.Split(strings.Split(node.Value,`allocatable":{"alpha.kubernetes.io/nvidia-gpu":"`)[1],`"`)[0])
	      }*/

	/*    client,err:=client.LoadClient("./cfg/k8sconfig.yaml")
	if err!=nil {
		fmt.Println("err1: ",err)
	}
	fmt.Println("client: %v",client)
	resp,err:=client.CoreV1().ListEndpoints(context.TODO(),"default")
	if err!=nil {
		fmt.Println("list endpoints err:",err)
	}
	fmt.Println("list endpoints:",resp)
	resp1,err1:=client.CoreV1().ListNodes(context.TODO())
	if err1!=nil {
		fmt.Println("list nodes err:",err1)
	}
	//fmt.Println("list nodes gGetItems:",resp1.GetItems())
	fmt.Println("list nodes string:",resp1.String())
	//fmt.Println("list nodes GetMetadata :",resp1.GetMetadata())
	resp2,err2:=client.CoreV1().GetReplicationController(context.TODO(),"205rc","default")
	if err2!=nil {
		fmt.Println("GetReplicationController err:",err2)
	}
	fmt.Println("GetReplicationController:",resp2.String())
	//t.Println("list nodes string:",resp1.String())*/
}
