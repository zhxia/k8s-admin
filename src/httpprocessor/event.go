package httpprocessor

import (
    "kube-admin/src/util"
)

type EventProcessor struct {
    HttpRequestBaseProcessor
}

func (processor *EventProcessor) Process() (resp *util.JsonResponse, err error) {

    return
}
