package api

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/gorilla/websocket"
    log "github.com/sirupsen/logrus"
    v1 "k8s.io/api/core/v1"
    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/tools/remotecommand"
    "net/http"
    "strconv"
    "strings"
)

type WebSshTerminator struct {
    Context   context.Context
    Namespace string
    Pod       string
    Container string
    SizeChan  chan *remotecommand.TerminalSize
    Conn      *websocket.Conn
}

func (wt *WebSshTerminator) Read(p []byte) (int, error) {
    log.Info("read message!")
    var reply []byte
    var msg map[string]string
    _, reply, err := wt.Conn.ReadMessage()
    if err != nil {
        return 0, err
    }
    log.Info("received:", string(reply))
    if err := json.Unmarshal(reply, &msg); err != nil {
        return 0, nil
    }
    if _, ok := msg["resize"]; ok {
        dataArr := strings.SplitN(msg["resize"], ":", 2)
        w, err := strconv.Atoi(dataArr[0])
        if err != nil {
            return 0, nil
        }
        width := uint16(w)
        h, err := strconv.Atoi(dataArr[1])
        if err != nil {
            return 0, nil
        }
        height := uint16(h)
        wt.SizeChan <- &remotecommand.TerminalSize{
            Width:  width,
            Height: height,
        }
        return 0, nil
    } else if _, ok := msg["data"]; ok {
        return copy(p, msg["data"]), nil
    } else {
        log.Error("invalid data format!")
        return 0, nil
    }
}

func (wt *WebSshTerminator) Write(p []byte) (int, error) {
    log.Info("write message!")
    err := wt.Conn.WriteMessage(websocket.BinaryMessage, p)
    log.Info("write:", string(p))
    return len(p), err
}

// 实现tty size queue
func (wt *WebSshTerminator) Next() *remotecommand.TerminalSize {
    size := <-wt.SizeChan
    log.Info(fmt.Sprintf("terminal size to width: %d height: %d", size.Width, size.Height))
    return size
}

//===========WebSshApi==============
type WebSshApi struct {
    k8sOperator *K8sOperator
}

func NewWebSshApi(namespace, kubeConfig string) *WebSshApi {
    return &WebSshApi{k8sOperator: NewOperator(kubeConfig, namespace)}
}

func (webssh *WebSshApi) Handle(t *WebSshTerminator, cmd []string) error {
    fn := func() error {
        req := webssh.k8sOperator.client.CoreV1().RESTClient().Post().Resource("pods").Name(t.Pod).
            Namespace(t.Namespace).SubResource("exec")
        req.VersionedParams(&v1.PodExecOptions{
            Container: t.Container,
            Stdin:     true,
            Stdout:    true,
            Stderr:    true,
            Command:   cmd,
            TTY:       true,
        }, scheme.ParameterCodec)
        log.Info("req.URL:", req.URL().String())
        executor, err := remotecommand.NewSPDYExecutor(webssh.k8sOperator.config, http.MethodPost, req.URL())
        if err != nil {
            log.Error(err.Error())
            return err
        }
        log.Info("execute")
        return executor.Stream(remotecommand.StreamOptions{
            Stdin:             t,
            Stdout:            t,
            Stderr:            t,
            Tty:               true,
            TerminalSizeQueue: t,
        })
    }
    return fn()
    //log.Info("handle:")
    //inFd, ok := term.GetFdInfo(t.Conn)
    //log.Info("isTerminalIn:", ok)
    //state, err := term.SaveState(inFd)
    //if err != nil {
    //	log.Error(err.Error())
    //	return err
    //}
    //return interrupt.Chain(nil, func() {
    //	term.RestoreTerminal(inFd, state)
    //}).Run(fn)
}
