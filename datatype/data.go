package datatype

import (
	"container/list"
	//log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Task struct {
	Task_id   string `json:"task_id"`
	User_name string `json:"user_name"`
	//Status      string `json:"status"`
	Gpu_num     int    `json:"gpu_num"`
	Script_path string `json:"script_path"`
	Image       string `json:"image_name"`
}

type TaskPool struct {
	Queue *list.List
}

type Worker struct {
	Task_id    string
	Machine_ip string
	Gpu_num    int
	Master     bool
}

type SysConfig struct {
	K8s_addr   string
	K8s_cfg    string
	Etcd_addr  string
	Nfs_server string
}

type JobInfo struct {
	Job_id      string
	Task_id     string
	User_name   string
	Pod_num     int
	Image       string
	Script_path string
	Gpu_num     int
	Etcd_addr   string
	K8s_addr    string
	Start_path  string
	Nfs_server  string
	Nfs_path    string
	Host_node   string
}

/*
type Element struct {
	//Value interface{} //在元素中存储的值
	Task Task
}

func (t *TaskPool) PushFront(task Task) {
	t.Queue.PushFront(task)
}

func (t *TaskPool) PushBack(task Task) {
	t.Queue.PushBack(task)
}

func (t *TaskPool) Front() (*Task, bool) {
	var task Task
	if !t.IsEmpty() {
		e := t.Queue.Front()
		task = e.Value.Task
		log.Infof("Get task from front of task_queue successful: %+v", task)
		return task, true
	}
	return task, false

}

func (t *TaskPool) Back() *Task {
	var task Task
	if !t.IsEmpty() {
		e := t.Queue.Back()
		task = e.Value.Task
		log.Infof("Get task from back of task_queue successful: %+v", task)
		return task, true
	}
	return task, false

}

func (t *TaskPool) Remove(task *Task) {

}

func (l *TaskPool) IsEmpty() bool {
	return l.Queue.root.next == &l.root
}
*/
func UnmarshalSysConfig(path string, data SysConfig) (SysConfig, error) {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return data, err
	}
	err = yaml.Unmarshal(in, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func UnmarshalTask(path string, data Task) (Task, error) {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return data, err
	}
	err = yaml.Unmarshal(in, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
