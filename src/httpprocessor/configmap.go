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

type ConfigProcessor struct {
    HttpRequestBaseProcessor
}

func (processor *ConfigProcessor) Process() (*util.JsonResponse, error) {
    var err error
    var resp *util.JsonResponse
    r := processor.req.Request
    action := strings.ToUpper(processor.req.GetPath("action", ""))
    cfg := processor.cfg
    if action == "LIST" {
        strLabelSelector := r.URL.Query().Get("lbs")
        strFieldSelector := r.URL.Query().Get("fes")
        labelSelector := util.SelectorConvertToMap(strLabelSelector)
        fieldSelector := util.SelectorConvertToMap(strFieldSelector)
        var cms []v1.ConfigMap
        cms, err = api.NewConfigMapApi(cfg.K8sNamespace, cfg.K8sConfig).GetConfigMapList(labelSelector, fieldSelector, 500)
        if err != nil {
            log.Error(err.Error())
            err = errors.New("get cfg map list failed")
        } else {
            resp = util.NewJsonResponse(cms, "ok", util.ResultOk)
        }
    } else if action == "GET" {
        name := strings.Trim(r.URL.Query().Get("name"), " ")
        var cm *v1.ConfigMap
        cm, err = api.NewConfigMapApi(cfg.K8sNamespace, cfg.K8sConfig).GetConfigMap(name)
        if err != nil {
            log.Error(err.Error())
            err = errors.New("get cfg map failed")
        } else {
            resp = util.NewJsonResponse(cm, "ok", util.ResultOk)
        }
    } else {
        err = errors.New(fmt.Sprintf("unsupported action:%s", action))
        log.Error(err.Error())
    }
    return resp, err
}
