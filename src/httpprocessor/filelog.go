package httpprocessor

import (
    "errors"
    "fmt"
    log "github.com/sirupsen/logrus"
    "kube-admin/src/api"
    "kube-admin/src/api/crd"
    "kube-admin/src/util"
    "strconv"
    "strings"
)

type FileLogProcessor struct {
    HttpRequestBaseProcessor
}

// /log/action/
func (processor *FileLogProcessor) Process() (*util.JsonResponse, error) {
    var resp *util.JsonResponse
    var err error
    r := processor.req.Request
    cfg := processor.cfg
    action := strings.ToUpper(r.URL.Query().Get("action"))
    logApi := api.NewLogApi(cfg.K8sNamespace, cfg.K8sConfig)
    pod := strings.Trim(r.URL.Query().Get("pod"), "")
    container := strings.Trim(r.URL.Query().Get("container"), "")
    if action == "LIST" {
        clonesetName := strings.Trim(r.URL.Query().Get("csn"), "")
        cloneset, err := crd.NewCloneSetApi(cfg.K8sNamespace, cfg.K8sConfig).GetCloneSet(clonesetName)
        if err == nil {
            for _, vm := range cloneset.Spec.Template.Spec.Containers[0].VolumeMounts {
                if vm.Name == "log-dir" {
                    logDir := vm.MountPath
                    cmd := []string{"/bin/bash", "-c", fmt.Sprintf("ls %s", logDir)}
                    out, err := logApi.ExecRemoteCmd(cfg.K8sNamespace, pod, container, cmd)
                    if err == nil {
                        data := make(map[string]string)
                        data["log_dir"] = logDir
                        data["log_files"] = string(out)
                        resp = util.NewJsonResponse(data, "success", util.ResultOk)
                    }
                }
            }
        }
    } else if action == "VIEW" {
        limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
        if err != nil {
            limit = 500
        }
        logFile := strings.Trim(r.URL.Query().Get("logfile"), "")
        cmd := []string{"/bin/bash", "-c", fmt.Sprintf("tail -n%d %s", limit, logFile)}
        log.Info("pod:", pod, "container:", container, "namespace:", cfg.K8sNamespace)
        out, err := logApi.ExecRemoteCmd(cfg.K8sNamespace, pod, container, cmd)
        log.Info("over!")
        if err != nil {
            return nil, err
        }
        resp = util.NewJsonResponse(string(out), "success", util.ResultOk)
    } else {
        err = errors.New(fmt.Sprintf("unsupported action:%s", action))
    }
    return resp, err
}
