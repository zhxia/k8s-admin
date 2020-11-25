package main

import (
    "fmt"
    "github.com/gorilla/mux"
    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    log "github.com/sirupsen/logrus"
    nettools "github.com/toolkits/net"
    "gopkg.in/alecthomas/kingpin.v2"
    "io"
    "kube-admin/src/httpprocessor"
    "kube-admin/src/middleware"
    "kube-admin/src/util"
    "net"
    "net/http"
    "os"
    "os/exec"
    "path"
    "strings"
)

var (
    host     = kingpin.Arg("host", "Server bind host").Default("127.0.0.1").IP()
    port     = kingpin.Arg("port", "Server listen port.").Default("8487").Int()
    debug    = kingpin.Flag("debug", "Enable debug mode.").Bool()
    daemon   = kingpin.Flag("daemon", "Enable daemon mode.").Bool()
    confFile = kingpin.Flag("conf", "Set config file.").Default("./config.yml").ExistingFileOrDir()
    version  = "1.0.0"
)

var conf *util.Config

func init() {
    kingpin.Version(version)
    kingpin.Parse()
    if *daemon {
        args := os.Args[1:]
        i := 0
        argList := args[:0]
        for ; i < len(args); i++ {
            if args[i] == "--daemon" {
                args[i] = ""
                continue
            }
            argList = append(argList, args[i])
        }
        fmt.Println(os.Args[0], args)
        cmd := exec.Command(os.Args[0], argList...)
        cmd.Start()
        log.Info("[PID]:", cmd.Process.Pid)
        os.Exit(0)
    }
    conf = util.NewConfig(*confFile)
    log.SetLevel(conf.LogLevel)
    var writer io.Writer
    if *debug {
        writer = os.Stdout
    } else {
        var err error
        logFile := path.Join(conf.LogDir, "jr.log")
        writer, err = rotatelogs.New(
            logFile+".%Y%m%d",
            rotatelogs.WithLinkName(logFile),
            rotatelogs.WithRotationCount(10),
        )
        if err != nil {
            log.Fatalln("log file open failed!")
            os.Exit(0)
        }
    }
    log.SetOutput(writer)
    log.SetFormatter(&log.TextFormatter{
        DisableColors: false,
        FullTimestamp: true,
    })
    if host.Equal(net.ParseIP("127.0.0.1")) && conf.ServerHost != "" {
        *host = net.ParseIP(conf.ServerHost)
    }
    if *port == 8487 && conf.ServerPort != 0 {
        *port = conf.ServerPort
    }
    ips, _ := nettools.IntranetIP()
    log.Debug("server ip:", ips[0], ",host:", *host, ",port:", *port, ",debug:", *debug)
}

func main() {
    r := mux.NewRouter()
    r.Use(mux.CORSMethodMiddleware(r))
    r.HandleFunc("/", httpprocessor.DefaultHandler).Methods("GET")
    s := r.PathPrefix("/api").Subrouter()
    s.Use(middleware.LoggingMiddleware)
    am := middleware.AuthMiddleware{}
    am.Parse(conf.AuthUsers)
    s.Use(am.Middleware)
    apiDispatcher := httpprocessor.ApiHandlerDispatcher{
        Handlers: map[string]httpprocessor.BaseHandler{
            "DEPLOY":    {Config: conf, RequestProcessor: &httpprocessor.DeploymentProcessor{}},
            "POD":       {Config: conf, RequestProcessor: &httpprocessor.PodProcessor{}},
            "NODE":      {Config: conf, RequestProcessor: &httpprocessor.NodeProcessor{}},
            "SERVICE":   {Config: conf, RequestProcessor: &httpprocessor.ServiceProcessor{}},
            "CONFIGMAP": {Config: conf, RequestProcessor: &httpprocessor.ConfigProcessor{}},
            "RELEASE":   {Config: conf, RequestProcessor: &httpprocessor.ReleaseProcessor{}},
            "CLONESET":  {Config: conf, RequestProcessor: &httpprocessor.CloneSetProcessor{}},
            "RELEASEV2": {Config: conf, RequestProcessor: &httpprocessor.ReleaseV2Processor{}},
            "FILELOG":   {Config: conf, RequestProcessor: &httpprocessor.FileLogProcessor{}},
            "NAMESPACE": {Config: conf, RequestProcessor: &httpprocessor.NamespaceProcessor{}},
            "EVENT":     {Config: conf, RequestProcessor: &httpprocessor.EventProcessor{}},
            "VALIDATOR": {Config: conf, RequestProcessor: &httpprocessor.ValidatorProcessor{}},
            "SECRET":    {Config: conf, RequestProcessor: &httpprocessor.SecretProcessor{}},
        },
    }
    s.HandleFunc("/{type:namespace|deploy|service|pod|node|configmap|cloneset|secret}/{action:list|create|update|update-scale|get-scale|get|logs}", apiDispatcher.Handle).Methods(http.MethodGet, http.MethodPost)
    s.HandleFunc("/{type:release|releasev2}/{action:deploy|status}", apiDispatcher.Handle).Methods(http.MethodGet, http.MethodPost)
    s.HandleFunc("/common/{type:filelog|validator}", apiDispatcher.Handle).Methods(http.MethodGet, http.MethodPost)
    r.Handle("/ws/logs", httpprocessor.LogHandler{Config: conf}).Methods(http.MethodConnect, http.MethodGet, http.MethodPost)
    r.Handle("/ws/webssh", httpprocessor.WebSshHandler{Config: conf}).Methods(http.MethodConnect, http.MethodGet, http.MethodPost)
    r.Handle("/ws/event/watch/{type:event|deployment}", httpprocessor.EventWatcherHandler{Config: conf}).Methods(http.MethodConnect, http.MethodGet, http.MethodPost)
    r.NotFoundHandler = http.HandlerFunc(httpprocessor.NotfoundHandler)
    if *debug {
        err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
            pathTemplate, err := route.GetPathTemplate()
            if err == nil {
                fmt.Println("ROUTE:", pathTemplate)
            }
            pathRegexp, err := route.GetPathRegexp()
            if err == nil {
                fmt.Println("Path regexp:", pathRegexp)
            }
            queriesTemplates, err := route.GetQueriesTemplates()
            if err == nil {
                fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
            }
            queriesRegexps, err := route.GetQueriesRegexp()
            if err == nil {
                fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
            }
            methods, err := route.GetMethods()
            if err == nil {
                fmt.Println("Methods:", strings.Join(methods, ","))
            }
            fmt.Println()
            return nil
        })
        if err != nil {
            log.Error(err)
        }
    }
    log.Info("server is running...")
    log.Fatalln(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), r))
}
