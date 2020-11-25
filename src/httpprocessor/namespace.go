package httpprocessor

import (
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    "kube-admin/src/api"
    "kube-admin/src/util"
    "strings"
)

type NamespaceProcessor struct {
    HttpRequestBaseProcessor
}

func (processor *NamespaceProcessor) Process() (*util.JsonResponse, error) {
    var resp *util.JsonResponse
    var err error
    action := strings.ToUpper(processor.req.GetPath("action", ""))
    cfg := processor.cfg
    if action == "LIST" {
        strLabelSelector := processor.req.GetQuery("lbs", "")
        strFieldSelector := processor.req.GetQuery("fes", "")
        labelSelector := util.SelectorConvertToMap(strLabelSelector)
        fieldSelector := util.SelectorConvertToMap(strFieldSelector)
        nsList, err := api.NewNamespaceApi(cfg.K8sNamespace, cfg.K8sConfig).GetNamespaceList(labelSelector, fieldSelector, 100)
        if err != nil {
            log.Error(err.Error())
            err = errors.New("get namespace list failed!")
        } else {
            resp = util.NewJsonResponse(nsList, "success", util.ResultOk)
        }
    } else {
        err = errors.New("unsupported action")
    }
    return resp, err
}
