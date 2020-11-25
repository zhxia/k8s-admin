package api

import (
    "context"
    "github.com/pkg/errors"
    "k8s.io/api/extensions/v1beta1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "kube-admin/src/util"
)

type IngressApi struct {
    operator *K8sOperator
}

func (api *IngressApi) CreateOrUpdate(strYaml string) (ingress *v1beta1.Ingress, err error) {
    apiObj, err := util.YamlStrToApiObject(strYaml)
    if err != nil {
        err = errors.Wrap(err, "get ingress object failed")
        return
    }
    ingress = apiObj.(*v1beta1.Ingress)
    oldIngress, err := api.operator.GetIngress().Get(context.TODO(), ingress.Name, metav1.GetOptions{})
    if err != nil {
        ingress, err = api.operator.GetIngress().Create(context.TODO(), ingress, metav1.CreateOptions{})
    } else {
        ingress.ResourceVersion = oldIngress.ResourceVersion
        ingress, err = api.operator.GetIngress().Update(context.TODO(), ingress, metav1.UpdateOptions{})
    }
    return
}

func (api *IngressApi) GetIngress(name string) (ingress *v1beta1.Ingress, err error) {
    ingress, err = api.operator.GetIngress().Get(context.TODO(), name, metav1.GetOptions{})
    return
}
