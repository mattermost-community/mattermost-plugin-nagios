package nagios

import (
	"encoding/json"
	"net/url"
	"strconv"
)

const statusEndpoint = "statusjson.cgi"

type HostStatus struct {
	Up          bool
	Down        bool
	Unreachable bool
	Pending     bool
}

func (h HostStatus) String() string {
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
	if h.Pending {
		options = append(options, "pending")
	}

	return buildOptions(options)
}

type GeneralHostRequest struct {
	FormatOptions                  FormatOptions
	Start                          int
	Count                          int
	ParentHost                     string
	ChildHost                      string
	ShowDetails                    bool
	DateFormat                     string
	HostGroup                      string
	HostStatus                     HostStatus
	ContactGroup                   string
	CheckTimeperiodName            string
	HostNotificationTimeperiodName string
	CheckCommandName               string
	EventHandlerName               string
	ContactName                    string
	HostTimeField                  string
	StartTime                      int64
	EndTime                        int64
}

func (g GeneralHostRequest) build(query string, includeStartCount bool) Query {
	q := Query{
		Endpoint: statusEndpoint,
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

	q.SetNonEmpty("parenthost", g.ParentHost)
	q.SetNonEmpty("childhost", g.ChildHost)

	if g.ShowDetails {
		q.URLQuery.Set("details", strconv.FormatBool(g.ShowDetails))
	}

	q.SetNonEmpty("dateformat", g.DateFormat)
	q.SetNonEmpty("hostgroup", g.HostGroup)
	q.SetNonEmpty("hoststatus", g.HostStatus.String())
	q.SetNonEmpty("contactgroup", g.ContactGroup)
	q.SetNonEmpty("checktimeperiod", g.CheckTimeperiodName)
	q.SetNonEmpty("hostnotificationtimeperiod", g.HostNotificationTimeperiodName)
	q.SetNonEmpty("checkcommand", g.CheckCommandName)
	q.SetNonEmpty("eventhandler", g.EventHandlerName)
	q.SetNonEmpty("contactname", g.ContactName)

	q.SetNonEmpty("hosttimefield", g.HostTimeField)
	q.SetNonEmpty("starttime", strconv.FormatInt(g.StartTime, 10))
	q.SetNonEmpty("endtime", strconv.FormatInt(g.EndTime, 10))

	return q
}

type HostCountRequest struct {
	GeneralHostRequest
}

func (h HostCountRequest) Build() Query {
	return h.build("hostcount", false)
}

type HostListRequest struct {
	GeneralHostRequest
}

func (h HostListRequest) Build() Query {
	return h.build("hostlist", true)
}

type HostStatusCount struct {
	Up          int `json:"up"`
	Down        int `json:"down"`
	Unreachable int `json:"unreachable"`
	Pending     int `json:"pending"`
}

type HostCountData struct {
	Selectors map[string]json.RawMessage `json:"selectors"`
	Count     HostStatusCount            `json:"count"`
}

type HostCount struct {
	FormatVersion int           `json:"format_version"`
	Result        Result        `json:"result"`
	Data          HostCountData `json:"data"`
}

// TODO(DanielSz50): Handle showDetails check properly.
type HostListData struct {
	Selectors map[string]json.RawMessage `json:"selectors"`
	HostList  map[string]json.RawMessage `json:"hostlist"`
}

type HostList struct {
	FormatVersion int          `json:"format_version"`
	Result        Result       `json:"result"`
	Data          HostListData `json:"data"`
}

type ServiceStatus struct {
	Ok       bool
	Warning  bool
	Critical bool
	Unknown  bool
	Pending  bool
}

func (s ServiceStatus) String() string {
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
	if s.Pending {
		options = append(options, "pending")
	}

	return buildOptions(options)
}

type GeneralServiceRequest struct {
	FormatOptions                     FormatOptions
	Start                             int
	Count                             int
	ParentHost                        string
	ChildHost                         string
	ShowDetails                       bool
	DateFormat                        string
	HostName                          string
	HostGroup                         string
	HostStatus                        HostStatus
	ServiceGroup                      string
	ServiceStatus                     ServiceStatus
	ParentService                     string
	ChildService                      string
	ContactGroup                      string
	ServiceDescription                string
	CheckTimeperiodName               string
	ServiceNotificationTimeperiodName string
	CheckCommandName                  string
	EventHandlerName                  string
	ContactName                       string
	ServiceTimeField                  string
	StartTime                         int64
	EndTime                           int64
}

func (g GeneralServiceRequest) build(query string, includeStartCount bool) Query {
	q := Query{
		Endpoint: statusEndpoint,
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

	q.SetNonEmpty("parenthost", g.ParentHost)
	q.SetNonEmpty("childhost", g.ChildHost)

	if g.ShowDetails {
		q.URLQuery.Set("details", strconv.FormatBool(g.ShowDetails))
	}

	q.SetNonEmpty("dateformat", g.DateFormat)
	q.SetNonEmpty("hostname", g.HostName)
	q.SetNonEmpty("hostgroup", g.HostGroup)
	q.SetNonEmpty("hoststatus", g.HostStatus.String())
	q.SetNonEmpty("servicegroup", g.ServiceGroup)
	q.SetNonEmpty("servicestatus", g.ServiceStatus.String())
	q.SetNonEmpty("parentservice", g.ParentService)
	q.SetNonEmpty("childservice", g.ChildService)
	q.SetNonEmpty("contactgroup", g.ContactGroup)
	q.SetNonEmpty("servicedescription", g.ServiceDescription)
	q.SetNonEmpty("checktimeperiod", g.CheckTimeperiodName)
	q.SetNonEmpty("servicenotificationtimeperiod", g.ServiceNotificationTimeperiodName)
	q.SetNonEmpty("checkcommand", g.CheckCommandName)
	q.SetNonEmpty("eventhandler", g.EventHandlerName)
	q.SetNonEmpty("contactname", g.ContactName)

	q.SetNonEmpty("servicetimefield", g.ServiceTimeField)
	q.SetNonEmpty("starttime", strconv.FormatInt(g.StartTime, 10))
	q.SetNonEmpty("endtime", strconv.FormatInt(g.EndTime, 10))

	return q
}

type ServiceCountRequest struct {
	GeneralServiceRequest
}

func (s ServiceCountRequest) Build() Query {
	return s.build("servicecount", false)
}

type ServiceListRequest struct {
	GeneralServiceRequest
}

func (s ServiceListRequest) Build() Query {
	return s.build("servicelist", true)
}

type ServiceStatusCount struct {
	Ok       int `json:"ok"`
	Warning  int `json:"warning"`
	Critical int `json:"critical"`
	Unknown  int `json:"unknown"`
	Pending  int `json:"pending"`
}

type ServiceCountData struct {
	Selectors map[string]json.RawMessage `json:"selectors"`
	Count     ServiceStatusCount         `json:"count"`
}

type ServiceCount struct {
	FormatVersion int              `json:"format_version"`
	Result        Result           `json:"result"`
	Data          ServiceCountData `json:"data"`
}

type ServiceListData struct {
	Selectors   map[string]json.RawMessage `json:"selectors"`
	ServiceList map[string]json.RawMessage `json:"servicelist"`
}

type ServiceList struct {
	FormatVersion int             `json:"format_version"`
	Result        Result          `json:"result"`
	Data          ServiceListData `json:"data"`
}

type HostRequest struct {
	FormatOptions FormatOptions
	DateFormat    string
	HostName      string
}

func (h HostRequest) Build() Query {
	q := Query{
		Endpoint: statusEndpoint,
		URLQuery: make(url.Values),
	}

	q.SetNonEmpty("query", "host")
	q.SetNonEmpty("formatoptions", h.FormatOptions.String())
	q.SetNonEmpty("dateformat", h.DateFormat)

	q.URLQuery.Set("hostname", h.HostName)

	return q
}

type HostDetails struct {
	Name              string `json:"name"`
	PluginOutput      string `json:"plugin_output"`
	LongPluginOutput  string `json:"long_plugin_output"`
	PerfData          string `json:"perf_data"`
	Status            string `json:"status"`
	LastUpdate        int64  `json:"last_update"`
	HasBeenChecked    bool   `json:"has_been_checked"`
	ShouldBeScheduled bool   `json:"should_be_scheduled"`
	CurrentAttempt    int    `json:"current_attempt"`
	MaxAttempts       int    `json:"max_attempts"`
	LastCheck         int64  `json:"last_check"`
	NextCheck         int64  `json:"next_check"`

	// CheckOptions with bitmask turns into an array.
	CheckOptions               json.RawMessage `json:"check_options"`
	CheckType                  string          `json:"check_type"`
	LastStateChange            int64           `json:"last_state_change"`
	LastHardStateChange        int64           `json:"last_hard_state_change"`
	LastHardState              string          `json:"last_hard_state"`
	LastTimeUp                 int64           `json:"last_time_up"`
	LastTimeDown               int             `json:"last_time_down"`
	LastTimeUnreachable        int             `json:"last_time_unreachable"`
	StateType                  string          `json:"state_type"`
	LastNotification           int             `json:"last_notification"`
	NextNotification           int             `json:"next_notification"`
	NoMoreNotifications        bool            `json:"no_more_notifications"`
	NotificationsEnabled       bool            `json:"notifications_enabled"`
	ProblemHasBeenAcknowledged bool            `json:"problem_has_been_acknowledged"`
	AcknowledgementType        string          `json:"acknowledgement_type"`
	CurrentNotificationNumber  int             `json:"current_notification_number"`
	AcceptPassiveChecks        bool            `json:"accept_passive_checks"`
	EventHandlerEnabled        bool            `json:"event_handler_enabled"`
	ChecksEnabled              bool            `json:"checks_enabled"`
	FlapDetectionEnabled       bool            `json:"flap_detection_enabled"`
	IsFlapping                 bool            `json:"is_flapping"`
	PercentStateChange         float64         `json:"percent_state_change"`
	Latency                    float64         `json:"latency"`
	ExecutionTime              float64         `json:"execution_time"`
	ScheduledDowntimeDepth     int             `json:"scheduled_downtime_depth"`
	ProcessPerformanceData     bool            `json:"process_performance_data"`
	Obsess                     bool            `json:"obsess"`
}

type HostData struct {
	Details HostDetails `json:"host"`
}

type Host struct {
	FormatVersion int      `json:"format_version"`
	Result        Result   `json:"result"`
	Data          HostData `json:"data"`
}

type ServiceRequest struct {
	FormatOptions      FormatOptions
	DateFormat         string
	HostName           string
	ServiceDescription string
}

func (s ServiceRequest) Build() Query {
	q := Query{
		Endpoint: statusEndpoint,
		URLQuery: make(url.Values),
	}

	q.SetNonEmpty("query", "service")
	q.SetNonEmpty("formatoptions", s.FormatOptions.String())
	q.SetNonEmpty("dateformat", s.DateFormat)

	q.URLQuery.Set("hostname", s.HostName)
	q.URLQuery.Set("servicedescription", s.ServiceDescription)

	return q
}

type ServiceDetails struct {
	HostName          string `json:"host_name"`
	Description       string `json:"description"`
	PluginOutput      string `json:"plugin_output"`
	LongPluginOutput  string `json:"long_plugin_output"`
	PerfData          string `json:"perf_data"`
	MaxAttempts       int    `json:"max_attempts"`
	CurrentAttempt    int    `json:"current_attempt"`
	Status            string `json:"status"`
	LastUpdate        int64  `json:"last_update"`
	HasBeenChecked    bool   `json:"has_been_checked"`
	ShouldBeScheduled bool   `json:"should_be_scheduled"`
	LastCheck         int64  `json:"last_check"`

	// CheckOptions with bitmask turns to the array
	CheckOptions               json.RawMessage `json:"check_options"`
	CheckType                  string          `json:"check_type"`
	ChecksEnabled              bool            `json:"checks_enabled"`
	LastStateChange            int64           `json:"last_state_change"`
	LastHardStateChange        int64           `json:"last_hard_state_change"`
	LastHardState              string          `json:"last_hard_state"`
	LastTimeOk                 int64           `json:"last_time_ok"`
	LastTimeWarning            int64           `json:"last_time_warning"`
	LastTimeUnknown            int             `json:"last_time_unknown"`
	LastTimeCritical           int64           `json:"last_time_critical"`
	StateType                  string          `json:"state_type"`
	LastNotification           int64           `json:"last_notification"`
	NextNotification           int64           `json:"next_notification"`
	NextCheck                  int             `json:"next_check"`
	NoMoreNotifications        bool            `json:"no_more_notifications"`
	NotificationsEnabled       bool            `json:"notifications_enabled"`
	ProblemHasBeenAcknowledged bool            `json:"problem_has_been_acknowledged"`
	AcknowledgementType        string          `json:"acknowledgement_type"`
	CurrentNotificationNumber  int             `json:"current_notification_number"`
	AcceptPassiveChecks        bool            `json:"accept_passive_checks"`
	EventHandlerEnabled        bool            `json:"event_handler_enabled"`
	FlapDetectionEnabled       bool            `json:"flap_detection_enabled"`
	IsFlapping                 bool            `json:"is_flapping"`
	PercentStateChange         float64         `json:"percent_state_change"`
	Latency                    float64         `json:"latency"`
	ExecutionTime              float64         `json:"execution_time"`
	ScheduledDowntimeDepth     int             `json:"scheduled_downtime_depth"`
	ProcessPerformanceData     bool            `json:"process_performance_data"`
	Obsess                     bool            `json:"obsess"`
}

type ServiceData struct {
	Service ServiceDetails `json:"service"`
}

type Service struct {
	FormatVersion int         `json:"format_version"`
	Result        Result      `json:"result"`
	Data          ServiceData `json:"data"`
}

type PerformanceDataRequest struct {
	FormatOptions FormatOptions
	DateFormat    string
}

func (p PerformanceDataRequest) Build() Query {
	q := Query{
		Endpoint: statusEndpoint,
		URLQuery: make(url.Values),
	}

	q.SetNonEmpty("query", "performancedata")
	q.SetNonEmpty("formatoptions", p.FormatOptions.String())
	q.SetNonEmpty("dateformat", p.DateFormat)

	return q
}

type (
	Checks struct {
		OneMin     int `json:"1min"`
		FiveMin    int `json:"5min"`
		FifteenMin int `json:"15min"`
		OneHour    int `json:"1hour"`
		Start      int `json:"start"`
	}

	CheckExecutionTime struct {
		Min     float64 `json:"min"`
		Max     float64 `json:"max"`
		Average float64 `json:"average"`
	}

	CheckLatency struct {
		Min     float64 `json:"min"`
		Max     float64 `json:"max"`
		Average float64 `json:"average"`
	}

	PercentStateChange struct {
		Min     float64 `json:"min"`
		Max     float64 `json:"max"`
		Average float64 `json:"average"`
	}

	MetricsActive struct {
		CheckExecutionTime CheckExecutionTime `json:"check_execution_time"`
		CheckLatency       CheckLatency       `json:"check_latency"`
		PercentStateChange PercentStateChange `json:"percent_state_change"`
	}

	Active struct {
		Checks  Checks        `json:"checks"`
		Metrics MetricsActive `json:"metrics"`
	}

	MetricsPassive struct {
		PercentStateChange PercentStateChange `json:"percent_state_change"`
	}

	Passive struct {
		Checks  Checks         `json:"checks"`
		Metrics MetricsPassive `json:"metrics"`
	}

	ServiceChecks struct {
		Active  Active  `json:"active"`
		Passive Passive `json:"passive"`
	}

	HostChecks struct {
		Active  Active  `json:"active"`
		Passive Passive `json:"passive"`
	}

	ActiveScheduledHostChecks struct {
		OneMin     int `json:"1min"`
		FiveMin    int `json:"5min"`
		FifteenMin int `json:"15min"`
	}

	ActiveOnDemandHostChecks struct {
		OneMin     int `json:"1min"`
		FiveMin    int `json:"5min"`
		FifteenMin int `json:"15min"`
	}

	ParallelHostChecks struct {
		OneMin     int `json:"1min"`
		FiveMin    int `json:"5min"`
		FifteenMin int `json:"15min"`
	}

	SerialHostChecks struct {
		OneMin     int `json:"1min"`
		FiveMin    int `json:"5min"`
		FifteenMin int `json:"15min"`
	}

	CachedHostChecks struct {
		OneMin     int `json:"1min"`
		FiveMin    int `json:"5min"`
		FifteenMin int `json:"15min"`
	}

	PassiveHostChecks struct {
		OneMin     int `json:"1min"`
		FiveMin    int `json:"5min"`
		FifteenMin int `json:"15min"`
	}

	ActiveScheduledServiceChecks struct {
		OneMin     int `json:"1min"`
		FiveMin    int `json:"5min"`
		FifteenMin int `json:"15min"`
	}

	ActiveOnDemandServiceChecks struct {
		OneMin     int `json:"1min"`
		FiveMin    int `json:"5min"`
		FifteenMin int `json:"15min"`
	}

	CachedServiceChecks struct {
		OneMin     int `json:"1min"`
		FiveMin    int `json:"5min"`
		FifteenMin int `json:"15min"`
	}

	PassiveServiceChecks struct {
		OneMin     int `json:"1min"`
		FiveMin    int `json:"5min"`
		FifteenMin int `json:"15min"`
	}

	ExternalCommands struct {
		OneMin     int `json:"1min"`
		FiveMin    int `json:"5min"`
		FifteenMin int `json:"15min"`
	}

	CheckStatistics struct {
		ActiveScheduledHostChecks    ActiveScheduledHostChecks    `json:"active_scheduled_host_checks"`
		ActiveOnDemandHostChecks     ActiveOnDemandHostChecks     `json:"active_ondemand_host_checks"`
		ParallelHostChecks           ParallelHostChecks           `json:"parallel_host_checks"`
		SerialHostChecks             SerialHostChecks             `json:"serial_host_checks"`
		CachedHostChecks             CachedHostChecks             `json:"cached_host_checks"`
		PassiveHostChecks            PassiveHostChecks            `json:"passive_host_checks"`
		ActiveScheduledServiceChecks ActiveScheduledServiceChecks `json:"active_scheduled_service_checks"`
		ActiveOnDemandServiceChecks  ActiveOnDemandServiceChecks  `json:"active_ondemand_service_checks"`
		CachedServiceChecks          CachedServiceChecks          `json:"cached_service_checks"`
		PassiveServiceChecks         PassiveServiceChecks         `json:"passive_service_checks"`
		ExternalCommands             ExternalCommands             `json:"external_commands"`
	}

	ExternalCommandsBuffer struct {
		InUse          int `json:"in_use"`
		MaxUsed        int `json:"max_used"`
		TotalAvailable int `json:"total_available"`
	}

	BufferUsage struct {
		ExternalCommands ExternalCommandsBuffer `json:"external_commands"`
	}

	ProgramStatus struct {
		ServiceChecks   ServiceChecks   `json:"service_checks"`
		HostChecks      HostChecks      `json:"host_checks"`
		CheckStatistics CheckStatistics `json:"check_statistics"`
		BufferUsage     BufferUsage     `json:"buffer_usage"`
	}

	PerformanceData struct {
		ProgramStatus ProgramStatus `json:"programstatus"`
	}

	Performance struct {
		FormatVersion int             `json:"format_version"`
		Result        Result          `json:"result"`
		Data          PerformanceData `json:"data"`
	}
)
