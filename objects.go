package morganfield

/*
Legal types for JSON objects:
    bool, for JSON booleans
    float64, for JSON numbers
    string, for JSON strings
    []interface{}, for JSON arrays
    map[string]interface{}, for JSON objects
    nil for JSON null
    int, for JSON numbers

The "string" option signals that a field is stored as JSON inside a
JSON-encoded string. It applies only to fields of string, floating point,
or integer types.

Int64String int64 `json:",string"`

*/

// Example JSON objects as per https://docs.atlassian.com/jira/REST/latest/#d2e1869

type Auth_1_Session_Post_Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth_1_Session_Post_Response struct {
	Session struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"session"`
	LoginInfo struct {
		FailedLoginCount    int    `json:"failedLoginCount"`
		LoginCount          int    `json:"loginCount"`
		LastFailedLoginTime string `json:"lastFailedLoginTime"`
		PreviousLoginTime   string `json:"previousLoginTime"`
	} `json:"loginInfo"`
}

type Auth_1_Session_Get_Response struct {
	Self      string `json:"self"`
	Name      string `json:"name"`
	LoginInfo struct {
		FailedLoginCount    int    `json:"failedLoginCount"`
		LoginCount          int    `json:"loginCount"`
		LastFailedLoginTime string `json:"lastFailedLoginTime"`
		PreviousLoginTime   string `json:"previousLoginTime"`
	} `json:"loginInfo"`
}

func (a *Auth_1_Session_Post_Request) Validate() {
	// Extra validation!
}
