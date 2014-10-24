package morganfield

import (
	"strings"
)

type suffixlist struct {
	Internal_Host string
	External_Host string
}

// Fqdns only
// dev.localdomain OK
// .localdomain NOT ok
func (s suffixlist) PublicSuffix(domain string) string {
	// By default everything is public -> do not allow setting cookie
	result := domain

	// If it's the internal host, allow setting the cookie for this specific domain
	if domain == s.Internal_Host {
		result = strings.Join(strings.Split(s.Internal_Host, ".")[1:], ".")
	}

	// If it's the external host, allow setting the cookie for this specific domain
	if domain == s.External_Host {
		result = strings.Join(strings.Split(s.External_Host, ".")[1:], ".")
	}

	return result
}

// Informational method
func (s suffixlist) String() string {
	return "morganfield"
}

func get_suffix_list(s Service_Definition) suffixlist {
	// If there's a port, strip it off
	intt_host := s.Internal_Host_Without_Port()
	extt_host := s.External_Host_Without_Port()
	return suffixlist{
		Internal_Host: intt_host,
		External_Host: extt_host,
	}
}
