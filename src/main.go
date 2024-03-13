package main

import (
	"GatewayAuth/src/config"
	"GatewayAuth/src/login"
	"GatewayAuth/src/proxy"
	"GatewayAuth/src/util"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var configPath string
var printVersion bool

var Version = "v1.0.3"
var GoVersion = "not set"
var GitCommit = "not set"
var BuildTime = "not set"

func main() {

	start()

	if printVersion {
		version()
		return
	}

	conf := config.Get(configPath)
	jbyte, _ := json.Marshal(conf)
	log.Println(string(jbyte))

	configProxy(conf)

	login.HttpLogin(conf)

	log.Println("listen : " + strconv.Itoa(conf.Base.Port))

	if util.UseSsl(conf) {
		log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(conf.Base.Port), conf.Base.SslCertificate, conf.Base.SslCertificateKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.Base.Port), nil))
	}
}

func start() {
	flag.StringVar(&configPath, "c", "./config", "--c config file path / 配置文件路径")
	flag.BoolVar(&printVersion, "v", false, "--v 打印程序构建版本")
	flag.Parse()
}

func version() {
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Go Version: %s\n", GoVersion)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Build Time: %s\n", BuildTime)
}

func configProxy(conf config.Config) {

	for _, p := range conf.Base.ProxySort {
		n := conf.Proxy[p]
		path := n.Path
		proxyFunc(conf, path, n.Target)
	}

}

func proxyFunc(conf config.Config, path, target string) {
	proxy2, err := proxy.NewProxy(target)
	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	// 使用 proxy 处理所有请求到你的服务
	http.HandleFunc(path, proxy.ProxyRequestHandler(conf, proxy2))
}
