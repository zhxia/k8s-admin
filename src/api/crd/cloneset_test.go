package crd

import (
    "fmt"
    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
    "kube-admin/src/util"
    "testing"
)

func TestGetCloneSet(t *testing.T) {
    csApi := NewCloneSetApi("", "")
    cs, err := csApi.GetCloneSet("sample")
    if err != nil {
        fmt.Println("err:", err.Error())
    }
    ct := cs.Spec.Template.Spec.Containers[0]
    fmt.Println("result:", ct.Resources.Requests.Memory().String())
}

func TestGetCloneSetList(t *testing.T) {
    csApi := NewCloneSetApi("", "")
    labelSelector := make(map[string]string)
    labelSelector["release"] = "v2"
    fieldSelector := make(map[string]string)
    fieldSelector["metadata.name"] = "sample"
    csl, err := csApi.GetCloneSetList(nil, fieldSelector, 100)
    if err != nil {
        fmt.Println("err:", err.Error())
    }
    for _, cs := range csl {
        fmt.Println("out:", cs)
    }
}

func TestCreateCloneSet(t *testing.T) {
    strYaml := `apiVersion: apps.kruise.io/v1alpha1
kind: CloneSet
metadata:
  labels:
    app: sample
  name: sample
  namespace: test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: sample
  template:
    metadata:
      labels:
        app: sample
    spec:
      containers:
        - name: nginx
          image: nginx:alpine`
    csApi := NewCloneSetApi("", "")
    obj, err := csApi.CreateCloneSet(strYaml)
    if err != nil {
        fmt.Println("error:", err)
        return
    }
    fmt.Println(*(*obj).Spec.Replicas, (*obj).Kind)
}

func TestUpdateCloneSet(t *testing.T) {
    strYaml := `apiVersion: apps.kruise.io/v1alpha1
kind: CloneSet
metadata:
  labels:
    app: sample
    release: v1
  name: sample
  namespace: default
spec:
  replicas: 3
  updateStrategy:
    partition: 0
    maxSurge: 10%
    type: InPlaceIfPossible
    inPlaceUpdateStrategy:
        gracePeriodSeconds: 10
  selector:
    matchLabels:
      app: sample
  template:
    metadata:
      labels:
        app: sample
    spec:
      containers:
        - name: nginx
          image: nginx:mainline`
    csApi := NewCloneSetApi("", "")
    obj, err := csApi.UpdateCloneSet(strYaml)
    if err != nil {
        fmt.Println("error:", err)
        return
    }
    fmt.Println(*(*obj).Spec.Replicas, (*obj).Kind)
}

func TestCloneSetUpdate(t *testing.T) {
    csApi := NewCloneSetApi("", "")
    var rr = v1.ResourceRequirements{
        Limits: map[v1.ResourceName]resource.Quantity{
           v1.ResourceName("memory"): resource.MustParse("200Mi"),
           v1.ResourceName("cpu"):    resource.MustParse("2m"),
        },
        //Requests: map[v1.ResourceName]resource.Quantity{
        //    v1.ResourceName("memory"): resource.MustParse("55Mi"),
        //    v1.ResourceName("cpu"):    resource.MustParse("2m"),
        //},
    }
    cs, e := csApi.Update("sample", "nginx:mainline", util.Int32Ptr(2), util.Int32Ptr(1), &rr, false)
    if e != nil {
        fmt.Println(e)
        return
    }
    fmt.Println(*cs.Spec.Replicas)
}
