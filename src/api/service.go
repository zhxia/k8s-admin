package api

import (
	"context"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

type ServiceApi struct {
	k8sOperator *K8sOperator
}

func (secretApi *ServiceApi) GetServiceList(labelSelector, fieldSelector map[string]string, limit int64) ([]v1.Service, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSelector).String(),
		FieldSelector: fields.SelectorFromSet(fieldSelector).String(),
		Limit:         500,
	}
	servs, err := secretApi.k8sOperator.GetService().List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return servs.Items, nil
}

func (secretApi *ServiceApi) GetService(name string) (*v1.Service, error) {
	return secretApi.k8sOperator.GetService().Get(context.TODO(), name, metav1.GetOptions{})
}

func NewServiceApi(namespace, kubeConfig string) *ServiceApi {
	return &ServiceApi{
		k8sOperator: NewOperator(kubeConfig, namespace),
	}
}
