package watcher

import (
    "github.com/gorilla/websocket"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/util/json"
    "k8s.io/apimachinery/pkg/watch"
    "kube-admin/src/api"
    "kube-admin/src/util"
    "time"
)

type EventWatcher struct {
}

func (w EventWatcher) Watch(conn *websocket.Conn, cfg *util.Config, labelSelector, fieldSelector map[string]string) {
    watcher, err := api.NewEventApi(cfg.K8sNamespace, cfg.K8sConfig).Watch(labelSelector, fieldSelector)
    if err != nil {
        log.Error(errors.Wrap(err, "get event watcher failed"))
        return
    }
    go w.reader(conn)
    w.writer(conn, watcher)
}

func (w EventWatcher) reader(conn *websocket.Conn) {
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

func (w EventWatcher) writer(conn *websocket.Conn, watcher watch.Interface) {
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
            event, ok := e.Object.(*v1.Event)
            if !ok {
                continue
            }
            data, err := json.Marshal(event)
            if err != nil {
                log.Error(errors.Wrap(err, "event object json encoded failed!"))
                continue
            }
            conn.SetWriteDeadline(time.Now().Add(writeWait))
            conn.WriteJSON(event)
            log.Info("event:", string(data))
        }
    }
}
