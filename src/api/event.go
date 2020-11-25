package api

import (
    "context"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/fields"
    "k8s.io/apimachinery/pkg/labels"
    "k8s.io/apimachinery/pkg/watch"
)

type EventApi struct {
    operator *K8sOperator
}

func NewEventApi(namespace, kubeConfig string) *EventApi {
    return &EventApi{operator: NewOperator(kubeConfig, namespace)}
}

func (api *EventApi) Watch(labelSelector, fieldSelector map[string]string) (watch.Interface, error) {
    listOptions := metav1.ListOptions{
        LabelSelector: labels.SelectorFromSet(labelSelector).String(),
        FieldSelector: fields.SelectorFromSet(fieldSelector).String(),
        Watch:         true,
    }
    return api.operator.GetEvent().Watch(context.TODO(), listOptions)
}
