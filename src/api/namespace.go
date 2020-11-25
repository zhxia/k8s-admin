package api

import (
    "context"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/fields"
    "k8s.io/apimachinery/pkg/labels"
)

type NamespaceApi struct {
    k8sOperator *K8sOperator
}

func NewNamespaceApi(namespace,kubeConfig string) *NamespaceApi {
    return &NamespaceApi{k8sOperator: NewOperator(kubeConfig,namespace)}
}

func (api *NamespaceApi) GetNamespaceList(labelSelector, fieldSelector map[string]string, limit int64) ([]v1.Namespace, error) {
    listOptions := metav1.ListOptions{
        LabelSelector: labels.SelectorFromSet(labelSelector).String(),
        FieldSelector: fields.SelectorFromSet(fieldSelector).String(),
    }
    nsList, err := api.k8sOperator.GetNamespace().List(context.TODO(), listOptions)
    if err != nil {
        return nil, err
    }
    return nsList.Items, nil

}
