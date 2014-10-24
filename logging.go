package morganfield

import (
	"github.com/blackjack/syslog"
	"net/http"
)

type RequestLog struct {
	Method        string
	Host          string
	Path          string
	Proto         string
	ContentLength int64
	ContentType   string `json:"Content-Type"`
	UserAgent     string `json:"User-Agent"`
	RemoteAddr    string
	InCookies     []string
	OutCookies    []string
	InJson        string
	OutJson       string
	Status        int
}

// Print object to syslog as JSON string
func (reqlog *RequestLog) Log() {
	r_json, _ := Marshal(reqlog)
	syslog.Infof("%v", r_json)
}

// Initialize new RequestLog and fill in some details from request
func RequestLog_from_request(r *http.Request) RequestLog {
	reqlog := RequestLog{}
	reqlog.Method = r.Method
	reqlog.Host = r.Host
	if r.Host == "" {
		reqlog.Host = r.URL.Host
	}
	reqlog.Path = r.URL.Path
	reqlog.Proto = r.Proto
	reqlog.ContentLength = r.ContentLength
	if ct, ok := r.Header["Content-Type"]; ok {
		reqlog.ContentType = ct[0]
	} else {
		reqlog.ContentType = ""
	}
	if ua, ok := r.Header["User-Agent"]; ok {
		reqlog.UserAgent = ua[0]
	} else {
		reqlog.UserAgent = ""
	}
	reqlog.RemoteAddr = r.RemoteAddr
	return reqlog
}
