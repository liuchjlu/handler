//package main

package client

import (
	//"flag"
	"os"
	//"path/filepath"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	//"time"
	//"sync"

	"github.com/ghodss/yaml"
	"io/ioutil"
	//meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	//"k8s.io/apimachinery/pkg/api/errors"
	log "github.com/Sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	//clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	watch "k8s.io/apimachinery/pkg/watch"
	//restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//restclient "k8s.io/client-go/rest"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

//var wg sync.WaitGroup
var DeleteOptions metav1.DeletionPropagation = "Background"

type KUBE struct {
	Client *kubernetes.Clientset
}

var Quantity resource.Quantity

const (
	ResourceName corev1.ResourceName = "alpha.kubernetes.io/nvidia-gpu"
)

func NewK8sClient(kubeconfig string) (KUBE, error) {
	var kube KUBE
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Errorf("client.NewK8sClient(). clientcmd.BuildConfigFromFlags:%+v \n", err)
		return kube, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("client.NewK8sClient(). kubernetes.NewForConfig:%+v \n", err)
	}
	kube.Client = clientset
	return kube, nil
}
func (kube KUBE) GetJob(jobname string) (*v1.Job, error) {
	batchClient := kube.Client.BatchV1()
	jobsClient := batchClient.Jobs("default")
	//metav1.GetOptions{ResourceVersion: 1}
	result, err := jobsClient.Get(jobname, metav1.GetOptions{ResourceVersion: "1"})
	if err != nil {
		log.Errorf("client.GetJob(). When get the job of %+v : %+v", jobname, err)
		return result, err
	}
	return result, err
}

func (kube KUBE) DeleteJob(jobname string) error {
	batchClient := kube.Client.BatchV1()
	jobsClient := batchClient.Jobs("default")
	err := jobsClient.Delete(jobname, &metav1.DeleteOptions{PropagationPolicy: &DeleteOptions})
	if err != nil {
		log.Errorf("client.DeleteJob() err:%+v", err)
	}
	return err
}

//func (kube KUBE) WatchNode() (watch.Interface, error) {
func (kube KUBE) WatchNode() (watch.Interface, error) {
	resp, err := kube.Client.CoreV1().Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		log.Errorf("client.WatchNode() err:%+v", err)
		return resp, err
	}
	log.Infof("client.WatchNode() resp:%+v", resp.ResultChan())
	return resp, err
}

func (kube KUBE) ListPod() (result *corev1.PodList, err error) {
	resp, err := kube.Client.CoreV1().Pods("default").List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("client.ListPod() err:%+v", err)
	}
	return resp, err
}

func (kube KUBE) GetPod(name string) (result *corev1.Pod, err error) {
	resp, err := kube.Client.CoreV1().Pods("default").Get(name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("client.GetPod() err:%+v", err)
	}
	return resp, err
}

func (kube KUBE) GetResourceQuotas(name string) (result *corev1.ResourceQuota, err error) {
	resp, err := kube.Client.CoreV1().ResourceQuotas("default").Get(name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("client.GetResourceQuotas() err:%+v", err)
	}
	log.Infof("client.GetResourceQuotas() resp:%+v", resp)
	return resp, err
}

func (kube KUBE) ListResourceQuotas() (result *corev1.ResourceQuotaList, err error) {
	resp, err := kube.Client.CoreV1().ResourceQuotas("default").List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("client.ListResourceQuotas() err:%+v", err)
	}
	log.Infof("client.ListResourceQuotas() resp:%+v", resp)
	return resp, err
}

func (kube KUBE) GetNode(name string) (result *corev1.Node, err error) {
	resp, err := kube.Client.CoreV1().Nodes().Get(name, metav1.GetOptions{ResourceVersion: "1"})
	if err != nil {
		log.Errorf("client.ListResourceQuotas() err:%+v", err)
	}
	//log.Infof("client.ListResourceQuotas() resp:%+v", resp)
	return resp, err
}

func (kube KUBE) ListNode() (result *corev1.NodeList, err error) {
	resp, err := kube.Client.CoreV1().Nodes().List(metav1.ListOptions{Watch: true})
	if err != nil {
		log.Errorf("client.ListResourceQuotas() err:%+v", err)
	}
	//log.Infof("client.ListResourceQuotas() resp:%+v", resp)
	return resp, err
}

