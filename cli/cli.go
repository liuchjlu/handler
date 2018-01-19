package cli

import (
	log "github.com/Sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"

	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/liuchjlu/handler/client"
	"github.com/liuchjlu/handler/datatype"
)

var (
	Task      datatype.Task      = datatype.Task{}
	TaskPool  datatype.TaskPool  = datatype.TaskPool{}
	Worker    datatype.Worker    = datatype.Worker{}
	SysConfig datatype.SysConfig = datatype.SysConfig{}
	JobInfo   datatype.JobInfo   = datatype.JobInfo{}
	//save the gpu_num info of the cluster
	Gmap map[string]int
)

// These are the valid statuses of pods.
const (
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending corev1.PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning corev1.PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	PodSucceeded corev1.PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed corev1.PodPhase = "Failed"
	// PodUnknown means that for some reason the state of the pod could not be obtained, typically due
	// to an error in communicating with the host of the pod.
	PodUnknown corev1.PodPhase = "Unknown"
)

const (
	// PodScheduled represents status of the scheduling process for this pod.
	PodScheduled corev1.PodConditionType = "PodScheduled"
	// PodReady means the pod is able to service requests and should be added to the
	// load balancing pools of all matching services.
	PodReady corev1.PodConditionType = "Ready"
	// PodInitialized means that all init containers in the pod have started successfully.
	PodInitialized corev1.PodConditionType = "Initialized"
	// PodReasonUnschedulable reason in PodScheduled PodCondition means that the scheduler
	// can't schedule the pod right now, for example due to insufficient resources in the cluster.
	PodReasonUnschedulable = "Unschedulable"
)

func InitConfig(sysconfig_path string) {
	var err error
	SysConfig, err = datatype.UnmarshalSysConfig(sysconfig_path, SysConfig)
	if err != nil {
		log.Fatalln("Failed to read sysconfig:", err)
	}
	log.Debugf("InitConfig success:", SysConfig)
}

// A data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value int
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.

var Total int

func sortMapByValue(m map[string]int) PairList {
	p := make(PairList, len(m))
	Total = 0
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		Total += v
		i++
	}
	log.Infof("Total gpu card:%+v", Total)
	sort.Sort(sort.Reverse(p))
	return p
}

func UpdateGmap(Etcd *client.Etcd, podlist *corev1.PodList) (map[string]int, error) {
	Gmap = make(map[string]int)
	if resp, err := Etcd.GetAbsoluteDir("/registry/minions"); err != nil {
		log.Errorf("cli.UpdateGmap() Etcd.GetAbsoluteDir:", err)
		return Gmap, err
	} else {
		for _, node := range resp.Node.Nodes {
			log.Debugf("The info of node %s :%+v", node.Key, node.Value)
			//capacity":{"alpha.kubernetes.io/nvidia-gpu":"
			gpu_num, err := strconv.Atoi(strings.Split(strings.Split(node.Value, `capacity":{"alpha.kubernetes.io/nvidia-gpu":"`)[1], `"`)[0])
			if err != nil {
				log.Errorf("Update gpu resource failed:", err)
			}
			host := strings.Split(node.Key, "/")[3]
			Gmap[host] = gpu_num
			//Gmap[node.Key]=2
		}
	}
	for _, pod := range podlist.Items {
		if pod.Status.Phase == PodRunning || pod.Status.Phase == PodUnknown || pod.Status.Phase == PodPending {
			node_name := pod.Spec.NodeName
			client.Quantity = pod.Spec.Containers[0].Resources.Requests[client.ResourceName]
			gpu_requests, _ := client.Quantity.AsInt64()
			Gmap[node_name] -= int(gpu_requests)
		}

	}
	log.Infof("******** Gpu resource Updated:%+v", Gmap)
	return Gmap, nil
}

//func Binding(gnum int, gmap map[string]int)(map[string]int,error){
func Binding(gpu_need int, gmap map[string]int) map[string]int {
	gmap_bind := make(map[string]int)
	glist := sortMapByValue(gmap)
	if gpu_need > Total {
		log.Errorf("cli.Binding() GPU is not enough")
		return gmap_bind
	}
	for i, list := range glist {
		if gpu_need > list.Value {
			gmap_bind[list.Key] = list.Value
			gpu_need -= list.Value
		}
		if gpu_need == list.Value {
			gmap_bind[list.Key] = list.Value
			gpu_need -= list.Value
			return gmap_bind
		}
		if gpu_need < list.Value {
			list = match(i, glist, gpu_need)
			gmap_bind[list.Key] = gpu_need
			return gmap_bind
		}

	}
	log.Infof("Completed Gpu Binding:%+v", gmap_bind)
	return gmap_bind
}

