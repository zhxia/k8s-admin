package api

import (
	"bytes"
	"context"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	restclient "k8s.io/client-go/rest"
)

type PodApi struct {
	k8sOperator *K8sOperator
}

func (api *PodApi) GetPodList(labelSelector, fieldSelector map[string]string, limit int64) ([]v1.Pod, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSelector).String(),
		FieldSelector: fields.SelectorFromSet(fieldSelector).String(),
		Limit:         limit,
	}
	pods, err := api.k8sOperator.GetPod().List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

func (api *PodApi) GetPod(name string) (*v1.Pod, error) {
	return api.k8sOperator.GetPod().Get(context.TODO(), name, metav1.GetOptions{})
}

func (api *PodApi) GetPodLogs(name, container string, lines int64) (string, error) {
	if lines == 0 {
		lines = 1000
	}
	podLogOptions := v1.PodLogOptions{
		TailLines: &lines,
	}
	if container != "" {
		podLogOptions.Container = container
	}
	req := api.k8sOperator.GetPod().GetLogs(name, &podLogOptions)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "", err
	}
	defer func() {
		podLogs.Close()
	}()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (api *PodApi) GetPodLogsRequest(name, container string) *restclient.Request {
	tailLines := int64(500)
	podLogOptions := v1.PodLogOptions{
		Follow:    true,
		TailLines: &tailLines,
	}
	if container != "" {
		podLogOptions.Container = container
	}
	return api.k8sOperator.GetPod().GetLogs(name, &podLogOptions)
}

func (api *PodApi) Watch(labelSelector, fieldSelector map[string]string) (watch.Interface, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSelector).String(),
		FieldSelector: fields.SelectorFromSet(fieldSelector).String(),
		Watch:         true,
	}
	return api.k8sOperator.GetPod().Watch(context.TODO(),listOptions)
}

func NewPodApi(namespace, kubeConfig string) *PodApi {
	return &PodApi{
		k8sOperator: NewOperator(kubeConfig, namespace),
	}
}