func (kube KUBE) CreateJob(jobyamlpath string) (*v1.Job, error) {

	//var listoption metav1.ListOptions

	batchClient := kube.Client.BatchV1()
	jobsClient := batchClient.Jobs("default")

	/*  piJob, err := jobsClient.Get("testname-videokeywords-v1-0",getoption)
	    check(err)
	    fmt.Printf("piJob Name: %v\n", piJob.Name)

	    jobsList, err := jobsClient.List(listoption)
	    check(err)

	    // Loop over all jobs and print their name
	    for i, job := range jobsList.Items {
	        fmt.Printf("Job %d: %s\n", i, job.Name)
	    }*/

	var job *v1.Job
	var result *v1.Job
	data, err := ioutil.ReadFile(jobyamlpath)
	if err != nil {
		log.Errorf("client.CreateJob() ioutil.ReadFile: %+v", err)
		return job, err
	}
	if err := yaml.Unmarshal(data, &job); err != nil {
		log.Errorf("client.CreateJob(). unmarshal jobyaml: %+v ", err)
		return job, err
	}
	result, err = jobsClient.Create(job)
	if err != nil {
		log.Errorf("client.CreateJob(). jobsClient.Create: %+v", err)
		log.Infof("client.CreateJob(). result:%+v", result)
		return result, err
	}
	log.Debugf("Create job: %+v", result)
	return result, err
}

func (kube KUBE) GetPodFullName(name string) (string, error) {
	result, err := kube.ListPod()
	if err != nil {
		log.Errorf("client.GetPodFullName() kube.ListPod:%+v", err)
	}

	for _, pod := range result.Items {
		if strings.HasPrefix(pod.ObjectMeta.Name, name+"-m") {
			return pod.ObjectMeta.Name, nil
		}
	}
	return "", err
}

func (kube KUBE) GetPodLog(name string) (string, error) {
	//var option corev1.PodLogOptions{Follow: follow}
	//var req restclient.Request
	podName, err := kube.GetPodFullName(name)
	if err != nil {
		log.Errorf("client.GetPodLog() kube.GetPodFullName:%+v", err)
		return "", err
	} else if name == "" {
		log.Errorf("the pod %+v is not found.", name)
		return "", err
	}
	log.Infof("client.GetPodLog() podName:%+v", podName)
	req := kube.Client.CoreV1().Pods("default").GetLogs(podName, &corev1.PodLogOptions{Follow: false})
	//go func(kube KUBE) *restclient.Request {

	//	return kube.Client.CoreV1().Pods("default").GetLogs(podName, &corev1.PodLogOptions{Follow: true})
	//}(kube)

	//log.Infof("Logs of the pod %s :%+v", podName, req)
	readcloser, err := req.Stream()
	if err != nil {
		log.Errorf("client.GetPodLog() err:%+v", err)
		return "", err
	}
	log.Infof("io.ReadCloser:%+v", readcloser)
	buf := new(bytes.Buffer)
	// reader := bufio.NewReader(stdout)

	// cmd := exec.Command("", "")
	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	log.Errorf("cmd.ExecCommand():%+v\n", err)
	// }
	// cmd.Start()
	// for {

	// }

	buf.ReadFrom(readcloser)
	s := buf.String()
	//log.Infof("********s:%+v", s)
	//fmt.Printf(" **pod log:%+v", buf)
	fmt.Printf("pod log:%+v", s)
	// for {
	// 	line, err := buf.ReadString('\n')
	// 	if err != nil {
	// 		log.Errorf("client.GetPodLog() err:$+v", err)
	// 		break
	// 	}
	// 	fmt.Println(line)
	// }
	return s, nil

}
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
func ExecCommand(commandName string, logpath string, params []string) error {
	cmd := exec.Command(commandName, params...)

	//显示运行的命令
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatalf("cmd.ExecCommand():%+v\n", err)
		return err
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)
	// var fo *os.File
	// if checkFileIsExist(logpath) {
	// 	fo, err = os.OpenFile(logpath, os.O_APPEND|os.O_RDWR, 0755)
	// } else {
	// 	fo, err = os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	// }
	// //fo, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	// //os.Create(logpath)
	// if err != nil {
	// 	panic(err)
	// }
	// defer fo.Close()
	// writer := bufio.NewWriter(fo) //创建输出缓冲流

	// buf := make([]byte, 1024)
	// for {
	// 	n, err := reader.Read(buf)
	// 	if err != nil && err != io.EOF {
	// 		panic(err)
	// 	}
	// 	if n == 0 {
	// 		break
	// 	}

	// 	if n2, err := writer.Write(buf[:n]); err != nil {
	// 		panic(err)
	// 	} else if n2 != n {
	// 		panic("error in writing")
	// 	}
	// }

	// if err = writer.Flush(); err != nil {
	// 	panic(err)
	// }

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}

		// n, err := writer.WriteString(line)
		// if err != nil || n == -1 {
		// 	panic("error in writing")
		// }

		fmt.Printf(line)
	}

	cmd.Wait()
	return nil
}

