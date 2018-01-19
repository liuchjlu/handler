package cli

import (
	"github.com/liuchjlu/handler/client"

	log "github.com/Sirupsen/logrus"
	//"github.com/liuchjlu/handler/datatype"
	"time"
)

func IsTaskSchduled(jobsname []string, kube client.KUBE) bool {

	// podlist, err := kube.ListPod()
	// if err != nil {
	// 	log.Errorf("cli.IsTaskSchduled() kube.ListPod:%+v", err)
	// 	return false
	// }
	for i := 0; i < 20; i++ {
		var bl bool = true
		for _, name := range jobsname {
			if job, err := kube.GetJob(name); err != nil {
				log.Errorf("cli.IsTaskSchdlued() kube.GetJob %+v : %+v", name, err)
				bl = false
			} else {
				log.Infof("@@@@@@@@@@@@ Jobname:%+v   job.Status.Active:%+v, job.Status.Succeeded:%+v, job.Status.Failed:%+v", name, job.Status.Active, job.Status.Succeeded, job.Status.Failed)
				if job.Status.Active+job.Status.Succeeded+job.Status.Failed == 0 {
					bl = false
				}
			}
		}
		if bl {
			return true
		}
		time.Sleep(3 * time.Second)
	}
	return false
}

func JobStart(jobyaml []string, kube client.KUBE) error {
	log.Infof("Starting the job:%+v", jobyaml)

	//log.Infoln("$$$$$$$$$$$$$$$$ kube.client:",kube.Client)

	/*	resp1,err1:=kube.Client.CoreV1().ListNodes(context.TODO())
		if err1!=nil {
			fmt.Println("list nodes err:",err1)
		}
		//fmt.Println("list nodes gGetItems:",resp1.GetItems())
		fmt.Println("list nodes string:",resp1.String())

	*/
	for _, job := range jobyaml {

		resp, err := kube.CreateJob(job)
		if err != nil {
			log.Errorf("cli.JobStart CreateJob:", err)
			return err
		} else {
			log.Infoln("Job created:", resp)
		}
		time.Sleep(4 * time.Second)
	}

	return nil
}
