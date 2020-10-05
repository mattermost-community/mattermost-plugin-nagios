package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ulumuri/go-nagios/nagios"
)

func Test_gettingReportUnsuccessfulMessage(t *testing.T) {
	type args struct {
		reportPart string
		message    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				reportPart: "a part",
				message:    "a message",
			},
			want: "Getting monitoring report unsuccessful (a part): a message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gettingReportUnsuccessfulMessage(tt.args.reportPart, tt.args.message); got != tt.want {
				t.Errorf("gettingReportUnsuccessfulMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reportPreamble(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		t    time.Time
		want string
	}{
		{
			name: "basic",
			t:    now,
			want: "#### :bar_chart: System monitoring report (" + now.Format(time.UnixDate) + ")\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reportPreamble(tt.t); got != tt.want {
				t.Errorf("reportPreamble() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatHostCount(t *testing.T) {
	tests := []struct {
		name  string
		count nagios.HostCount
		want  string
	}{
		{
			name: "basic",
			count: nagios.HostCount{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
				Data: nagios.HostCountData{
					Count: nagios.HostStatusCount{
						Up:          1,
						Down:        2,
						Unreachable: 3,
						Pending:     4,
					},
				},
			},
			want: "##### HOST SUMMARY\n\n:up: Up: **1**  :small_red_triangle_" +
				"down: Down: **2**  :mailbox_with_no_mail: Unreachable: **3**" +
				"  :hourglass_flowing_sand: Pending: **4**",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatHostCount(tt.count); got != tt.want {
				t.Errorf("formatHostCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func generateHomogenousHostListData(state string, n int) nagios.HostListData {
	ret := nagios.HostListData{
		HostList: make(map[string]json.RawMessage),
	}

	m := json.RawMessage(fmt.Sprintf(`"%s"`, state))

	for i := 0; i < n; i++ {
		ret.HostList[strconv.Itoa(i)] = m
	}

	return ret
}

// Test_formatHostList only tests the edge cases to make sure we don't hit the
// Mattermost's message limits.
func Test_formatHostList(t *testing.T) {
	tests := []struct {
		name string
		list nagios.HostList
		want string
	}{
		{
			name: "empty",
			list: nagios.HostList{},
			want: gettingReportUnsuccessfulMessage("host list", ""),
		},
		{
			name: "empty successful",
			list: nagios.HostList{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
			},
			want: "##### HOST LIST\n\nNo hosts to show.",
		},
		{
			name: "too many hosts (all UP)",
			list: nagios.HostList{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
				Data: generateHomogenousHostListData(upState, maximumReportLength),
			},
			want: "##### HOST LIST\n\n**Too many hosts. Showing only abnormal" +
				" state hosts.**\n\nNo hosts to show.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatHostList(tt.list); got != tt.want {
				t.Errorf("formatHostList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatServiceCount(t *testing.T) {
	tests := []struct {
		name  string
		count nagios.ServiceCount
		want  string
	}{
		{
			name: "basic",
			count: nagios.ServiceCount{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
				Data: nagios.ServiceCountData{
					Count: nagios.ServiceStatusCount{
						Ok:       1,
						Warning:  2,
						Critical: 3,
						Unknown:  4,
						Pending:  5,
					},
				},
			},
			want: "##### SERVICE SUMMARY\n\n:white_check_mark: OK: **1**  :wa" +
				"rning: Warning: **2**  :bangbang: Critical: **3**  :question" +
				": Unknown: **4**  :hourglass_flowing_sand: Pending: **5**",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatServiceCount(tt.count); got != tt.want {
				t.Errorf("formatServiceCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func extractedServicesEqual(a, b []extractedService) bool {
	m := make(map[string]string)

	for _, v := range a {
		m[v.name] = v.state
	}

	for _, v := range b {
		if v.state != m[v.name] {
			return false
		}
	}

	return true
}

func Test_extractServices(t *testing.T) {
	tests := []struct {
		name       string
		rawMessage json.RawMessage
		want       []extractedService
	}{
		{
			name: "part of a real response",
			rawMessage: []byte(`{"Bandwidth Spike":"ok","Ping":"ok","Port 1 B` +
				`andwidth":"ok","Port 1 Status":"ok","Port 10 Bandwidth":"ok"` +
				`,"Port 10 Status":"critical","Port 11 Bandwidth":"ok","Port ` +
				`11 Status":"critical","Port 12 Bandwidth":"ok","Port 12 Stat` +
				`us":"warning","Port 13 Bandwidth":"ok","Port 13 Status":"cri` +
				`tical","Port 14 Bandwidth":"ok","Port 14 Status":"critical",` +
				`"Port 15 Bandwidth":"ok","Port 15 Status":"critical","Port 1` +
				`6 Bandwidth":"ok","Port 16 Status":"warning","Port 17 Bandwi` +
				`dth":"ok","Port 17 Status":"warning","Port 18 Bandwidth":"ok` +
				`","Port 18 Status":"critical","Port 19 Bandwidth":"ok","Port` +
				` 19 Status":"critical","Port 2 Bandwidth":"ok","Port 2 Statu` +
				`s":"critical","Port 20 Bandwidth":"ok","Port 20 Status":"ok"` +
				`,"Port 21 Bandwidth":"ok","Port 21 Status":"ok","Port 22 Ban` +
				`dwidth":"ok","Port 22 Status":"warning","Port 23 Bandwidth":` +
				`"ok","Port 23 Status":"warning","Port 24 Bandwidth":"ok","Po` +
				`rt 24 Status":"ok","Port 25 Bandwidth":"ok","Port 25 Status"` +
				`:"ok","Port 3 Bandwidth":"ok","Port 3 Status":"ok","Port 4 B` +
				`andwidth":"ok","Port 4 Status":"warning","Port 5 Bandwidth":` +
				`"ok","Port 5 Status":"ok","Port 6 Bandwidth":"ok","Port 6 St` +
				`atus":"critical","Port 7 Bandwidth":"ok","Port 7 Status":"cr` +
				`itical","Port 8 Bandwidth":"ok","Port 8 Status":"critical","` +
				`Port 9 Bandwidth":"ok","Port 9 Status":"ok","Youtube Usage":` +
				`"warning"}`),
			want: []extractedService{
				{
					name:  "Bandwidth Spike",
					state: okState,
				},
				{
					name:  "Ping",
					state: okState,
				},
				{
					name:  "Port 1 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 1 Status",
					state: okState,
				},
				{
					name:  "Port 10 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 10 Status",
					state: criticalState,
				},
				{
					name:  "Port 11 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 11 Status",
					state: criticalState,
				},
				{
					name:  "Port 12 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 12 Status",
					state: warningState,
				},
				{
					name:  "Port 13 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 13 Status",
					state: criticalState,
				},
				{
					name:  "Port 14 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 14 Status",
					state: criticalState,
				},
				{
					name:  "Port 15 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 15 Status",
					state: criticalState,
				},
				{
					name:  "Port 16 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 16 Status",
					state: warningState,
				},
				{
					name:  "Port 17 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 17 Status",
					state: warningState,
				},
				{
					name:  "Port 18 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 18 Status",
					state: criticalState,
				},
				{
					name:  "Port 19 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 19 Status",
					state: criticalState,
				},
				{
					name:  "Port 2 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 2 Status",
					state: criticalState,
				},
				{
					name:  "Port 20 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 20 Status",
					state: okState,
				},
				{
					name:  "Port 21 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 21 Status",
					state: okState,
				},
				{
					name:  "Port 22 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 22 Status",
					state: warningState,
				},
				{
					name:  "Port 23 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 23 Status",
					state: warningState,
				},
				{
					name:  "Port 24 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 24 Status",
					state: okState,
				},
				{
					name:  "Port 25 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 25 Status",
					state: okState,
				},
				{
					name:  "Port 3 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 3 Status",
					state: okState,
				},
				{
					name:  "Port 4 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 4 Status",
					state: warningState,
				},
				{
					name:  "Port 5 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 5 Status",
					state: okState,
				},
				{
					name:  "Port 6 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 6 Status",
					state: criticalState,
				},
				{
					name:  "Port 7 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 7 Status",
					state: criticalState,
				},
				{
					name:  "Port 8 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 8 Status",
					state: criticalState,
				},
				{
					name:  "Port 9 Bandwidth",
					state: okState,
				},
				{
					name:  "Port 9 Status",
					state: okState,
				},
				{
					name:  "Youtube Usage",
					state: warningState,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractServices(tt.rawMessage); !extractedServicesEqual(got, tt.want) {
				t.Errorf("extractServices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mustGenerateHomogenousServiceListData(state string, outer, inner int) nagios.ServiceListData {
	services := make(map[string]json.RawMessage)

	for i := 0; i < inner; i++ {
		services[strconv.Itoa(i)] = json.RawMessage(fmt.Sprintf(`"%s"`, state))
	}

	rawInner, err := json.Marshal(services)
	if err != nil {
		panic(fmt.Sprintf("Marshal: %v", err))
	}

	ret := nagios.ServiceListData{
		ServiceList: make(map[string]json.RawMessage),
	}

	for i := 0; i < outer; i++ {
		ret.ServiceList[strconv.Itoa(i)] = rawInner
	}

	return ret
}

// Test_formatServiceList only tests the edge cases to make sure we don't hit
// the Mattermost's message limits.
func Test_formatServiceList(t *testing.T) {
	tests := []struct {
		name string
		list nagios.ServiceList
		want string
	}{
		{
			name: "empty",
			list: nagios.ServiceList{},
			want: gettingReportUnsuccessfulMessage("service list", ""),
		},
		{
			name: "empty successful",
			list: nagios.ServiceList{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
			},
			want: "##### SERVICE LIST\n\nNo services to show.",
		},
		{
			name: "too many services (all OK)",
			list: nagios.ServiceList{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
				Data: mustGenerateHomogenousServiceListData(okState, 1, maximumReportLength-1),
			},
			want: "##### SERVICE LIST\n\n**Too many services. Showing only ab" +
				"normal state services.**\n\nNo services to show.",
		},
		{
			name: "too many services (all OK [hosts only])",
			list: nagios.ServiceList{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
				Data: mustGenerateHomogenousServiceListData(okState, maximumReportLength, 0),
			},
			want: "##### SERVICE LIST\n\n**Too many services. Showing only ab" +
				"normal state services.**\n\nNo services to show.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatServiceList(tt.list); got != tt.want {
				t.Errorf("formatServiceList() = %v, want %v", got, tt.want)
			}
		})
	}
}
