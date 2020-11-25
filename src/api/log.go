package api

import (
	"bytes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
)

type LogApi struct {
	k8sOperator *K8sOperator
}

type LogWriter struct {
	buffer bytes.Buffer
}

func (writer *LogWriter) Write(p []byte) (int, error) {
	log.Info("write data:", string(p))
	return writer.buffer.Write(p)
}

func (writer *LogWriter) Read(p []byte) (int, error) {
	var out []byte
	out = writer.buffer.Bytes()
	log.Info("read data:", string(out))
	if len(out) > 0 {
		writer.buffer.Reset()
		return len(out), nil
	} else {
		return 0, io.EOF
	}
}

func NewLogApi(namespace, kubeConfig string) *LogApi {
	return &LogApi{k8sOperator: NewOperator(kubeConfig, namespace)}
}

func (api *LogApi) ExecRemoteCmd(namespace, pod, container string, cmd []string) ([]byte, error) {
	req := api.k8sOperator.client.CoreV1().RESTClient().Post().Resource("pods").Name(pod).
		Namespace(namespace).SubResource("exec")
	req.VersionedParams(&v1.PodExecOptions{
		Container: container,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		Command:   cmd,
		TTY:       false,
	}, scheme.ParameterCodec)
	log.Debug("req.URL:", req.URL().String())
	executor, err := remotecommand.NewSPDYExecutor(api.k8sOperator.config, http.MethodPost, req.URL())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var stdout, stderr bytes.Buffer
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return nil, errors.Wrap(err, string(stderr.Bytes()))
	}
	return stdout.Bytes(), nil
}
