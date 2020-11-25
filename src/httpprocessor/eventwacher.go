package httpprocessor

import (
    "kube-admin/src/httpprocessor/watcher"
    "kube-admin/src/util"
    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
    log "github.com/sirupsen/logrus"
    "net/http"
    "strings"
)

var eventWatcherUpgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

var eventWatchers = map[string]watcher.BaseWatcher{
    "EVENT":   watcher.EventWatcher{},
    "DEPLOYMENT":watcher.DevelopmentWatcher{},
}

type EventWatcherHandler struct {
    Config *util.Config
}

func (handler EventWatcherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    conn, err := eventWatcherUpgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Error("websocket upgrade failed:", err.Error())
        return
    }
    defer func() {
        conn.Close()
        log.Info("connection closed!")
    }()
    log.Info("client [", conn.RemoteAddr().String(), "] connected!")
    strLabelSelector := r.URL.Query().Get("lbs")
    strFieldSelector := r.URL.Query().Get("fes")
    labelSelector := util.SelectorConvertToMap(strLabelSelector)
    fieldSelector := util.SelectorConvertToMap(strFieldSelector)
    rtype := strings.ToUpper(vars["type"])
    if watcher, ok := eventWatchers[rtype]; ok {
        log.Info("rtype:", rtype)
        watcher.Watch(conn, handler.Config, labelSelector, fieldSelector)
    } else {
        log.Error("can not found watcher type:", rtype)
    }
}
