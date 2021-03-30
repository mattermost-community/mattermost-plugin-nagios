package nagios

import (
	"encoding/json"
	"net/url"
	"strconv"
)

const archiveEndpoint = "archivejson.cgi"

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

type GeneralAlertRequest struct {
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

func (g GeneralAlertRequest) build(query string, includeStartCount bool) Query {
	q := Query{
		Endpoint: archiveEndpoint,
		URLQuery: make(url.Values),
	}

	q.SetNonEmpty("query", query)
	q.SetNonEmpty("formatoptions", g.FormatOptions.String())

	if includeStartCount {
		q.URLQuery.Set("start", strconv.Itoa(g.Start))
		if g.Count > 0 {
			q.URLQuery.Set("count", strconv.Itoa(g.Count))
		}
	}

	q.SetNonEmpty("dateformat", g.DateFormat)
	q.SetNonEmpty("objecttypes", g.ObjectTypes.String())
	q.SetNonEmpty("statetypes", g.StateTypes.String())
	q.SetNonEmpty("hoststates", g.HostStates.String())
	q.SetNonEmpty("servicestates", g.ServiceStates.String())
	q.SetNonEmpty("parenthost", g.ParentHost)
	q.SetNonEmpty("childhost", g.ChildHost)
	q.SetNonEmpty("hostname", g.HostName)
	q.SetNonEmpty("hostgroup", g.HostGroup)
	q.SetNonEmpty("servicegroup", g.ServiceGroup)
	q.SetNonEmpty("servicedescription", g.ServiceDescription)
	q.SetNonEmpty("contactname", g.ContactName)
	q.SetNonEmpty("contactgroup", g.ContactGroup)
	q.SetNonEmpty("backtrackedarchives", g.BacktrackedArchives)

	q.URLQuery.Set("starttime", strconv.FormatInt(g.StartTime, 10))
	q.URLQuery.Set("endtime", strconv.FormatInt(g.EndTime, 10))

	return q
}

type AlertCountRequest struct {
	GeneralAlertRequest
}

func (a AlertCountRequest) Build() Query {
	return a.build("alertcount", false)
}

type AlertListRequest struct {
	GeneralAlertRequest
}

func (a AlertListRequest) Build() Query {
	return a.build("alertlist", true)
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
	Name         string `json:"name"`
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

type GeneralNotificationRequest struct {
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

func (g GeneralNotificationRequest) build(query string, includeStartCount bool) Query {
	q := Query{
		Endpoint: archiveEndpoint,
		URLQuery: make(url.Values),
	}

	q.SetNonEmpty("query", query)
	q.SetNonEmpty("formatoptions", g.FormatOptions.String())

	if includeStartCount {
		q.URLQuery.Set("start", strconv.Itoa(g.Start))
		if g.Count > 0 {
			q.URLQuery.Set("count", strconv.Itoa(g.Count))
		}
	}

	q.SetNonEmpty("dateformat", g.DateFormat)
	q.SetNonEmpty("objecttypes", g.ObjectTypes.String())
	q.SetNonEmpty("hostnotificationtypes", g.HostNotificationTypes.String())
	q.SetNonEmpty("servicenotificationtypes", g.ServiceNotificationTypes.String())
	q.SetNonEmpty("parenthost", g.ParentHost)
	q.SetNonEmpty("childhost", g.ChildHost)
	q.SetNonEmpty("hostname", g.HostName)
	q.SetNonEmpty("hostgroup", g.HostGroup)
	q.SetNonEmpty("servicegroup", g.ServiceGroup)
	q.SetNonEmpty("servicedescription", g.ServiceDescription)
	q.SetNonEmpty("contactname", g.ContactName)
	q.SetNonEmpty("contactgroup", g.ContactGroup)
	q.SetNonEmpty("notificationmethod", g.NotificationMethod)
	q.SetNonEmpty("backtrackedarchives", g.BacktrackedArchives)

	q.URLQuery.Set("starttime", strconv.FormatInt(g.StartTime, 10))
	q.URLQuery.Set("endtime", strconv.FormatInt(g.EndTime, 10))

	return q
}

type NotificationCountRequest struct {
	GeneralNotificationRequest
}

func (n NotificationCountRequest) Build() Query {
	return n.build("notificationcount", false)
}

type NotificationListRequest struct {
	GeneralNotificationRequest
}

func (n NotificationListRequest) Build() Query {
	return n.build("notificationlist", true)
}

type NotificationCountData struct {
	Selectors map[string]json.RawMessage `json:"selectors"`
	Count     int                        `json:"count"`
}

type NotificationCount struct {
	FormatVersion int                   `json:"format_version"`
	Result        Result                `json:"result"`
	Data          NotificationCountData `json:"data"`
}

type NotificationListEntry struct {
	Timestamp        int64  `json:"timestamp"`
	ObjectType       string `json:"object_type"`
	HostName         string `json:"host_name"`
	Description      string `json:"description"`
	Name             string `json:"name"`
	Contact          string `json:"contact"`
	NotificationType string `json:"notification_type"`
	Method           string `json:"method"`
	Message          string `json:"message"`
}

type NotificationListData struct {
	Selectors        map[string]json.RawMessage `json:"selectors"`
	NotificationList []NotificationListEntry    `json:"notificationlist"`
}

type NotificationList struct {
	FormatVersion int                  `json:"format_version"`
	Result        Result               `json:"result"`
	Data          NotificationListData `json:"data"`
}
