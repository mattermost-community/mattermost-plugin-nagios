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

type HostNotificationTypes struct {
	NoData        bool
	Down          bool
	Unreachable   bool
	Recovery      bool
	HostCustom    bool
	HostAck       bool
	HostFlapStart bool
	HostFlapStop  bool
}

func (h HostNotificationTypes) String() string {
	var options []string

	if h.NoData {
		options = append(options, "nodata")
	}
	if h.Down {
		options = append(options, "down")
	}
	if h.Unreachable {
		options = append(options, "unreachable")
	}
	if h.Recovery {
		options = append(options, "recovery")
	}
	if h.HostCustom {
		options = append(options, "hostcustom")
	}
	if h.HostAck {
		options = append(options, "hostack")
	}
	if h.HostFlapStart {
		options = append(options, "hostflapstart")
	}
	if h.HostFlapStop {
		options = append(options, "hostflapstop")
	}

	return buildOptions(options)
}

type ServiceNotificationTypes struct {
	NoData           bool
	Critical         bool
	Warning          bool
	Recovery         bool
	Custom           bool
	ServiceAck       bool
	ServiceFlapStart bool
	ServiceFlapStop  bool
	Unknown          bool
}

func (s ServiceNotificationTypes) String() string {
	var options []string

	if s.NoData {
		options = append(options, "nodata")
	}
	if s.Critical {
		options = append(options, "critical")
	}
	if s.Warning {
		options = append(options, "warning")
	}
	if s.Recovery {
		options = append(options, "recovery")
	}
	if s.Custom {
		options = append(options, "custom")
	}
	if s.ServiceAck {
		options = append(options, "serviceack")
	}
	if s.ServiceFlapStart {
		options = append(options, "serviceflapstart")
	}
	if s.ServiceFlapStop {
		options = append(options, "serviceflapstop")
	}
	if s.Unknown {
		options = append(options, "unknown")
	}

	return buildOptions(options)
}

type notificationRequest struct {
	FormatOptions            FormatOptions
	Start                    int
	Count                    int
	DateFormat               string
	ObjectTypes              ObjectTypes
	HostNotificationTypes    HostNotificationTypes
	ServiceNotificationTypes ServiceNotificationTypes
	ParentHost               string
	ChildHost                string
	HostName                 string
	HostGroup                string
	ServiceGroup             string
	ServiceDescription       string
	ContactName              string
	ContactGroup             string
	NotificationMethod       string
	BacktrackedArchives      string
	StartTime                int64
	EndTime                  int64
}

func (n notificationRequest) build(includeStartCount bool) Query {
	q := Query{
		Endpoint: archiveEndpoint,
		URLQuery: make(url.Values),
	}

	q.SetNonEmpty("formatoptions", n.FormatOptions.String())

	if includeStartCount {
		q.SetNonEmpty("start", strconv.Itoa(n.Start))
		q.SetNonEmpty("count", strconv.Itoa(n.Count))
	}

	q.SetNonEmpty("dateformat", n.DateFormat)
	q.SetNonEmpty("objecttypes", n.ObjectTypes.String())
	q.SetNonEmpty("hostnotificationtypes", n.HostNotificationTypes.String())
	q.SetNonEmpty("servicenotificationtypes", n.ServiceNotificationTypes.String())
	q.SetNonEmpty("parenthost", n.ParentHost)
	q.SetNonEmpty("childhost", n.ChildHost)
	q.SetNonEmpty("hostname", n.HostName)
	q.SetNonEmpty("hostgroup", n.HostGroup)
	q.SetNonEmpty("servicegroup", n.ServiceGroup)
	q.SetNonEmpty("servicedescription", n.ServiceDescription)
	q.SetNonEmpty("contactname", n.ContactName)
	q.SetNonEmpty("contactgroup", n.ContactGroup)
	q.SetNonEmpty("notificationmethod", n.NotificationMethod)
	q.SetNonEmpty("backtrackedarchives", n.BacktrackedArchives)
	q.SetNonEmpty("starttime", strconv.FormatInt(n.StartTime, 10))
	q.SetNonEmpty("endtime", strconv.FormatInt(n.EndTime, 10))

	return q
}

type NotificationCountRequest struct {
	notificationRequest
}

func (n NotificationCountRequest) Build() Query {
	return n.build(false)
}

type NotificationListRequest struct {
	notificationRequest
}

func (n NotificationListRequest) Build() Query {
	return n.build(true)
}
