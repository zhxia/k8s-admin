package api

import (
	"kube-admin/src/util"
	"fmt"
	"k8s.io/apimachinery/pkg/fields"
	"testing"
)

func TestFieldSelector(t *testing.T) {
	selector := map[string]string{"name": "zhxia", "age": "18"}
	out := fields.SelectorFromSet(selector)
	fmt.Println(out)
}

func TestConvertToMap(t *testing.T) {
	str := "name=zhxia,spec.name=redis"
	out := util.SelectorConvertToMap(str)
	fmt.Println(out)

}
