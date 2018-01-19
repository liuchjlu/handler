package template

import (
	"os"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/liuchjlu/handler/datatype"
)

func BuildYaml(job datatype.JobInfo) (string, error) {
	yamlpath := "/tmp/yaml/" + job.User_name + "/"
	os.MkdirAll(yamlpath, 0755)
	filename := yamlpath + "job-" + job.Job_id + ".yaml"
	dstFile, err := os.Create(filename)
	defer dstFile.Close()
	check(err)

	t, err := template.New("job.template").ParseFiles("/home/liuchjlu/workspace/go/src/github.com/liuchjlu/handler/template/job.template")
	check(err)
	data := struct {
		Job_id      string
		Task_id     string
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
		Mount_path  string
	}{
		Task_id:     job.Task_id,
		Job_id:      job.Job_id,
		Pod_num:     job.Pod_num,
		Image:       job.Image,
		Script_path: job.Script_path,
		Gpu_num:     job.Gpu_num,
		Etcd_addr:   job.Etcd_addr,
		K8s_addr:    job.K8s_addr,
		Start_path:  job.Start_path,
		Nfs_server:  job.Nfs_server,
		Nfs_path:    job.Nfs_path,
		Host_node:   job.Host_node,
		Mount_path:  "/home/" + job.User_name,
	}
	err = t.Execute(dstFile, data)
	check(err)
	dstFile.WriteString("\n")
	dstFile.WriteString("---\n")

	return filename, nil
}
func BuilTest(task datatype.Task) error {
	yamlpath := "/tmp/testyaml/"
	os.MkdirAll(yamlpath, 0755)
	filename := yamlpath + task.User_name + ".yaml"
	dstFile, err := os.Create(filename)
	defer dstFile.Close()
	check(err)

	t, err := template.New("test.template").ParseFiles("/home/liuchjlu/workspace/go/src/github.com/liuchjlu/handler/template/test.template")
	check(err)
	data := struct {
		User_name string
	}{
		User_name: task.User_name,
	}
	err = t.Execute(dstFile, data)
	check(err)
	dstFile.WriteString("\n")
	dstFile.WriteString("---\n")

	return nil
}

func check(err error) {
	if err != nil {
		log.Errorf("err:%+v", err)
	}
}
