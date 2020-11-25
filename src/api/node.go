package api

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

type NodeApi struct {
	operator *K8sOperator
}

func (nodeApi *NodeApi) GetNodeList(labelSelector, fieldSelector map[string]string, limit int64) ([]v1.Node, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSelector).String(),
		FieldSelector: fields.SelectorFromSet(fieldSelector).String(),
		Limit:         limit,
	}
	nodes, err := nodeApi.operator.GetNode().List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return nodes.Items, nil
}

func (nodeApi *NodeApi) GetNode(name string) (*v1.Node, error) {
	return nodeApi.operator.GetNode().Get(context.TODO(), name, metav1.GetOptions{})
}

func NewNodeApi(namespace, kubeConfig string) *NodeApi {
	return &NodeApi{
		operator: NewOperator(kubeConfig, namespace),
	}
}
