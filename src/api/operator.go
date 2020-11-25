package api

import (
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    apiv1 "k8s.io/api/core/v1"
    "k8s.io/client-go/kubernetes"
    appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
    corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
    extensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
    restclient "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    "kube-admin/src/util"
    "os"
    "path/filepath"
)

type K8sOperator struct {
    client    *kubernetes.Clientset
    namespace string
    config    *restclient.Config
}

func (operator *K8sOperator) GetPersistentVolume() corev1.PersistentVolumeInterface {
    return operator.client.CoreV1().PersistentVolumes()
}

func (operator *K8sOperator) GetServiceAccount() corev1.ServiceAccountInterface {
    return operator.client.CoreV1().ServiceAccounts(operator.namespace)
}

func (operator *K8sOperator) GetSecret() corev1.SecretInterface {
    return operator.client.CoreV1().Secrets(operator.namespace)
}
func (operator *K8sOperator) GetConfigMap() corev1.ConfigMapInterface {
    return operator.client.CoreV1().ConfigMaps(operator.namespace)
}

func (operator *K8sOperator) GetNode() corev1.NodeInterface {
    return operator.client.CoreV1().Nodes()
}

func (operator *K8sOperator) GetNamespace() corev1.NamespaceInterface {
    return operator.client.CoreV1().Namespaces()
}

func (operator *K8sOperator) GetService() corev1.ServiceInterface {
    return operator.client.CoreV1().Services(operator.namespace)
}

func (operator *K8sOperator) GetDeployment() appsv1.DeploymentInterface {
    return operator.client.AppsV1().Deployments(operator.namespace)
}

func (operator *K8sOperator) GetPod() corev1.PodInterface {
    return operator.client.CoreV1().Pods(operator.namespace)
}

func (operator *K8sOperator) GetIngress() extensionsv1beta1.IngressInterface {
    return operator.client.ExtensionsV1beta1().Ingresses(operator.namespace)
}

func (operator *K8sOperator) GetEvent() corev1.EventInterface {
    return operator.client.CoreV1().Events(operator.namespace)
}

func NewOperator(kubeConfig, namespace string) *K8sOperator {
    var config *restclient.Config
    var err error
    if namespace == "" {
        namespace = apiv1.NamespaceDefault
    }
    if kubeConfig == "incluster" {
        config, err = restclient.InClusterConfig() //incluster mode
    } else {
        if kubeConfig == "" {
            home := homedir.HomeDir()
            kubeConfig = filepath.Join(home, ".kube", "config")
        }
        exist, _ := util.PathExists(kubeConfig)
        if !exist {
            log.Error("kubectl config file:\"", kubeConfig, "\" not found!")
            os.Exit(0)
        }
        config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
    }
    if err != nil {
        panic(errors.Wrap(err, "get k8s config failed"))
    }
    client, err := kubernetes.NewForConfig(config)
    return &K8sOperator{
        client:    client,
        namespace: namespace,
        config:    config,
    }
}
