package httpprocessor

import (
    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    "kube-admin/src/util"
    "net/http"
    "strings"
)

type HttpRequestBaseProcessor struct {
    req *util.HttpRequest
    cfg *util.Config
}

func (bp *HttpRequestBaseProcessor) InitRequest(req *http.Request, cfg *util.Config) {
    bp.req = &util.HttpRequest{Request: req}
    bp.cfg = cfg
}

type HttpRequestProcessor interface {
    // need http req processor to implements this method
    Process() (*util.JsonResponse, error)
    InitRequest(req *http.Request, cfg *util.Config)
}

type BaseHandler struct {
    Config           *util.Config
    RequestProcessor HttpRequestProcessor
}

func (handler *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var err error
    var httpResponse *util.JsonResponse
    defer func(er *error) {
        if *er != nil {
            util.NewJsonResponse(nil, (*er).Error(), util.ResultError).Output(w)
        } else {
            if httpResponse == nil {
                httpResponse = util.NewJsonResponse(nil, "", util.ResultOk)
            }
            httpResponse.Output(w)
        }
    }(&err)
    namespace := r.URL.Query().Get("ns")
    if namespace != "" {
        handler.Config.K8sNamespace = namespace
    }
    // init params
    handler.RequestProcessor.InitRequest(r, handler.Config)
    // process req
    httpResponse, err = handler.RequestProcessor.Process()
}

type ApiHandlerDispatcher struct {
    Handlers map[string]BaseHandler
}

func (dispatcher *ApiHandlerDispatcher) Handle(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    apiType := strings.ToUpper(vars["type"])
    if handler, ok := dispatcher.Handlers[apiType]; ok {
        handler.ServeHTTP(w, r)
    } else {
        log.Error("illegal req!")
    }
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("<h1>welcome to use JobRunner 2.0!<h1>"))
}

func NotfoundHandler(w http.ResponseWriter, r *http.Request) {
    log.Error("req url:[", r.RequestURI, "] not found!")
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("404,sorry!"))
}
