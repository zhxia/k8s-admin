package httpprocessor

import (
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    "kube-admin/src/util"
)

type ValidatorProcessor struct {
    HttpRequestBaseProcessor
}

func (processor *ValidatorProcessor) Process() (resp *util.JsonResponse, err error) {
    result := util.GetJson(processor.req.Request)
    if !result.Exists() {
        err = errors.New("post data empty!")
        log.Error(err.Error())
        return
    }
    if !result.Get("data").Exists() {
        err = errors.New("data key \"data\" not exist!")
        log.Error(err.Error())
        return
    }
    _, err = util.YamlStrToApiObject(result.Get("data").String())
    return
}
