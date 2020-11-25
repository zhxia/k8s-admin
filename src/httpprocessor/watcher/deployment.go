package watcher

import (
    "github.com/gorilla/websocket"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    appsv1 "k8s.io/api/apps/v1"
    "k8s.io/apimachinery/pkg/util/json"
    "k8s.io/apimachinery/pkg/watch"
    "kube-admin/src/api"
    "kube-admin/src/util"
    "time"
)

type DevelopmentWatcher struct {
}

func (w DevelopmentWatcher) Watch(conn *websocket.Conn, cfg *util.Config, labelSelector, fieldSelector map[string]string) {
    watcher, err := api.NewDeploymentApi(cfg.K8sNamespace, cfg.K8sConfig).Watch(labelSelector, fieldSelector)
    if err != nil {
        log.Error(errors.Wrap(err, "get deployment watcher failed"))
        return
    }
    go w.reader(conn)
    w.writer(conn, watcher)
}

func (w DevelopmentWatcher) reader(conn *websocket.Conn) {
    conn.SetReadDeadline(time.Now().Add(pongWait))
    conn.SetPongHandler(func(pong string) error {
        // response to other peer's pong
        log.Info("pong:", pong)
        conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Error(errors.Wrap(err, "connection error"))
            }
            break
        }
    }
}

func (w DevelopmentWatcher) writer(conn *websocket.Conn, watcher watch.Interface) {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        watcher.Stop()
    }()
    for {
        select {
        case <-ticker.C:
            log.Info("send ping...")
            conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                log.Error(errors.Wrap(err, "ping error"))
                return
            }
        case e := <-watcher.ResultChan():
            if e.Object == nil {
                continue
            }
            log.Info("object:", e.Object)
            event, ok := e.Object.(*appsv1.Deployment)
            if !ok {
                continue
            }
            data, err := json.Marshal(event)
            if err != nil {
                log.Error(errors.Wrap(err, "event object json encoded failed!"))
                continue
            }
            conn.WriteJSON(event)
            log.Info("event:", string(data))
        }
    }
}
