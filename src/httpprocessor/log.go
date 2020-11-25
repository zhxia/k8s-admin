package httpprocessor

import (
	"kube-admin/src/api"
	"kube-admin/src/util"
	"bufio"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type LogHandler struct {
	Config *util.Config
}

func (handler LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logConnected := false
	ctx, cancel := context.WithCancel(context.Background())
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("websocket upgrade failed:", err.Error())
		return
	}
	log.Info("client [", conn.RemoteAddr().String(), "] connected!")
	defer func() {
		log.Info("client [", conn.RemoteAddr().String(), "] closed!")
		conn.Close()
	}()
	conn.SetCloseHandler(func(code int, text string) error {
		log.Info("code:", code, "text:", text, ",connection closed!")
		cancel()
		return nil
	})
	var strMsg = ""

	for {
		mt, msg, err := conn.ReadMessage()
		log.Info("message received!")
		if err != nil {
			log.Error("read message error:", err)
			return
		}
		strMsg = string(msg)
		log.Info("received command:", strMsg)
		if !strings.Contains(strMsg, "@") {
			err = errors.New("invalid command format! eg:pod@container@container")
			log.Error(err)
			conn.WriteMessage(mt, []byte(err.Error()))
			continue
		}
		// pod@container@namespace
		arrMsg := strings.SplitN(strMsg, "@", 3)
		namespace := handler.Config.K8sNamespace
		if len(arrMsg) == 3 {
			namespace = arrMsg[2]
		}
		podName := arrMsg[0]
		containerName := arrMsg[1]
		if !logConnected {
			log.Info(fmt.Sprintf("connected to :[%s-%s-%s]", namespace, podName, containerName))
			logConnected = true
			go func() {
				handler.readLog(ctx, conn, podName, containerName, namespace)
			}()
		}
	}

}

func (handler LogHandler) readLog(ctx context.Context, conn *websocket.Conn, podName, containerName, namespace string) {
	req := api.NewPodApi(namespace, handler.Config.K8sConfig).GetPodLogsRequest(podName, containerName)
	rs, err := req.Stream(ctx)
	if err != nil {
		log.Error("get pod stream failed!")
		return
	}
	closed := false
	defer func() {
		log.Info("pod logs stream closed!")
		rs.Close()
	}()
	go func() {
		<-ctx.Done()
		closed = true
	}()

	scanner := bufio.NewScanner(rs)
	for scanner.Scan() {
		if closed {
			log.Info("connection closed!")
			break
		}
		data := scanner.Bytes()
		log.Debug(string(data))
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Error("message send failed!")
			return
		}
	}
}
