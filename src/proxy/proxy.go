package proxy

import (
	"GatewayAuth/src/config"
	"GatewayAuth/src/login"
	"GatewayAuth/src/util"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

// NewProxy takes target host and creates a reverse proxy
// NewProxy 拿到 targetHost 后，创建一个反向代理
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyRequestHandler handles the http request using proxy
// ProxyRequestHandler 使用 proxy 处理请求
func ProxyRequestHandler(conf config.Config, proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		loginState, cacheMaxAge, err := login.Login(conf, r)
		if err != nil && err.Error() != "http: named cookie not present" {
			log.Println(err)
		}
		switch loginState {
		case login.NotLogin:
			login.ClearCookie(w)
			if r.Method != "GET" {
				w.Header().Set("Content-Type", "application-json")
				io.Copy(w, strings.NewReader(string(`{"code":401,"msg":"Login session expired / 登录已失效"}`)))
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				http.SetCookie(w, &http.Cookie{Name: login.CookieKey, Path: "/", Value: "", HttpOnly: true, MaxAge: -1})
				w.Header().Set("Location", "/login?"+paramEncode(r))
				w.WriteHeader(http.StatusFound)
			}
		case login.NoPermission:
			if r.Method != "GET" {
				w.Header().Set("Content-Type", "application-json")
				io.Copy(w, strings.NewReader(string(`{"code":403,"msg":"Account has no permission / 账号无权限"}`)))
				w.WriteHeader(http.StatusForbidden)
			} else {
				w.Header().Set("Location", "/login?type=nopermission&"+paramEncode(r))
				w.WriteHeader(http.StatusFound)
			}
		case login.AlreadyLogin:
			serveHttp(proxy, cacheMaxAge, w, r)
		case login.NoLogin:
			serveHttp(proxy, cacheMaxAge, w, r)
		}
	}
}

func serveHttp(proxy *httputil.ReverseProxy, cacheMaxAge int64, w http.ResponseWriter, r *http.Request) {
	if cacheMaxAge > 0 {
		w.Header().Set("Cache-Control", "max-age="+strconv.FormatInt(cacheMaxAge, 10))
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}
	proxy.ServeHTTP(w, r)
}

func paramEncode(r *http.Request) string {

	u := util.GetURL(r)
	var param = url.Values{}
	param.Add("url", u)

	return param.Encode()
}
