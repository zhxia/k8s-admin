package api

import (
	"context"
	appsv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

type ConfigMapApi struct {
	k8sOperator *K8sOperator
}

func (configMapApi *ConfigMapApi) GetConfigMapList(labelSelector, fieldSelector map[string]string, limit int64) ([]appsv1.ConfigMap, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSelector).String(),
		FieldSelector: fields.SelectorFromSet(fieldSelector).String(),
		Limit:         limit,
	}
	cms, err := configMapApi.k8sOperator.GetConfigMap().List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return cms.Items, nil
}

func (configMapApi *ConfigMapApi) GetConfigMap(name string) (*appsv1.ConfigMap, error) {
	return configMapApi.k8sOperator.GetConfigMap().Get(context.TODO(), name, metav1.GetOptions{})
}

func NewConfigMapApi(namespace, kubeConfig string) *ConfigMapApi {
	return &ConfigMapApi{
		k8sOperator: NewOperator(kubeConfig, namespace),
	}
}
