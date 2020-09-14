package nagios

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

const archiveEndpoint = "archivejson.cgi"

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

type ObjectTypes struct {
	Host    bool
	Service bool
}

func (o ObjectTypes) String() string {
	var options []string

	if o.Host {
		options = append(options, "host")
	}
	if o.Service {
		options = append(options, "service")
	}

	return buildOptions(options)
}

type StateTypes struct {
	Soft bool
	Hard bool
}

func (s StateTypes) String() string {
	var options []string

	if s.Soft {
		options = append(options, "soft")
	}
	if s.Hard {
		options = append(options, "hard")
	}

	return buildOptions(options)
}

type HostStates struct {
	Up          bool
	Down        bool
	Unreachable bool
}

func (h HostStates) String() string {
	var options []string

	if h.Up {
		options = append(options, "up")
	}
	if h.Down {
		options = append(options, "down")
	}
	if h.Unreachable {
		options = append(options, "unreachable")
	}

	return buildOptions(options)
}

type ServiceStates struct {
	Ok       bool
	Warning  bool
	Critical bool
	Unknown  bool
}

func (s ServiceStates) String() string {
	var options []string

	if s.Ok {
		options = append(options, "ok")
	}
	if s.Warning {
		options = append(options, "warning")
	}
	if s.Critical {
		options = append(options, "critical")
	}
	if s.Unknown {
		options = append(options, "unknown")
	}

	return buildOptions(options)
}

type alertRequest struct {
	FormatOptions       FormatOptions
	Start               int
	Count               int
	DateFormat          string
	ObjectTypes         ObjectTypes
	StateTypes          StateTypes
	HostStates          HostStates
	ServiceStates       ServiceStates
	ParentHost          string
	ChildHost           string
	HostName            string
	HostGroup           string
	ServiceGroup        string
	ServiceDescription  string
	ContactName         string
	ContactGroup        string
	BacktrackedArchives string
	StartTime           int64
	EndTime             int64
}

func (a alertRequest) build(includeStartCount bool) Query {
	q := Query{
		Endpoint: archiveEndpoint,
		URLQuery: make(url.Values),
	}

	if v := a.FormatOptions.String(); len(v) > 0 {
		q.URLQuery.Add("formatoptions", v)
	}

	if includeStartCount {
		if a.Start > 0 {
			q.URLQuery.Add("start", strconv.Itoa(a.Start))
		}
		if a.Count > 0 {
			q.URLQuery.Add("count", strconv.Itoa(a.Count))
		}
	}

	if len(a.DateFormat) > 0 {
		q.URLQuery.Add("dateformat", a.DateFormat)
	}

	if v := a.ObjectTypes.String(); len(v) > 0 {
		q.URLQuery.Add("objecttypes", v)
	}
	if v := a.StateTypes.String(); len(v) > 0 {
		q.URLQuery.Add("statetypes", v)
	}
	if v := a.HostStates.String(); len(v) > 0 {
		q.URLQuery.Add("hoststates", v)
	}
	if v := a.ServiceStates.String(); len(v) > 0 {
		q.URLQuery.Add("servicestates", v)
	}

	if len(a.ParentHost) > 0 {
		q.URLQuery.Add("parenthost", a.ParentHost)
	}
	if len(a.ChildHost) > 0 {
		q.URLQuery.Add("childhost", a.ChildHost)
	}
	if len(a.HostName) > 0 {
		q.URLQuery.Add("hostname", a.HostName)
	}
	if len(a.HostGroup) > 0 {
		q.URLQuery.Add("hostgroup", a.HostGroup)
	}
	if len(a.ServiceGroup) > 0 {
		q.URLQuery.Add("servicegroup", a.ServiceGroup)
	}
	if len(a.ServiceDescription) > 0 {
		q.URLQuery.Add("servicedescription", a.ServiceDescription)
	}
	if len(a.ContactName) > 0 {
		q.URLQuery.Add("contactname", a.ContactName)
	}
	if len(a.ContactGroup) > 0 {
		q.URLQuery.Add("contactgroup", a.ContactGroup)
	}
	if len(a.BacktrackedArchives) > 0 {
		q.URLQuery.Add("backtrackedarchives", a.BacktrackedArchives)
	}

	q.URLQuery.Add("starttime", strconv.FormatInt(a.StartTime, 10))
	q.URLQuery.Add("endtime", strconv.FormatInt(a.EndTime, 10))

	return q
}

type AlertCountRequest struct {
	alertRequest
}

func (a AlertCountRequest) Build() Query {
	return a.build(false)
}

type AlertListRequest struct {
	alertRequest
}

func (a AlertListRequest) Build() Query {
	return a.build(true)
}

type AlertCountData struct {
	Selectors map[string]json.RawMessage `json:"selectors"`
	Count     int                        `json:"count"`
}

type AlertCount struct {
	FormatVersion int            `json:"format_version"`
	Result        Result         `json:"result"`
	Data          AlertCountData `json:"data"`
}

type AlertListEntry struct {
	Timestamp    int64  `json:"timestamp"`
	ObjectType   int    `json:"object_type"`
	HostName     string `json:"host_name"`
	Description  string `json:"description"`
	StateType    int    `json:"state_type"`
	State        int    `json:"state"`
	PluginOutput string `json:"plugin_output"`
}

type AlertListData struct {
	Selectors map[string]json.RawMessage `json:"selectors"`
	Entries   []AlertListEntry           `json:"alertlist"`
}

type AlertList struct {
	FormatVersion int           `json:"format_version"`
	Result        Result        `json:"result"`
	Data          AlertListData `json:"data"`
}
