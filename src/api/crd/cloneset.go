package crd

import (
    "context"
    "fmt"
    kruiseapps "github.com/openkruise/kruise-api/apps/v1alpha1"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    apiv1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/equality"
    "k8s.io/apimachinery/pkg/fields"
    "k8s.io/apimachinery/pkg/labels"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
    kscheme "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    "kube-admin/src/util"
    "os"
    "path/filepath"
)
import "sigs.k8s.io/controller-runtime/pkg/client"

type CloneSetApi struct {
    namespace  string
    kubeConfig string
}

func NewCloneSetApi(namespace, kubeConfig string) *CloneSetApi {
    return &CloneSetApi{
        namespace:  namespace,
        kubeConfig: kubeConfig,
    }
}
func (cloneSetApi *CloneSetApi) getClient() (client.Client, error) {
    if cloneSetApi.namespace == "" {
        cloneSetApi.namespace = apiv1.NamespaceDefault
    }
    if cloneSetApi.kubeConfig == "" {
        home := homedir.HomeDir()
        cloneSetApi.kubeConfig = filepath.Join(home, ".kube", "config")
    }
    exist, _ := util.PathExists(cloneSetApi.kubeConfig)
    if !exist {
        log.Error("kubectl config file:\"", cloneSetApi.kubeConfig, "\" not found!")
        os.Exit(0)
    }
    config, err := clientcmd.BuildConfigFromFlags("", cloneSetApi.kubeConfig)
    if err != nil {
        panic(err)
    }
    scheme := runtime.NewScheme()
    kruiseapps.AddToScheme(scheme)
    kscheme.AddToScheme(scheme)
    return client.New(config, client.Options{Scheme: scheme})
}

func (cloneSetApi *CloneSetApi) GetCloneSetList(labelSelector, fieldSelector map[string]string, limit int64) ([]kruiseapps.CloneSet, error) {
    cloneSetList := kruiseapps.CloneSetList{}
    cli, err := cloneSetApi.getClient()
    if err != nil {
        panic(err)
    }
    listOptions := client.ListOptions{
        LabelSelector: labels.SelectorFromSet(labelSelector),
        FieldSelector: fields.SelectorFromSet(fieldSelector),
        Limit:         limit,
    }
    err = cli.List(context.TODO(), &cloneSetList, &listOptions, client.InNamespace(cloneSetApi.namespace))
    if err != nil {
        return nil, err
    }
    return cloneSetList.Items, nil
}

func (cloneSetApi *CloneSetApi) GetCloneSet(name string) (kruiseapps.CloneSet, error) {
    cloneSet := kruiseapps.CloneSet{}
    cli, err := cloneSetApi.getClient()
    if err != nil {
        panic(err)
    }
    err = cli.Get(context.TODO(), types.NamespacedName{Namespace: cloneSetApi.namespace, Name: name}, &cloneSet)
    return cloneSet, err
}

func (cloneSetApi *CloneSetApi) CreateCloneSet(strYaml string) (*kruiseapps.CloneSet, error) {
    apiObj, err := util.YamlStrToApiObject(strYaml)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("get cloneset api object failed:%s", err.Error()))
    }
    cloneSet := apiObj.(*kruiseapps.CloneSet)
    cli, err := cloneSetApi.getClient()
    err = cli.Create(context.TODO(), cloneSet)
    if err != nil {
        return nil, err
    }
    return cloneSet, nil
}

func (cloneSetApi *CloneSetApi) CreateOrUpdate(strYaml string) (*kruiseapps.CloneSet, error) {
    apiObj, err := util.YamlStrToApiObject(strYaml)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("get cloneset api object failed:%s", err.Error()))
    }
    cloneSet := apiObj.(*kruiseapps.CloneSet)
    name := cloneSet.Name
    cli, err := cloneSetApi.getClient()
    oldCloneSet, err := cloneSetApi.GetCloneSet(name)
    if err != nil {
        log.Info(fmt.Sprintf("get old cloneset failed:%s,need to create new one", name))
        err = cli.Create(context.TODO(), cloneSet)
        if err != nil {
            log.Error(fmt.Sprintf("create new cloneset failed:%s", name))
            return nil, err
        }
        return cloneSet, nil
    } else {
        cloneSet.ResourceVersion = oldCloneSet.ResourceVersion
        err = cli.Update(context.TODO(), cloneSet)
        if err != nil {
            log.Error(fmt.Sprintf("update cloneset:%s failed,error:%s", name, err.Error()))
            return nil, err
        }
        return cloneSet, nil
    }
}

func (cloneSetApi *CloneSetApi) UpdateCloneSet(strYaml string) (*kruiseapps.CloneSet, error) {
    apiObj, err := util.YamlStrToApiObject(strYaml)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("get cloneset api object failed:%s", err.Error()))
    }
    cloneSet := apiObj.(*kruiseapps.CloneSet)
    name := cloneSet.Name
    oldCloneSet, err := cloneSetApi.GetCloneSet(name)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("get old cloneset failed:%s", oldCloneSet.Name))
    }
    // update cloneset need old resource version
    cloneSet.ResourceVersion = oldCloneSet.ResourceVersion
    cli, err := cloneSetApi.getClient()
    err = cli.Update(context.TODO(), cloneSet)
    if err != nil {
        return nil, err
    }
    return cloneSet, nil
}

func (cloneSetApi *CloneSetApi) Update(name, image string, replicas, partition *int32, res *apiv1.ResourceRequirements, paused bool) (*kruiseapps.CloneSet, error) {
    cloneSet, err := cloneSetApi.GetCloneSet(name)
    if err != nil {
        return nil, errors.Wrap(err, fmt.Sprintf("get old cloneset:%s failed", cloneSet.Name))
    }
    if image != "" {
        cloneSet.Spec.Template.Spec.Containers[0].Image = image
    }
    if replicas != nil {
        cloneSet.Spec.Replicas = replicas
    }
    if partition != nil {
        cloneSet.Spec.UpdateStrategy.Partition = partition
    }
    cloneSet.Spec.UpdateStrategy.Paused = paused
    if res != nil {
        // old resource value override by new resource value
        oldResource := cloneSet.Spec.Template.Spec.Containers[0].Resources
        if res.Requests != nil {
            oldResource.Requests = res.Requests
        }
        if res.Limits != nil {
            oldResource.Limits = res.Limits
        }
        cloneSet.Spec.Template.Spec.Containers[0].Resources = oldResource
    }
    cli, err := cloneSetApi.getClient()
    err = cli.Update(context.TODO(), &cloneSet)
    if err != nil {
        return nil, errors.Wrap(err, "update cloneset failed")
    }
    return &cloneSet, nil
}

func (cloneSetApi *CloneSetApi) Compare(csNew, csOld kruiseapps.CloneSet) bool {
    return equality.Semantic.DeepDerivative(csNew, csOld)
}
