package httpprocessor

import (
    "github.com/pkg/errors"
    "kube-admin/src/api"
    "kube-admin/src/util"
    "net/http"
    "strings"
)

type SecretProcessor struct {
    HttpRequestBaseProcessor
}

func (processor *SecretProcessor) Process() (resp *util.JsonResponse, err error) {
    cfg := processor.cfg
    r := processor.req.Request
    secretApi := api.NewSecretApi(cfg.K8sNamespace, cfg.K8sConfig)
    action := strings.ToUpper(processor.req.GetPath("action", ""))
    if action == "CREATE" {
        resp, err = processor.create(r, secretApi)
    } else if action == "GET" {
        resp, err = processor.get(r, secretApi)
    } else {
        err = errors.New("unsupported action:" + action)
    }
    return
}

func (processor *SecretProcessor) create(r *http.Request, api *api.SecretApi) (resp *util.JsonResponse, err error) {
    jsonObject := util.GetJson(r)
    secret, err := api.CreateSecret(jsonObject.Get("yaml").String())
    if err == nil {
        resp = util.NewJsonResponse(secret, "", util.ResultOk)
    }
    return
}

func (processor *SecretProcessor) get(r *http.Request, api *api.SecretApi) (resp *util.JsonResponse, err error) {

    return
}
