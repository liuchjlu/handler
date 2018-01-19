package cli

import (
	"strconv"

	"github.com/liuchjlu/handler/datatype"
	"github.com/liuchjlu/handler/template"

	log "github.com/Sirupsen/logrus"
)

func LoadTask(cfg datatype.SysConfig, task datatype.Task, gmap_bind map[string]int) ([]string, []string, error) {
	var joblist = make([]datatype.JobInfo, len(gmap_bind))
	var jobnames = make([]string, len(joblist))
	var job datatype.JobInfo
	i := 0
	log.Infoln("Total pod:%v (master+slave)", len(gmap_bind))
	for k, v := range gmap_bind {
		if i == 0 {
			job.Job_id = task.Task_id + "-m"
			job.Start_path = "/start_m.sh"
		} else {
			job.Job_id = task.Task_id + "-s" + strconv.Itoa(i)
			job.Start_path = "/start_s.sh"
		}
		jobnames[i] = job.Job_id
		job.Task_id = task.Task_id
		job.User_name = task.User_name
		job.Pod_num = len(gmap_bind)
		job.Image = task.Image
		job.Script_path = task.Script_path
		job.Gpu_num = v
		job.Etcd_addr = cfg.Etcd_addr
		job.K8s_addr = cfg.K8s_addr
		job.Nfs_server = cfg.Nfs_server
		job.Nfs_path = "/nfsroot/" + task.User_name
		job.Host_node = k

		log.Infof("Jobinfo of %+v : %+v \n", job.Job_id, job)
		//joblist = append(joblist,job)
		joblist[i] = job
		i++
	}
	jobspath := make([]string, len(joblist))
	log.Infoln("len of joblist:", len(joblist))
	for j, jobinfo := range joblist {
		one_yaml, err := template.BuildYaml(jobinfo)
		if err != nil {
			log.Errorf("BuildYaml for the job:%v", jobinfo)
			return jobspath, jobnames, err
		} else {
			log.Infoln("Job yaml path of %s : %s", JobInfo.Job_id, one_yaml)
			jobspath[j] = one_yaml
		}
	}
	return jobspath, jobnames, nil

}
