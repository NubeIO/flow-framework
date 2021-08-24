package api

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

var CRON  *gocron.Scheduler

// The JobDatabase interface for encapsulating database access.
type JobDatabase interface {
	GetJob(uuid string) (*model.Job, error)
	GetJobs() ([]*model.Job, error)
	CreateJob(body *model.Job) error
	UpdateJob(uuid string, body *model.Job) (*model.Job, error)
	DeleteJob(uuid string) (bool, error)

}
type JobAPI struct {
	DB JobDatabase
}

func reposeHandlerError(err error, ctx *gin.Context) {

	res := BadEntity(err.Error())
	ctx.JSON(res.GetStatusCode(), res.GetResponse())

}


func reposeHandler(body interface{}, err error, ctx *gin.Context) {
	if err != nil {
		res := BadEntity(err.Error())
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	} else {
		res := Data(body)
		ctx.JSON(res.GetStatusCode(), res.GetResponse())
	}
}

func (j *JobAPI) GetJobs(ctx *gin.Context) {
	q, err := j.DB.GetJobs()
	reposeHandler(q, err, ctx)

}

func (j *JobAPI) CreateJob(ctx *gin.Context) {
	body, _ := getBODYJobs(ctx)
	_, err := govalidator.ValidateStruct(body)
	if err != nil {
		reposeHandler(nil, err, ctx)
	}
	err = j.DB.CreateJob(body)
	reposeHandler(body, err, ctx)
}


func (j *JobAPI) UpdateJob(ctx *gin.Context) {
	body, _ := getBODYJobs(ctx)
	uuid := resolveID(ctx)
	q, err := j.DB.UpdateJob(uuid, body)
	reposeHandler(q, err, ctx)
}


func (j *JobAPI) GetJob(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.GetJob(uuid)
	reposeHandler(q, err, ctx)
}


func (j *JobAPI) DeleteJob(ctx *gin.Context) {
	uuid := resolveID(ctx)
	q, err := j.DB.DeleteJob(uuid)
	reposeHandler(q, err, ctx)
}


/*
add job
- don't start job until it has one or more subscribers

edit job
- if is set disable then notify all subscriber's and for enable do the same

delete job
- notify all subscriber's, and they will unsubscribe

job subscriber
- On add: make sure the job uuid is valid
- On delete: update the subscriber's list and unsubscribe

remote subscriber
- required: rubix-uuid
- optional: point uuid (required network_name and device_name and point_name)

*/

func (j *JobAPI) initCron() {
	CRON = 	gocron.NewScheduler(time.UTC)
	CRON.StartAsync()
	//j.syncJobs()
}


//syncJobs start all the jobs on start of the app
//func (j *JobAPI) syncJobs()  {
//	q, err := j.DB.GetJobs()
//	for _, job := range q {
//		if job.Enable {
//			for _, jobSub := range job.JobSubscriber {
//				if jobSub.Enable {
//					err = j.jobAdd(job.UUID, &job)
//					if err != nil {
//						log.Println("error on read job")
//					}
//				}
//			}
//		}
//	}
//}

//jobAdd add a job
func (j *JobAPI) jobAdd(uuid string, body *model.Job) error {
	if body.Frequency == "" {
		return errors.New("invalid time frequency, example 5m")
	}
	_, err := CRON.Every(body.Frequency).Tag(uuid).Do(taskWithParams, uuid, body)
	if err != nil {
		return err
	}
	return nil
}


func (j *JobAPI) jobRemover(uuid string) error {
	err := CRON.RemoveByTag(uuid)
	if err != nil {
		return errors.New("error on remove job")
	}
	return nil

}

func (j *JobAPI) NewJobEngine() {
	j.initCron()
	log.Println("INIT CRON")
}





func taskWithParams(uuid string, body *model.Job) {
	fmt.Println(uuid)
	//payload := new(payloadBody)
	//payload.UUID = uuid
	//payload.Delete = false
	//payload.MessageString = "what up"
	//payload.MessageTS = time.Now().Format(time.RFC850)
	//topic := fmt.Sprintf("%s:%s", "job",uuid)
	//
	//BUS.RegisterTopics(topic)
	//err := BUS.Emit(BusBackground, topic, payload)
	//
	//fmt.Println("topics", BUS.Topics())
	//if err != nil {
	//	fmt.Println("error", err)
	//}

}
