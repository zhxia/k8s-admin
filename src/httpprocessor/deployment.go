package httpprocessor

import (
    "errors"
    "fmt"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    appsv1 "k8s.io/api/apps/v1"
    autoscalingv1 "k8s.io/api/autoscaling/v1"
    "kube-admin/src/api"
    "kube-admin/src/util"
    "strconv"
    "strings"
)

type DeploymentProcessor struct {
    HttpRequestBaseProcessor
}

func (processor *DeploymentProcessor) Process() (*util.JsonResponse, error) {
    var resp *util.JsonResponse
    var err error
    var r = processor.req.Request
    var cfg = processor.cfg
    vars := mux.Vars(r)
    action := strings.ToUpper(vars["action"])
    if action == "LIST" {
        strLabelSelector := r.URL.Query().Get("lbs")
        strFieldSelector := r.URL.Query().Get("fes")
        labelSelector := util.SelectorConvertToMap(strLabelSelector)
        fieldSelector := util.SelectorConvertToMap(strFieldSelector)
        var deps []appsv1.Deployment
        deps, err = api.NewDeploymentApi(cfg.K8sNamespace, cfg.K8sConfig).GetDeploymentList(labelSelector, fieldSelector,
            500)
        if err != nil {
            log.Error(err.Error())
            err = errors.New("get deployment list failed")
        } else {
            log.Info("get des list len:", len(deps), "namespace:", cfg.K8sNamespace, ",cfg:", cfg.K8sConfig)
            resp = util.NewJsonResponse(deps, "ok", 0)
        }
    } else if action == "CREATE" {
        jsonObject := util.GetJson(r)
        deployYaml := jsonObject.Get("deploy").String()
        var dep *appsv1.Deployment
        dep, err = api.NewDeploymentApi(cfg.K8sNamespace, cfg.K8sConfig).CreateDeployment(deployYaml)
        if err != nil {
            log.Error(err.Error())
            err = errors.New("create deployment failed")
        } else {
            resp = util.NewJsonResponse(dep, "ok", 0)
        }
    } else if action == "UPDATE" {
        jsonObject := util.GetJson(r)
        deployYaml := jsonObject.Get("deploy").String()
        var dep *appsv1.Deployment
        dep, err = api.NewDeploymentApi(cfg.K8sNamespace, cfg.K8sConfig).UpdateDeployment(deployYaml)
        if err != nil {
            log.Error(err.Error())
            err = errors.New("update deployment failed")
        } else {
            resp = util.NewJsonResponse(dep, "ok", 0)
        }
    } else if action == "GET-SCALE" {
        name := r.URL.Query().Get("name")
        var scale *autoscalingv1.Scale
        scale, err = api.NewDeploymentApi(cfg.K8sNamespace, cfg.K8sConfig).GetScale(name)
        if err != nil {
            log.Error(err.Error())
            err = errors.New("get scale failed")
        } else {
            resp = util.NewJsonResponse(scale, "ok", 0)
        }
    } else if action == "UPDATE-SCALE" {
        jsonObject := util.GetJson(r)
        name := jsonObject.Get("name").String()
        var replicas int
        replicas, err = strconv.Atoi(jsonObject.Get("replicas").String())
        if err != nil {
            log.Error(err.Error())
            err = errors.New("get replicas params failed")
        } else {
            var scale *autoscalingv1.Scale
            scale, err = api.NewDeploymentApi(cfg.K8sNamespace, cfg.K8sConfig).UpdateScale(name, int32(replicas))
            if err != nil {
                log.Error(err.Error())
                err = errors.New("update scale failed")
            } else {
                resp = util.NewJsonResponse(scale, "ok", 0)
            }
        }
    } else if action == "GET" {
        name := r.URL.Query().Get("name")
        var dep *appsv1.Deployment
        dep, err = api.NewDeploymentApi(cfg.K8sNamespace, cfg.K8sConfig).GetDeployment(name)
        if err != nil {
            log.Error(err.Error())
            err = errors.New("get deployment failed")
        } else {
            resp = util.NewJsonResponse(dep, "ok", 0)
        }
    } else {
        err = errors.New(fmt.Sprintf("unsupported action:%s", action))
        log.Error(err.Error())
    }
    return resp, err
}
