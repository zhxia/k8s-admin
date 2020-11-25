package httpprocessor

import (
    "errors"
    "fmt"
    log "github.com/sirupsen/logrus"
    v1 "k8s.io/api/core/v1"
    "kube-admin/src/api"
    "kube-admin/src/util"
    "strings"
)

type ServiceProcessor struct {
    HttpRequestBaseProcessor
}

func (processor ServiceProcessor) Process() (*util.JsonResponse, error) {
    var err error
    var resp *util.JsonResponse
    action := strings.ToUpper(processor.req.GetPath("action", ""))
    cfg := processor.cfg
    if action == "LIST" {
        strLabelSelector := processor.req.GetQuery("lbs")
        strFieldSelector := processor.req.GetQuery("fes")
        labelSelector := util.SelectorConvertToMap(strLabelSelector)
        fieldSelector := util.SelectorConvertToMap(strFieldSelector)
        var servs []v1.Service
        servs, err := api.NewServiceApi(cfg.K8sNamespace, cfg.K8sConfig).GetServiceList(labelSelector, fieldSelector, 500)
        if err != nil {
            log.Error(err)
            err = errors.New("get service list failed")
        } else {
            resp = util.NewJsonResponse(servs, "ok", util.ResultOk)
        }
    } else if action == "GET" {
        name := strings.Trim(processor.req.GetQuery("name"), " ")
        var serv *v1.Service
        serv, err = api.NewServiceApi(cfg.K8sNamespace, cfg.K8sConfig).GetService(name)
        if err != nil {
            log.Error(err.Error())
            err = errors.New("get service failed")
        } else {
            resp = util.NewJsonResponse(serv, "ok", util.ResultOk)
        }
    } else {
        err = errors.New(fmt.Sprintf("unsupported action:%s", action))
        log.Error(err.Error())
    }
    return resp, err
}