func main() {
	params := []string{"--server", "http://192.168.12.49:8080", "logs", "-f", os.Args[1]}
	ExecCommand("/home/liuchjlu/kubernetes1.6/kube-master/kubectl", "/home/liuchjlu/kubernetes1.6/kube-master/log", params)
	// var kubeconfig string
	// /* if home := homeDir(); home != "" {
	//        kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	//    } else {
	//        kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	//    }
	//    flag.Parse()*/
	// kubeconfig = "/home/liuchjlu/.kube/config"
	// kube, err := NewK8sClient(kubeconfig)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// //kube.GetPodLog("liuchjlu2017121223242912-m-2l7dv ")
	// for {
	// 	result, err := kube.GetJob("liucaihong-1513666667-m")
	// 	if err != nil {
	// 		log.Errorf("err:%+v", err)
	// 	}
	// 	log.Infof("getr job result of liucaihong-1513666667-m:%+v", result)
	// 	time.Sleep(5)
	// }

	// for {
	//     pods, err := kube.Client.CoreV1().Pods("").List(metav1.ListOptions{})
	//     if err != nil {
	//         panic(err.Error())
	//     }
	//     fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	//     // Examples for error handling:
	//     // - Use helper functions like e.g. errors.IsNotFound()
	//     // - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	//     _, err = kube.Client.CoreV1().Pods("default").Get("205rc1-ghn2p", metav1.GetOptions{})
	//     if errors.IsNotFound(err) {
	//         fmt.Printf("Pod not found\n")
	//     } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
	//         fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	//     } else if err != nil {
	//         panic(err.Error())
	//     } else {
	//         fmt.Printf("Found pod\n")
	//     }

	//     time.Sleep(10 * time.Second)
	// }
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

/*package client
import (
    "io/ioutil"
    "context"

    "github.com/ericchiang/k8s"
    batchv1 "github.com/ericchiang/k8s/apis/batch/v1"
    "github.com/ghodss/yaml"
    log "github.com/Sirupsen/logrus"
)

// loadClient parses a kubeconfig from a file and returns a Kubernetes
// client. It does not support extensions or client auth providers.
type KUBE  struct {
    Client *k8s.Client
}

func LoadClient(kubeconfigPath string) (KUBE, error) {
    var kube KUBE
    data, err := ioutil.ReadFile(kubeconfigPath)
    if err != nil {
        log.Errorf("read kubeconfig: %v", err)
        return kube,err
    }

    // Unmarshal YAML into a Kubernetes config object.
    var config k8s.Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        log.Errorf("unmarshal kubeconfig: %v", err)
        return kube,err
    }
    log.Debugln("k8sconfig:",config)
    kube.Client,err=k8s.NewClient(&config)
    log.Infoln("**11111111**")
    if err != nil {
        log.Errorf("Client.LoadClient. k8s.NewClient:",err)
    }
    log.Infoln("****")
    return kube,err
}

//obj *batchv1.Job
func (kube KUBE)CreateJob(ctx context.Context, jobyamlpath string) (*batchv1.Job, error){
    // var c *k8s.BatchV1
    // //c.client=client
    // c.client=client
    var job *batchv1.Job
    data,err := ioutil.ReadFile(jobyamlpath)
    if err != nil {
        log.Errorf("1111read jobyamlpath: %v", err)
        return job,err

    }
    if err := yaml.Unmarshal(data, &job); err != nil {
        log.Errorf("2222unmarshal jobyaml: %v", err)
        return job,err
    }
    jobres,err := kube.Client.CreateJob(ctx, job)
    if err!=nil{
        log.Errorf("3333CreateJob: ",err)
    }
    return jobres,nil

}*/
