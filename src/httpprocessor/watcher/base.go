package watcher

import (
    "kube-admin/src/util"
    "github.com/gorilla/websocket"
    "time"
)

var (
    writeWait  = 10 * time.Second
    pingPeriod = (pongWait * 9) / 10
    pongWait   = 60 * time.Second
)

type BaseWatcher interface {
    Watch(conn *websocket.Conn, cfg *util.Config, labelSelector, fieldSelector map[string]string)
}
