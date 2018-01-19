package main

import (
	"fmt"
	// log "github.com/Sirupsen/logrus"
	// "github.com/liuchjlu/handler/client"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync"
	"time"
)

var m *sync.Mutex = new(sync.Mutex)

type haha struct {
	lock *sync.Mutex
	name string
}

func (h *haha) Lock() {
	h.lock.Lock()
}

func (h *haha) Unlock() {
	h.lock.Unlock()
}

func (h *haha) t1(i int) {
	h.Lock()
	defer h.Unlock()
	fmt.Println("t1():", i)
	time.Sleep(10 * time.Second)

}

func (h *haha) t2(i int) {
	h.Lock()
	defer h.Unlock()
	fmt.Println("t2():", i)
	time.Sleep(10 * time.Second)
}

// func main() {
// 	var kubeconfig string
// 	kubeconfig = "/home/liuchjlu/.kube/config"
// 	kube, err := client.NewK8sClient(kubeconfig)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	//kube.GetPodLog("liuchjlu2017121223242912-m-2l7dv ")
// 	// for {
// 	// 	result, err := kube.GetJob("liucaihong-1513666667-m")
// 	// 	if err != nil {
// 	// 		log.Errorf("err:%+v", err)
// 	// 	}
// 	// 	log.Infof("getr job result of liucaihong-1513666667-m:%+v", result)
// 	// 	time.Sleep(5)
// 	// }

// 	//for {
// 	pods, err := kube.Client.CoreV1().Pods("default").List(metav1.ListOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
// 	//log.Infof("podlist:%+v", pods)

// 	// for pod := range pods.Items {
// 	// 	log.Infof("pod info:%+v", pod)
// 	// }
// 	for {
// 		if pod, err := kube.Client.CoreV1().Pods("default").Get("5-1513753896-m-v7bb6", metav1.GetOptions{}); err != nil {
// 			log.Infof("err%+v", err)
// 		} else {
// 			log.Infof("pod: %+v", pod)
// 		}
// 		time.Sleep(3 * time.Second)

// 	}

// }

func main() {
	//fmt.Println("time:", time.UnixNano())
	var h haha
	fmt.Println("111")
	h.name = "liucaihong"
	fmt.Println("333")
	h.lock = new(sync.Mutex)
	fmt.Println("222")
	go h.t1(1)
	time.Sleep(2 * time.Second)
	h.t2(2)
}
