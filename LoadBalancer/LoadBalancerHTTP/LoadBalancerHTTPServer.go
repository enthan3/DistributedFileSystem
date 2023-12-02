package LoadBalancerHTTP

import (
	"DistributedFileSystem/LoadBalancer/LoadBalancerDefinition"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// RedirectToFrontendService Setting this server as a Reverse Proxy then redirect the request to Frontend Service.
func RedirectToFrontendService(w http.ResponseWriter, r *http.Request, l *LoadBalancerDefinition.LoadBalancerServer) {
	FrontendServiceURL, err := url.Parse("http://" + l.ServiceHTTP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(FrontendServiceURL)
	proxy.ServeHTTP(w, r)
}
