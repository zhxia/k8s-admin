package httpprocessor

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "kube-admin/src/api/crd"
    "kube-admin/src/util"
    "net/http"
    "strings"
)

type ReleaseV2Processor struct {
    HttpRequestBaseProcessor
}

func (processor *ReleaseV2Processor) Process() (resp *util.JsonResponse, err error) {
    action := strings.ToUpper(processor.req.GetPath("action", ""))
    log.Info("action:", action)
    r := processor.req.Request
    cfg := processor.cfg
    api := crd.NewCloneSetApi(cfg.K8sNamespace, cfg.K8sConfig)
    if action == "DEPLOY" {
        resp, err = processor.deploy(r, api)
    } else if action == "STATUS" {
        resp, err = processor.status(r, api)
    }
    return
}

func (processor *ReleaseV2Processor) deploy(r *http.Request, api *crd.CloneSetApi) (resp *util.JsonResponse, err error) {
    data := util.GetJson(r)
    strYaml := data.Get("yaml").String()
    cloneset, err := api.CreateOrUpdate(strYaml)
    if err != nil {
        log.Error(fmt.Sprintf("cloneset deploy failed:%s", err.Error()))
    } else {
        resp = util.NewJsonResponse(cloneset, "success", util.ResultOk)
    }
    return
}

func (processor *ReleaseV2Processor) status(r *http.Request, api *crd.CloneSetApi) (resp *util.JsonResponse, err error) {
    cs_name := strings.Trim(r.URL.Query().Get("name"), "")
    cs, err := api.GetCloneSet(cs_name)
    if err != nil {
        return
    }
    data := map[string]int32{
        "replicas":               cs.Status.Replicas,
        "ready_replicas":         cs.Status.ReadyReplicas,
        "updated_ready_replicas": cs.Status.UpdatedReadyReplicas,
        "updated_replicas":       cs.Status.UpdatedReplicas,
        "percentage":             (cs.Status.UpdatedReadyReplicas / cs.Status.Replicas) * 100,
        "partition":              *cs.Spec.UpdateStrategy.Partition,
    }
    resp = util.NewJsonResponse(data, "success", util.ResultOk)
    return
}
