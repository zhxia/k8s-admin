package httpprocessor

import (
    "errors"
    "fmt"
    uuid "github.com/satori/go.uuid"
    log "github.com/sirupsen/logrus"
    appsv1 "k8s.io/api/apps/v1"
    "kube-admin/src/api"
    "kube-admin/src/util"
    "strings"
)

type ReleaseProcessor struct {
    HttpRequestBaseProcessor
}

func (processor *ReleaseProcessor) Process() (*util.JsonResponse, error) {
    var err error
    var resp *util.JsonResponse
    jsonObject := processor.req.GetJson()
    cfg := processor.cfg
    action := strings.ToUpper(jsonObject.Get("action").String())
    reqId := uuid.NewV4().String()
    newDeployName := jsonObject.Get("new_deploy").String()
    lastDeployName := jsonObject.Get("last_deploy").String()
    maxReplicas := jsonObject.Get("max_replicas").Int()
    parallel := jsonObject.Get("parallel").Int()
    if parallel == 0 {
        parallel = 1
    }
    releaseApi := api.NewReleaseApi(reqId, action, cfg.K8sNamespace, cfg.K8sConfig, int32(maxReplicas))
    if action == "DEPLOY" {
        var dep *appsv1.Deployment
        dep, err = releaseApi.Deploy(jsonObject.Get("yaml").String())
        if err != nil {
            log.Error(err)
            err = errors.New("release deploy failed")
        } else {
            resp = util.NewJsonResponse(map[string]interface{}{"dep": dep, "reqId": reqId}, "ok", util.ResultOk)
        }
    } else if action == "SCALE" {
        replicas := jsonObject.Get("replicas").Int()
        releaseApi.ScaleTimeout = jsonObject.Get("scale_timeout").Int()
        go func() {
            err = releaseApi.DeployScale(newDeployName, lastDeployName, int32(replicas), int32(parallel))
            if err != nil {
                log.Error(err)
                err = errors.New("deploy scale update failed")
            }
        }()
        result := map[string]interface{}{
            "reqId": reqId,
        }
        resp = util.NewJsonResponse(result, "ok", util.ResultOk)
    } else if action == "RESULT" {

    } else {
        err = errors.New(fmt.Sprintf("unsupported action:%s", action))
        log.Error(err.Error())
    }
    return resp, err
}
