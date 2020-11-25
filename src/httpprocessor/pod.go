package httpprocessor

import (
    "errors"
    "fmt"
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    v1 "k8s.io/api/core/v1"
    "kube-admin/src/api"
    "kube-admin/src/util"
    "strconv"
    "strings"
)

type PodProcessor struct {
    HttpRequestBaseProcessor
}

func (processor PodProcessor) Process() (*util.JsonResponse, error) {
    var err error
    var resp *util.JsonResponse
    r := processor.req.Request
    cfg := processor.cfg
    vars := mux.Vars(r)
    action := strings.ToUpper(vars["action"])
    if action == "LIST" {
        strLabelSelector := r.URL.Query().Get("lbs")
        strFieldSelector := r.URL.Query().Get("fes")
        labelSelector := util.SelectorConvertToMap(strLabelSelector)
        fieldSelector := util.SelectorConvertToMap(strFieldSelector)
        var pods []v1.Pod
        pods, err = api.NewPodApi(cfg.K8sNamespace, cfg.K8sConfig).GetPodList(labelSelector, fieldSelector, 5000)
        if err != nil {
            log.Error(err.Error())
            err = errors.New("get pod list failed")
        } else {
            resp = util.NewJsonResponse(pods, "ok", util.ResultOk)
        }
    } else if action == "GET" {
        name := strings.Trim(r.URL.Query().Get("name"), " ")
        var pod *v1.Pod
        pod, err = api.NewPodApi(cfg.K8sNamespace, cfg.K8sConfig).GetPod(name)
        if err != nil {
            log.Error(err)
            err = errors.New("get pod failed")
        } else {
            resp = util.NewJsonResponse(pod, "ok", util.ResultOk)
        }
    } else if action == "LOGS" {
        name := strings.Trim(r.URL.Query().Get("name"), " ")
        container := strings.Trim(r.URL.Query().Get("container"), " ")
        limit, _ := strconv.Atoi(strings.Trim(r.URL.Query().Get("limit"), " "))
        strLogs, err := api.NewPodApi(cfg.K8sNamespace, cfg.K8sConfig).GetPodLogs(name, container, int64(limit))
        if err != nil {
            log.Error(err)
            err = errors.New("get pod logs failedS")
        } else {
            resp = util.NewJsonResponse(strLogs, "ok", util.ResultOk)
        }
    } else {
        err = errors.New(fmt.Sprintf("unsupported action:%s", action))
        log.Error(err.Error())
    }
    return resp, err
}
