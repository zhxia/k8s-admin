package util

import (
    "bufio"
    kruiseapps "github.com/openkruise/kruise-api/apps/v1alpha1"
    "github.com/pkg/errors"
    "github.com/tidwall/gjson"
    "io/ioutil"
    "k8s.io/apimachinery/pkg/util/yaml"
    "k8s.io/client-go/kubernetes/scheme"
    "net/http"
    "os"
    "strings"
)

func YamlStrToApiObject(strYaml string) (interface{}, error) {
    jsonData, err := yaml.ToJSON([]byte(strYaml))
    if err != nil {
        return nil, errors.Wrap(err, "str yaml convert to json failed")
    }
    // CRD scheme added here
    kruiseapps.AddToScheme(scheme.Scheme)
    decode := scheme.Codecs.UniversalDeserializer().Decode
    apiObject, _, err := decode([]byte(jsonData), nil, nil)
    if err != nil {
        return nil, errors.Wrap(err, "scheme.codecs.decode failed")
    }
    return apiObject, err
}

func YamlFileToJson(filename string) gjson.Result {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }
    jsonData, err := yaml.ToJSON(data)
    if err != nil {
        panic(err)
    }
    return gjson.Parse(string(jsonData))
}

func Int32Ptr(i int32) *int32 {
    return &i
}

func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

func SelectorConvertToMap(strSelector string) map[string]string {
    strSelector = strings.Trim(strSelector, "")
    if strSelector == "" || strSelector == "*" {
        return nil
    }
    mp := make(map[string]string)
    parts := strings.Split(strSelector, ",")
    for _, v := range parts {
        arr := strings.Split(v, "=")
        mp[arr[0]] = arr[1]
    }
    return mp
}

func GetJson(r *http.Request) gjson.Result {
    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return gjson.Result{}
    }
    result := gjson.Parse(string(data))
    return result
}

func ReadLine(r *bufio.Reader) ([]byte, error) {
    var (
        err      error = nil
        isPrefix       = true
        line, ln []byte
    )
    for isPrefix && err == nil {
        line, isPrefix, err = r.ReadLine()
        if err != nil {
            return nil, err
        }
        ln = append(ln, line...)
        if len(ln) > 8*1024 {
            return nil, err
        }
    }
    return ln, err
}
