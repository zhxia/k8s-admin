package httpprocessor

import (
	"kube-admin/src/api"
	"kube-admin/src/util"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
)

var websshUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSshHandler struct {
	Config *util.Config
}

func (handler WebSshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websshUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("websocket upgrade failed:", err.Error())
		return
	}
	log.Info("client [", conn.RemoteAddr().String(), "] connected!")
	ctx, cancel := context.WithCancel(context.Background())
	conn.SetCloseHandler(func(code int, text string) error {
		log.Info("code:", code, "text:", text, ",client side connection closed!")
		cancel()
		return nil
	})
	namespace := r.URL.Query().Get("ns")
	pod := r.URL.Query().Get("pod")
	container := r.URL.Query().Get("c")
	log.Info(fmt.Sprintf("ns:%s,pod:%s,container:%s", namespace, pod, container))
	t := &api.WebSshTerminator{
		Context:   ctx,
		Namespace: namespace,
		Pod:       pod,
		Container: container,
		Conn:      conn,
		SizeChan:  make(chan *remotecommand.TerminalSize),
	}
	webSshApi := api.NewWebSshApi(namespace, handler.Config.K8sConfig)
	webSshApi.Handle(t, []string{"/bin/bash"})
}
