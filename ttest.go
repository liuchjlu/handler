package main

import (
	//log "github.com/Sirupsen/logrus"
	// "github.com/liuchjlu/handler/datatype"
	// "github.com/liuchjlu/handler/template"
	"github.com/liuchjlu/handler/client"
	"os"
	//"os"
	// "strconv"
	//"time"
	"fmt"
)

func main() {
	/*	var task datatype.Task = datatype.Task{}
		for i := 0; i < 100; i++ {
			task.User_name = strconv.Itoa(i)
			log.Infof("task:%+v", task)
			err := template.BuilTest(task)
			if err != nil {
				log.Errorf("err:%+v", err)
			}
			time.Sleep(300 * time.Microsecond)
		}*/
	var kubeconfig string

	kubeconfig = "/home/liuchjlu/.kube/config"
	kube, err := client.NewK8sClient(kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	kube.GetPodLog(os.Args[1])
	if result, err := kube.ListPod(); err != nil {
		panic(err.Error())
	} else {
		for _, pod := range result.Items {
			fmt.Printf("****pod:", pod)
		}
	}

	//params := []string{"--server", "http://192.168.12.49:8080", "logs", "-f", os.Args[1]}
	//client.ExecCommand("/home/liuchjlu/kubernetes1.6/kube-master/kubectl", "/home/liuchjlu/kubernetes1.6/kube-master/log", params)
	// err = kube.DeleteJob(os.Args[1])
	// if err != nil {
	// 	log.Errorln("err:", err)
	// }

	/*	//result, err := kube.GetNode("cx49")
		kube.GetPodLog("liucaihong-1514860004-m-x7pss")
		//result, err := kube.ListResourceQuotas()
		if err != nil {
			log.Errorf("err:%+v", err)
		}*/

	//log.Infof("result.Items:%+v", result.Items)
	/*
		for i, pod := range result.Items {
			log.Infof("%+v:pod:%+v", i, pod)
			log.Infof("%+v:podspec:%+v", i, pod.Spec)
			log.Infof("%+v:podspec.containers0:%+v", i, pod.Spec.Containers[0])
			log.Infof("%+v:podstatus:%+v \n", i, pod.Status.Phase)
			resource_map := pod.Spec.Containers[0].Resources.Requests
			log.Infof("resource_map:%+v", resource_map)
			client.Quantity = pod.Spec.Containers[0].Resources.Requests[client.ResourceName]
			gpunum, flag := client.Quantity.AsInt64()
			//gpunum, _ := pod.Spec.Containers[0].Resources.Requests[client.ResourceName].AsInt64()
			log.Infof("gpunum:%+v flag:%+v", gpunum, flag)
			log.Infof("hostname:%+v", pod.Spec.NodeName)
			log.Infof("is schduled:%+v", pod.Status)
			log.Infof("%+v:podspec.GpuResourceNum:%+v \n", i, pod.Spec.Containers[0].Resources.Requests)

		}*/

}
