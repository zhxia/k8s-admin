package api

import (
    "kube-admin/src/util"
    "context"
    "errors"
    "fmt"
    appsv1 "k8s.io/api/apps/v1"
    autoscalingv1 "k8s.io/api/autoscaling/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/fields"
    "k8s.io/apimachinery/pkg/labels"
    "k8s.io/apimachinery/pkg/watch"
)

type DeploymentApi struct {
    k8sOperator *K8sOperator
}

func (depApi *DeploymentApi) GetDeploymentList(labelSelector, fieldSelector map[string]string, limit int64) ([]appsv1.Deployment, error) {
    listOptions := metav1.ListOptions{
        LabelSelector: labels.SelectorFromSet(labelSelector).String(),
        FieldSelector: fields.SelectorFromSet(fieldSelector).String(),
    }
    deps, err := depApi.k8sOperator.GetDeployment().List(context.TODO(), listOptions)
    if err != nil {
        return nil, err
    }
    return deps.Items, nil
}

func (depApi *DeploymentApi) GetDeployment(name string) (*appsv1.Deployment, error) {
    return depApi.k8sOperator.GetDeployment().Get(context.TODO(), name, metav1.GetOptions{})
}

func (depApi *DeploymentApi) CreateDeployment(strYaml string) (*appsv1.Deployment, error) {
    apiObj, err := util.YamlStrToApiObject(strYaml)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("get deployment api object failed:%s", err.Error()))
    }
    deployment := apiObj.(*appsv1.Deployment)
    dep, err := depApi.k8sOperator.GetDeployment().Create(context.TODO(), deployment, metav1.CreateOptions{})
    if err != nil {
        return nil, errors.New(fmt.Sprintf("create deployment api object failed:%s", err.Error()))
    }
    return dep, nil
}

func (depApi *DeploymentApi) UpdateDeployment(strYaml string) (*appsv1.Deployment, error) {
    apiObj, err := util.YamlStrToApiObject(strYaml)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("get deployment api object failed:%s", err.Error()))
    }
    deployment := apiObj.(*appsv1.Deployment)
    dep, err := depApi.k8sOperator.GetDeployment().Get(context.TODO(), deployment.Name, metav1.GetOptions{})
    if err != nil {
        return nil, errors.New(fmt.Sprintf("get old deployment api object:[%s] failed:%s, ", deployment.Name, err.Error()))
    }
    deployment.ResourceVersion = dep.ResourceVersion
    dep, err = depApi.k8sOperator.GetDeployment().Update(context.TODO(), deployment, metav1.UpdateOptions{})
    if err != nil {
        return nil, errors.New(fmt.Sprintf("update deployment api object failed:%s", err.Error()))
    }
    return dep, nil
}

func (depApi *DeploymentApi) UpdateScale(deploymentName string, replicas int32) (*autoscalingv1.Scale, error) {
    scale, err := depApi.GetScale(deploymentName)
    if err != nil {
        return nil, err
    }
    scale.Spec.Replicas = replicas
    updateOptions := metav1.UpdateOptions{}
    newScale, err := depApi.k8sOperator.GetDeployment().UpdateScale(context.TODO(), deploymentName, scale, updateOptions)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("update deployment scale failed:%s", err.Error()))
    }
    return newScale, nil
}

func (depApi *DeploymentApi) GetScale(deploymentName string) (*autoscalingv1.Scale, error) {
    scale, err := depApi.k8sOperator.GetDeployment().GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
    if err != nil {
        return nil, errors.New(fmt.Sprintf("get deployment scale failed:%s", err.Error()))
    }
    return scale, nil
}

func (depApi *DeploymentApi) Watch(labelSelector, fieldSelector map[string]string) (watch.Interface, error) {
    listOptions := metav1.ListOptions{
        LabelSelector: labels.SelectorFromSet(labelSelector).String(),
        FieldSelector: fields.SelectorFromSet(fieldSelector).String(),
        Watch:         true,
    }
    return depApi.k8sOperator.GetDeployment().Watch(context.TODO(), listOptions)
}

func NewDeploymentApi(namespace, kubeConfig string) *DeploymentApi {
    deploymentApi := new(DeploymentApi)
    deploymentApi.k8sOperator = NewOperator(kubeConfig, namespace)
    return deploymentApi
}
