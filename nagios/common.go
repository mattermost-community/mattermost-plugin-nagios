package nagios

import "strings"

func buildOptions(options []string) string {
	var b strings.Builder

	for i, o := range options {
		if i > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(o)
	}

	return b.String()
}

type FormatOptions struct {
	Whitespace bool
	Enumerate  bool
	Bitmask    bool
	Duration   bool
}

func (f FormatOptions) String() string {
	var options []string

	if f.Whitespace {
		options = append(options, "whitespace")
	}
	if f.Enumerate {
		options = append(options, "enumerate")
	}
	if f.Bitmask {
		options = append(options, "bitmask")
	}
	if f.Duration {
		options = append(options, "duration")
	}

	return buildOptions(options)
}

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
