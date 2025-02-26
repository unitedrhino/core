package proxy

import (
	"fmt"
	"gitee.com/unitedrhino/core/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	dir := http.Dir(utils.GerRealPwd(svcCtx.Config.Proxy.FileProxy.FrontDir))
	fileServer := http.FileServer(dir)
	return func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		upath := r.URL.Path
		for _, v := range svcCtx.Config.Proxy.StaticProxy {
			if strings.HasPrefix(upath, v.Router) {
				staticProxy(svcCtx, v, w, r)
				return
			}
		}
		header.Set("Access-Control-Allow-Headers", ctxs.HttpAllowHeader)
		header.Set("Access-Control-Allow-Origin", "*")
		f, err := dir.Open(upath)
		if err != nil {
			defaultHandle(svcCtx, upath, w, r, false)
			return
		} else {
			info, err := f.Stat()
			if err != nil || info.Mode().IsDir() {
				defaultHandle(svcCtx, upath, w, r, info.Mode().IsDir())
				return
			}
		}
		fileServer.ServeHTTP(w, r)
	}
}

func staticProxy(svcCtx *svc.ServiceContext, conf *conf.StaticProxyConf, w http.ResponseWriter, r *http.Request) {
	remote, err := url.Parse(conf.Dest)
	if err != nil {
		//defaultHandle(svcCtx, w, r)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	r.Host = remote.Host
	if conf.DeletePrefix {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, conf.Router)
	}
	//for k := range w.Header() {
	//	w.Header().Del(k)
	//}
	w.Header().Del("traceparent")
	proxy.ServeHTTP(w, r)
}
func defaultHandle(svcCtx *svc.ServiceContext, upath string, w http.ResponseWriter, r *http.Request, isDir bool) {
	if !svcCtx.Config.Proxy.FileProxy.IsEnable {
		http.NotFound(w, r)
		return
	}
	dir := http.Dir(utils.GerRealPwd(svcCtx.Config.Proxy.FileProxy.FrontDir))
	var (
		f   http.File
		err error
	)
	switch upath {
	case "/favicon.ico":
		f, err = dir.Open(fmt.Sprintf("%sfavicon.ico", svcCtx.Config.Proxy.FileProxy.CoreDir))
	case "/":
		http.Redirect(w, r, svcCtx.Config.Proxy.FileProxy.CoreDir, http.StatusMovedPermanently)
		return
	default:
		if strings.HasSuffix(upath, "/") {
			f, err = dir.Open(upath + "index.html")
		} else if isDir {
			f, err = dir.Open(upath + "/index.html")
		}
	}
	if err != nil || f == nil {
		http.NotFound(w, r)
		return
	}
	indexFile, err := io.ReadAll(f)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(indexFile)
	return
}
