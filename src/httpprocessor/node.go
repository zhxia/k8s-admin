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

type NodeProcessor struct {
    HttpRequestBaseProcessor
}

func (processor *NodeProcessor) Process() (*util.JsonResponse, error) {
    var err error
    var resp *util.JsonResponse
    r := processor.req.Request
    cfg := processor.cfg
    action := strings.ToUpper(processor.req.GetPath("action"))
    if action == "LIST" {
        strLabelSelector := r.URL.Query().Get("lbs")
        strFieldSelector := r.URL.Query().Get("fes")
        labelSelector := util.SelectorConvertToMap(strLabelSelector)
        fieldSelector := util.SelectorConvertToMap(strFieldSelector)
        nodes, err := api.NewNodeApi(cfg.K8sNamespace, cfg.K8sConfig).GetNodeList(labelSelector, fieldSelector, 500)
        if err != nil {
            log.Error(err.Error())
            err = errors.New("get node list failed")
        } else {
            resp = util.NewJsonResponse(nodes, "ok", util.ResultOk)
        }
    } else if action == "GET" {
        name := strings.Trim(r.URL.Query().Get("name"), " ")
        var node *v1.Node
        node, err = api.NewNodeApi(cfg.K8sNamespace, cfg.K8sConfig).GetNode(name)
        if err != nil {
            log.Error(err.Error())
            err = errors.New("get node detail failed")
        } else {
            resp = util.NewJsonResponse(node, "ok", util.ResultOk)
        }
    } else {
        err = errors.New(fmt.Sprintf("unsupported action:%s", action))
        log.Error(err.Error())
    }
    return resp, err
}
