package nagios

import (
	"strconv"
	"strings"
)

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
	q := Query{map[string][]string{}}

	q.Add("formatoptions", a.FormatOptions.String())

	if includeStartCount {
		q.Add("start", strconv.Itoa(a.Start))
		q.Add("count", strconv.Itoa(a.Count))
	}

	q.Add("dateformat", strconv.Itoa(a.Count))
	q.Add("objecttypes", a.ObjectTypes.String())
	q.Add("statetypes", a.StateTypes.String())
	q.Add("hoststates", a.HostStates.String())
	q.Add("servicestates", a.ServiceStates.String())
	q.Add("parenthost", a.ParentHost)
	q.Add("childhost", a.ChildHost)
	q.Add("hostname", a.HostName)
	q.Add("hostgroup", a.HostGroup)
	q.Add("servicegroup", a.ServiceGroup)
	q.Add("servicedescription", a.ServiceDescription)
	q.Add("contactname", a.ContactName)
	q.Add("contactgroup", a.ContactGroup)
	q.Add("backtrackedarchives", a.BacktrackedArchives)
	q.Add("starttime", strconv.FormatInt(a.StartTime, 10))
	q.Add("endtime", strconv.FormatInt(a.EndTime, 10))

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
