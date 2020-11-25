package util

import (
    "github.com/gorilla/mux"
    "github.com/tidwall/gjson"
    "io/ioutil"
    "net/http"
    "strings"
)

type HttpRequest struct {
    Request *http.Request
}

func (r *HttpRequest) GetQuery(args ...string) (value string) {
    if len(args) == 0 {
        value = ""
        return
    }
    name := args[0]
    defVal := ""
    if len(args) > 1 {
        defVal = args[1]
    }
    value = strings.Trim(r.Request.URL.Query().Get(name), "")
    if value == "" && defVal != "" {
        value = defVal
    }
    return
}

func (r *HttpRequest) GetForm(name, defVal string) (value string) {
    value = strings.Trim(r.Request.PostFormValue(name), "")
    if value == "" && defVal != "" {
        value = defVal
    }
    return
}

func (r *HttpRequest) GetJson() (value gjson.Result) {
    data, err := ioutil.ReadAll(r.Request.Body)
    if err != nil {
        return gjson.Result{}
    }
    value = gjson.Parse(string(data))
    return
}

func (r *HttpRequest) GetPath(args ...string) (value string) {
    vars := mux.Vars(r.Request)
    if len(args) == 0 {
        value = ""
        return
    }
    name := args[0]
    if v, ok := vars[name]; !ok {
        if len(args) > 1 {
            value = args[1]
        } else {
            value = ""
        }
    } else {
        value = v
    }
    return
}
