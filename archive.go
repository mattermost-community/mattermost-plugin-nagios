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

	q.SetNonEmpty("formatoptions", a.FormatOptions.String())

	if includeStartCount {
		q.SetNonEmpty("start", strconv.Itoa(a.Start))
		q.SetNonEmpty("count", strconv.Itoa(a.Count))
	}

	q.SetNonEmpty("dateformat", a.DateFormat)
	q.SetNonEmpty("objecttypes", a.ObjectTypes.String())
	q.SetNonEmpty("statetypes", a.StateTypes.String())
	q.SetNonEmpty("hoststates", a.HostStates.String())
	q.SetNonEmpty("servicestates", a.ServiceStates.String())
	q.SetNonEmpty("parenthost", a.ParentHost)
	q.SetNonEmpty("childhost", a.ChildHost)
	q.SetNonEmpty("hostname", a.HostName)
	q.SetNonEmpty("hostgroup", a.HostGroup)
	q.SetNonEmpty("servicegroup", a.ServiceGroup)
	q.SetNonEmpty("servicedescription", a.ServiceDescription)
	q.SetNonEmpty("contactname", a.ContactName)
	q.SetNonEmpty("contactgroup", a.ContactGroup)
	q.SetNonEmpty("backtrackedarchives", a.BacktrackedArchives)
	q.SetNonEmpty("starttime", strconv.FormatInt(a.StartTime, 10))
	q.SetNonEmpty("endtime", strconv.FormatInt(a.EndTime, 10))

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
	ObjectType   string `json:"object_type"`
	HostName     string `json:"host_name"`
	Description  string `json:"description"`
	StateType    string `json:"state_type"`
	State        string `json:"state"`
	PluginOutput string `json:"plugin_output"`
}

type AlertListData struct {
	Selectors map[string]json.RawMessage `json:"selectors"`
	AlertList []AlertListEntry           `json:"alertlist"`
}

type AlertList struct {
	FormatVersion int           `json:"format_version"`
	Result        Result        `json:"result"`
	Data          AlertListData `json:"data"`
}
