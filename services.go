package morganfield

import (
	"github.com/blackjack/syslog"
)

func Setup_services() {
	// Service definitions

	/*

	   Example services

	   c, err := Create_Service_Definition(
	       "POST",
	       "/rest/auth/1/session",
	       "http",
	       "localhost.localdomain",
	       "http",
	       "dev.localdomain",
	       ".*",
	       true,
	       Auth_1_Session_Post_Request{},
	       Auth_1_Session_Post_Response{})
	   if err != nil {
	       syslog.Errf("%v", err)
	       panic(err)
	   }

	   i, err := Create_Service_Definition(
	       "GET",
	       "/rest/auth/1/session",
	       "http",
	       "localhost.localdomain",
	       "http",
	       "dev.localdomain",
	       ".*",
	       true,
	       nil,
	       Auth_1_Session_Get_Response{})
	   if err != nil {
	       syslog.Errf("%v", err)
	       panic(err)
	   }
	   Services = append(Services, c)
	   Services = append(Services, i)

	*/
}
