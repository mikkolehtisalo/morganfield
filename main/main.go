package main

import (
	"github.com/blackjack/syslog"
	"github.com/mikkolehtisalo/morganfield"
	"net/http"
	"time"
)

func main() {
	// Set gomaxprocs automatically
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Open syslog
	syslog.Openlog("morganfield", syslog.LOG_PID, syslog.LOG_DAEMON)

	server := &http.Server{
		Addr:           ":80", // port
		Handler:        http.HandlerFunc(morganfield.Handler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	morganfield.Setup_services()

	syslog.Errf("%v", server.ListenAndServe())
	//syslog.Errf("%v", server.ListenAndServeTLS("/opt/morganfield/etc/morganfield.crt", "/opt/morganfield/etc/morganfield.key"))
	syslog.Emerg("morganfield stopped running!")
}
