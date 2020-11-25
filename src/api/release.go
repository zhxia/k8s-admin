package api

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	"time"
)

//ReleaseApi 应用发版服务类
type ReleaseApi struct {
	RequestId    string
	Action       string
	CreatedAt    string
	MaxReplicas  int32
	DepApi       *DeploymentApi
	ScaleTimeout int64
}

func NewReleaseApi(reqId, action, namespace, kubeConfig string, maxReplicas int32) *ReleaseApi {
	return &ReleaseApi{
		RequestId:   reqId,
		Action:      action,
		MaxReplicas: maxReplicas,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		DepApi:      NewDeploymentApi(namespace, kubeConfig),
	}
}

//Deploy 镜像部署
func (ra *ReleaseApi) Deploy(deploymentYaml string) (*appsv1.Deployment, error) {
	return ra.DepApi.CreateDeployment(deploymentYaml)
}

//Scale 镜像扩缩容
func (ra *ReleaseApi) DeployScale(newDeployName, lastDeployName string, replicas, parallel int32) error {
	channel := make(chan int32)
	overFlag := make(chan bool)
	if replicas > ra.MaxReplicas {
		return errors.New(fmt.Sprintf("replicas[%d] should less or equal than max replicas[%d]", replicas, ra.MaxReplicas))
	}
	scale, err := ra.DepApi.GetScale(newDeployName)
	if err != nil {
		return err
	}
	if lastDeployName == "" {
		_, err := ra.DepApi.UpdateScale(newDeployName, replicas)
		if err != nil {
			return err
		}
		return nil
	}
	// 扩容新版本
	go func() {
		var i, n int32
		defer func() {
			close(channel)
		}()
		i = scale.Spec.Replicas + 1
		n = 0
		for ; i <= replicas; i++ {
			n++
			if n%parallel == 0 {
				err := ra.ScaleUpdatedEnsure(newDeployName, i, true)
				if err != nil {
					return
				}
				channel <- i
				log.Info("push to channel:", i)
			}
		}
		if n%parallel != 0 {
			err := ra.ScaleUpdatedEnsure(newDeployName, replicas, true)
			if err != nil {
				return
			}
			channel <- replicas
			log.Info("push to channel:", replicas)
		}

	}()

	// 缩容旧版本
	go func() {
		defer func() {
			overFlag <- true
		}()
		for {
			n, ok := <-channel
			if !ok {
				log.Info("channel closed!")
				break
			}
			log.Info("pull from channel:", n)
			oldReplicas := ra.MaxReplicas - n
			ra.ScaleUpdatedEnsure(lastDeployName, oldReplicas, false)
		}
	}()
	<-overFlag
	log.Info("deploy scale finished[task_id:", ra.RequestId, "]!")
	return nil
}

func (ra *ReleaseApi) ScaleUpdatedEnsure(deployName string, replicas int32, isReady bool) (err error) {
	_, err = ra.DepApi.UpdateScale(deployName, replicas)
	if err != nil {
		log.Error("deployment:[", deployName, "] update scale[", replicas, "] failed:", err)
		return
	}
	if ra.ScaleTimeout == 0 {
		ra.ScaleTimeout = 300
	}
	start := time.Now().Unix()
	for {
		time.Sleep(time.Second * 2)
		if time.Now().Unix()-start >= ra.ScaleTimeout { //超过timeout扩容没有准备就绪，结束流程
			err = errors.New(fmt.Sprintf("deployment:[%s] scale up timeout!", deployName))
			log.Info(err.Error())
			break
		}
		dep, err := ra.DepApi.GetDeployment(deployName)
		if err != nil {
			log.Error("get deployment:[", deployName, "] detail failed!")
			continue
		}
		if isReady {
			if dep.Status.ReadyReplicas >= replicas {
				log.Info("deployment:[", deployName, "] scale to[", replicas, "] succeed and status ready!")
				break
			}
		} else {
			if dep.Status.Replicas >= replicas {
				log.Info("deployment:[", deployName, "] scale to[", replicas, "] succeed and status unknown!")
				break
			}
		}
		log.Info("trying to get deployment:[", deployName, "] status,expect replicas:", replicas, ",but current:", dep.Status.ReadyReplicas)
	}
	return
}
