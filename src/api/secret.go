package api

import (
	"kube-admin/src/util"
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

type SecretApi struct {
	k8sOperator *K8sOperator
}

func (secretApi *SecretApi) GetSecretList(labelSelector, FieldSelector map[string]string, limit int64) ([]v1.Secret, error) {
	listOptions := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelSelector).String(),
		FieldSelector: fields.SelectorFromSet(FieldSelector).String(),
		Limit:         limit,
	}
	ss, err := secretApi.k8sOperator.GetSecret().List(context.Background(), listOptions)
	if err != nil {
		return nil, err
	}
	return ss.Items, nil
}

func (secretApi *SecretApi) CreateSecret(strYaml string) (*v1.Secret, error) {
	apiObj, err := util.YamlStrToApiObject(strYaml)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("get secret api object failed:%s", err.Error()))
	}
	secret := apiObj.(*v1.Secret)
	return secretApi.k8sOperator.GetSecret().Create(context.TODO(), secret, metav1.CreateOptions{})
}

func (secretApi *SecretApi) UpdateSecret(strYaml string) (*v1.Secret, error) {
	apiObj, err := util.YamlStrToApiObject(strYaml)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("convert to secret api object failed:%s", err.Error()))
	}
	secret := apiObj.(*v1.Secret)
	sec, err := secretApi.k8sOperator.GetSecret().Get(context.TODO(), secret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secretApi.k8sOperator.GetSecret().Update(context.TODO(), sec, metav1.UpdateOptions{})
}

func (secretApi *SecretApi) DeleteSecret(name string) error {
	return secretApi.k8sOperator.GetSecret().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func NewSecretApi(namespace, kubeConfig string) *SecretApi {
	return &SecretApi{
		k8sOperator: NewOperator(kubeConfig, namespace),
	}
}
