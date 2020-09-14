package nagios

type Result struct {
	QueryTime      int64  `json:"query_time"`
	CGI            string `json:"cgi"`
	User           string `json:"user"`
	Query          string `json:"query"`
	QueryStatus    string `json:"query_status"`
	ProgramStart   int64  `json:"program_start"`
	LastDataUpdate int64  `json:"last_data_update"`
	TypeCode       int    `json:"type_code"`
	TypeText       string `json:"type_text"`
	Message        string `json:"message"`
}
