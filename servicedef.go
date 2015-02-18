package morganfield

import (
	"net/http"
	"regexp"
)

// Method: GET | PUT | POST | DELETE
// Uri: Regexp for path
// Internal_Protocol: http | https
// Internal_Host: Address used to connect to morganfield
// External_Protocol: http | https
// External_Host: Address used to connect to remote REST service
// Caller: Regexp for checking source of requests
// SetCookies: Should the cookie data be copied back and forth
// In_Object: Input JSON object (see objects.go)
// Out_Object: Output JSON object (see objects.go)
type Service_Definition struct {
	Method            string
	URI               *regexp.Regexp
	Internal_Protocol string
	Internal_Host     string
	External_Protocol string
	External_Host     string
	Caller            *regexp.Regexp // checked against RemoteAddr, form host:port!
	SetCookies        bool
	In_Object         interface{}
	Out_Object        interface{}
}

// Service is considered to match request, when the following match:
// method, url path, target host, remote address
func (s Service_Definition) Matches_Request(r *http.Request) bool {
	host := r.Host
	// request.Host may be empty - in that case try request url host
	if host == "" {
		host = r.URL.Host
	}
	return (r.Method == s.Method) && s.URI.MatchString(r.URL.Path) && (host == s.Internal_Host) && s.Caller.MatchString(r.RemoteAddr)
}

func Create_Service_Definition(
	method string,
	uri string,
	intprotocol string,
	inthost string,
	extprotocol string,
	exthost string,
	caller string,
	setcookies bool,
	inobj interface{},
	outobj interface{}) (Service_Definition, error) {

	service := Service_Definition{}

	service.Method = method
	// Will panic if doesn't compile
	service.URI = regexp.MustCompile(uri)
	service.Internal_Protocol = intprotocol
	service.Internal_Host = inthost
	service.External_Protocol = extprotocol
	service.External_Host = exthost
	// Will panic if doesn't compile
	service.Caller = regexp.MustCompile(caller)
	service.SetCookies = setcookies
	service.In_Object = inobj
	service.Out_Object = outobj

	return service, nil
}
