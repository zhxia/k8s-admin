用于操作k8s集群

项目初始化：
`go mod vendor`

编译打包：
`make`

服务启动:
` ./jobrunner --daemon`

配置文件:

    server:
        host: 0.0.0.0
        port: 9981
    k8s:
        namespace: default
        config: /Users/zhxia/workspace/go/k8s/config
    log:
        level: TRACE
        dir: /tmp
    redis:
        host: 127.0.0.1
        port: 6379
        password:

接口说明
    
 - 查询接口   
   1. 查询ingress、service、pod、configmap列表
   2.部署deployment
   3. 对已部署的deployment进行扩缩容调整
   http://localhost:9981/api/pod/list/\*/\*
   
 - 应用发布接口
    1. 增加基于Deployment的发布支持
    2. 增加基于OpenKurise的发布支持
 - 容器日志实时展示接口
 - 支持应用日志查看
 - 支持webssh方式进入POD 