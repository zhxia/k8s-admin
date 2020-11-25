package httpprocessor

import (
    "fmt"
    kruiseapps "github.com/openkruise/kruise-api/apps/v1alpha1"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
    "kube-admin/src/api/crd"
    "kube-admin/src/util"
    "net/http"
    "strings"
)

type CloneSetProcessor struct {
    HttpRequestBaseProcessor
}

func (processor *CloneSetProcessor) Process() (resp *util.JsonResponse, err error) {
    action := strings.ToUpper(processor.req.GetPath("action", ""))
    cfg := processor.cfg
    r := processor.req.Request
    api := crd.NewCloneSetApi(cfg.K8sNamespace, cfg.K8sConfig)
    if action == "LIST" {
        resp, err = processor.list(r, api)
    } else if action == "GET" {
        resp, err = processor.get(r, api)
    } else if action == "UPDATE" {
        resp, err = processor.update(r, api)
    } else {
        err = errors.New(fmt.Sprintf("unsupported action:%s", action))
    }
    return
}

func (processor *CloneSetProcessor) list(r *http.Request, api *crd.CloneSetApi) (resp *util.JsonResponse, err error) {
    strLabelSelector := r.URL.Query().Get("lbs")
    strFieldSelector := r.URL.Query().Get("fes")
    labelSelector := util.SelectorConvertToMap(strLabelSelector)
    fieldSelector := util.SelectorConvertToMap(strFieldSelector)
    var clonesetList []kruiseapps.CloneSet
    clonesetList, err = api.GetCloneSetList(labelSelector, fieldSelector, 500)
    if err != nil {
        log.Error(err.Error())
        err = errors.New("get cloneset list failed!")
    } else {
        resp = util.NewJsonResponse(clonesetList, "", util.ResultOk)
    }
    return
}

func (processor *CloneSetProcessor) get(r *http.Request, api *crd.CloneSetApi) (resp *util.JsonResponse, err error) {
    name := strings.Trim(r.URL.Query().Get("name"), "")
    cloneset, err := api.GetCloneSet(name)
    if err != nil {
        log.Error(err.Error())
        err = errors.New(fmt.Sprintf("get cloneset:%s failed!", name))
    } else {
        resp = util.NewJsonResponse(cloneset, "", util.ResultOk)
    }
    return
}

func (processor *CloneSetProcessor) update(r *http.Request, api *crd.CloneSetApi) (resp *util.JsonResponse, err error) {
    data := util.GetJson(r)
    name := data.Get("name").String()
    replicas := int32(data.Get("replicas").Int())
    requestCpu := data.Get("request_cpu").String()
    requestMem := data.Get("request_mem").String()
    limitCpu := data.Get("limit_cpu").String()
    limitMem := data.Get("limit_mem").String()
    var (
        ptrReplicas            *int32
        ptrPartition           *int32
        ptrResourceRequirement *v1.ResourceRequirements
    )
    if replicas > 0 {
        ptrReplicas = &replicas
    }
    if requestCpu != "" && requestMem != "" {
        ptrResourceRequirement.Requests = map[v1.ResourceName]resource.Quantity{
            v1.ResourceName("memory"): resource.MustParse(requestMem),
            v1.ResourceName("cpu"):    resource.MustParse(requestCpu),
        }
    }
    if limitCpu != "" && limitMem != "" {
        ptrResourceRequirement.Limits = map[v1.ResourceName]resource.Quantity{
            v1.ResourceName("memory"): resource.MustParse(limitMem),
            v1.ResourceName("cpu"):    resource.MustParse(limitCpu),
        }
    }
    cs, err := api.Update(name, "", ptrReplicas, ptrPartition, ptrResourceRequirement, false)
    if err == nil {
        resp = util.NewJsonResponse(cs, "", util.ResultOk)
    }
    return
}
