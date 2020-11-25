package api

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetScale(t *testing.T) {
	scale, err := NewDeploymentApi("default", "").GetScale("redis-v1")
	if err != nil {
		panic(err)
	}
	data, _ := json.Marshal(scale.Spec)
	fmt.Println(string(data))
}

func TestUpdateScale(t *testing.T) {
	scale, err := NewDeploymentApi("default", "").UpdateScale("redis-v1", 5)
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(scale)
	fmt.Println(string(data))

}
