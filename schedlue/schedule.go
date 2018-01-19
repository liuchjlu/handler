package schedule

import (
	"../client"
	"../data"
	"container/list"
	"encoding/json"
)

func (etcd *client.Etcd) GpuAvailableList() *list.List {

	return list.New()
}

func (etcd *client.Etcd) TaskSplit(gpu_available_list *list.List, task data.Task) (work_list *list.List) {
	return list.New()
}

func YamlBuilder(work_list *list.List) interface{} {

}

func (etcd *client.Etcd) TaskRegister(work_list *list.List) {

}

func (etcd *client.ETCD) TaskDeleter(task_id string) {

}

func (k8s *client.K8s) CreateJob(yaml_str string) {

}