//match the node that generate the least gpu gragment
func match(i int, p PairList, gpu_need int) Pair {
	for j := len(p) - 1; j >= i; j-- {
		if p[j].Value >= gpu_need {
			return p[j]
		}
	}
	return p[i]
}

func SchduleOneTask(task datatype.Task) (error, bool) {
	log.Infof("Start schduling the task:%+v", task)
	kube, err := client.NewK8sClient(SysConfig.K8s_cfg)
	if err != nil {
		log.Errorf("cli.SchduleOneTask()) client.NewK8sClient:", err)
		return err, false
	}
	etcd, err := client.NewEtcdClient(SysConfig.Etcd_addr)
	if err != nil {
		log.Errorf("cli.SchduleOneTask() NewEtcdClient:%+v", err)
		return err, false
	}
	podlist, err := kube.ListPod()
	if err != nil {
		log.Errorf("cli.SchduleOneTask() kube.ListPod:%+v", err)
		return err, false
	}
	gmap, err := UpdateGmap(etcd, podlist)
	if err != nil {
		log.Errorf("cli.SchduleOneTask() UpdateGmap:%+v", err)
		return err, false
	}
	gmap_bind := Binding(task.Gpu_num, gmap)
	jobspath, jobsname, err := LoadTask(SysConfig, task, gmap_bind)
	log.Infof("cli.SchduleOneTask() jobsname:%+v", jobsname)
	if err != nil {
		log.Errorf("cli.SchduleOneTask() LoadTask:%+v", err)
		return err, false
	}
	etcddir := "/usertask/" + task.Task_id + "/" + strconv.Itoa(len(gmap_bind))
	log.Infof("&&&&&&&&&&& etcddir:%+v", etcddir)
	err = etcd.CreateAbsoluteDir(etcddir)
	if err != nil {
		log.Errorf("cli.SchduleOneTask() etcd.CreateAbsoluteDir:%+v", err)
		return err, false
	}
	if err1 := JobStart(jobspath, kube); err != nil {
		log.Errorf("cli.SchduleOneTask() JobStart:%+v", err)
		return err1, false
	}
	if IsTaskSchduled(jobsname, kube) {
		return nil, true
	}
	//else {
	// 	if etcd.DeleteDir(etcddir); err != nil {
	// 		log.Errorf("cli.SchduleOneTask() etcd.DeleteDir:%+v", err)
	// 		return nil, false
	// 	}
	// }
	return nil, false
}
func GetPodLog(name string) (string, error) {
	kube, err := client.NewK8sClient(SysConfig.K8s_cfg)
	if err != nil {
		log.Errorf("cli.GetLog() client.NewK8sClient:%+v", err)
		return "", err
	}
	podlog, err1 := kube.GetPodLog(name)
	if err != nil {
		log.Errorf("cli.GetLog() kube.GetPodLog:%+v", err1)
		return "", err1
	}
	return podlog, nil

}

func Run(queue *datatype.LinkedQueue) {
	for {
		log.Infof("cli.Run() queue.size:%+v", queue.SizeOf())
		if queue.SizeOf() != 0 {
			log.Infof("@@@@@@ cli.Run() queue.Sizeof():", queue.SizeOf())
			task := queue.Peek()
			err, b := SchduleOneTask(task)
			if err != nil {
				log.Errorf("cli.Run() SchduleOneTask. Task info:%+v Error info:%+v", task, err)
			} else if b {
				log.Infof("@@@@@@Task run successful: %+v \n", task)
				log.Infoln("Waiting for another task to schdule...")
				queue.Remove()
			} else {
				log.Warnf("cli.Run() SchduleOneTask. Your task has been submit, but maybe schdule failer for the reason resource limit. Task info: %+v", task)
			}
			// for i := 0; i < 3; i++ {
			// 	err, b := SchduleOneTask(task)
			// 	if err != nil {
			// 		log.Errorf("cli.Run() SchduleOneTask. Task info:%+v Error info:%+v", task, err)
			// 		if i != 2 {
			// 			log.Infof("Now will try to schdule the task again.")
			// 		} else {
			// 			log.Infof("Fail to schdule the task after tried 3 times. The task will not be reschdule anain.")
			// 		}
			// 	} else {
			// 		queue.Remove()
			// 		log.Infof("Task run successful: %+v", task)
			// 		break
			// 	}
			// 	time.Sleep(1 * time.Second)
			// }
		}

		time.Sleep(5 * time.Second)
	}
}
