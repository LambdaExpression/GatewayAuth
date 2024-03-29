package util

import (
	"GatewayAuth/src/config"
	"net/http"
	"strings"
)

func GetURL(r *http.Request) (Url string) {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	return strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
}

func GetUrlArg(r *http.Request, name string) string {
	var arg string
	values := r.URL.Query()
	arg = values.Get(name)
	return arg
}

func UseSsl(conf config.Config) bool {
	return conf.Base.SslCertificate != "" && conf.Base.SslCertificateKey != ""
}
