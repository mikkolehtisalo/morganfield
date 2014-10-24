package morganfield

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/blackjack/syslog"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var (
	Services []Service_Definition = []Service_Definition{}
)

// Get the service by comparing request with service definitions
func get_service(r *http.Request) (Service_Definition, error) {
	x := Service_Definition{}
	found := false

	for _, s := range Services {
		if s.Matches_Request(r) {
			x = s
			found = true
		}
	}

	if !found {
		return x, fmt.Errorf("get_service: no service found for request")
	}

	return x, nil
}

// Initializes http client for forwarding the request
func get_http_client(s Service_Definition) *http.Client {
	jaropts := cookiejar.Options{
		PublicSuffixList: get_suffix_list(s),
	}

	jar, err := cookiejar.New(&jaropts)
	if err != nil {
		panic(err)
	}

	var client *http.Client
	if s.External_Protocol == "https" {
		client = &http.Client{
			Timeout:   60 * time.Second,
			Jar:       jar,
			Transport: get_transport(s.External_Host_Without_Port()),
		}
	} else {
		client = &http.Client{
			Timeout: 60 * time.Second,
			Jar:     jar,
		}
	}

	return client
}

// Get transport for https client
func get_transport(server_name string) *http.Transport {

	// Load client key and certificate (optional)
	cert, clerr := tls.LoadX509KeyPair("/opt/morganfield/etc/morganfield.crt", "/opt/morganfield/etc/morganfield.key")

	// Load the CA certificate for server certificate validation (mandatory)
	capool := x509.NewCertPool()
	cacert, err := ioutil.ReadFile("/opt/morganfield/etc/extca.crt")
	if err != nil {
		// proceeding without CA certificate checking would be bad
		syslog.Errf("%v", err)
		panic(err)
	}
	capool.AppendCertsFromPEM(cacert)

	// Prepare config and transport
	var config tls.Config
	if clerr != nil {
		// We did NOT have client certificate
		config = tls.Config{RootCAs: capool}
	} else {
		// We DID have client certificate so offer it
		config = tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      capool,
			ServerName:   server_name}
	}

	tr := &http.Transport{
		TLSClientConfig: &config,
	}

	return tr
}

// Handler for all requests
func Handler(w http.ResponseWriter, r *http.Request) {
	// Syslog connection should be ok for every request
	syslog.Openlog("morganfield", syslog.LOG_PID, syslog.LOG_DAEMON)
	// Initialize object for logging request
	reqlog := RequestLog_from_request(r)
	defer reqlog.Log()

	// Detect service
	s, err := get_service(r)
	if err != nil {
		reqlog.Status = http.StatusBadGateway
		w.WriteHeader(http.StatusBadGateway)
		panic(err)
	}

	// filter input JSON
	injson, err := filter_input_json(r, s)
	if err != nil {
		// The JSON will be empty string - no worries just log it
		syslog.Errf("%v", err)
		reqlog.InJson = fmt.Sprintf("%v", err)
	} else {
		reqlog.InJson = injson
	}

	// Build a new client & request
	client := get_http_client(s)

	request, err := http.NewRequest(s.Method, s.External_Protocol+"://"+s.External_Host+r.URL.Path, bytes.NewBufferString(injson))
	if err != nil {
		reqlog.Status = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")

	// Add cookies
	if s.SetCookies {
		origc := r.Cookies()
		for _, c := range origc {
			// When outgoing, cookie domain must match External_Host
			c.Domain = s.External_Host_Without_Port()
			reqlog.InCookies = append(reqlog.InCookies, fmt.Sprintf("%v", c))
			request.AddCookie(c)
		}
	}

	// Copy the user agent string
	if ua, ok := r.Header["User-Agent"]; ok {
		request.Header.Add("User-Agent", ua[0])
	}

	// Forward the request to target
	response, err := client.Do(request)
	if err != nil {
		reqlog.Status = http.StatusNotImplemented
		w.WriteHeader(http.StatusNotImplemented)
		panic(err)
	}

	// filter output JSON
	outjson, err := filter_output_json(response, s)
	if err != nil {
		// The JSON will be empty string - no worries just log it
		syslog.Errf("%v", err)
		reqlog.OutJson = fmt.Sprintf("%v", err)
	} else {
		reqlog.OutJson = outjson
	}

	// Cookies back
	if s.SetCookies {
		newc := response.Cookies()
		for _, c := range newc {
			// When incoming, cookie domain must match Internal_Host
			if len(c.Domain) > 0 {
				c.Domain = s.Internal_Host_Without_Port()
			}
			w.Header().Add("Set-Cookie", c.String())
			reqlog.OutCookies = append(reqlog.OutCookies, fmt.Sprintf("%v", c))
		}
	}

	// Status number
	w.WriteHeader(response.StatusCode)
	reqlog.Status = response.StatusCode

	// Return the JSON
	fmt.Fprintf(w, "%s", outjson)
}
