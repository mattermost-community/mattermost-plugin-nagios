package nagios

import (
	"net/url"
	"reflect"
	"testing"
)

func TestHostStatus_String(t *testing.T) {
	tests := []struct {
		name       string
		hostStatus HostStatus
		want       string
	}{
		{
			name:       "none",
			hostStatus: HostStatus{},
			want:       "",
		},
		{
			name: "up only",
			hostStatus: HostStatus{
				Up:          true,
				Down:        false,
				Unreachable: false,
				Pending:     false,
			},
			want: "up",
		},
		{
			name: "down only",
			hostStatus: HostStatus{
				Up:          false,
				Down:        true,
				Unreachable: false,
				Pending:     false,
			},
			want: "down",
		},
		{
			name: "unreachable only",
			hostStatus: HostStatus{
				Up:          false,
				Down:        false,
				Unreachable: true,
				Pending:     false,
			},
			want: "unreachable",
		},
		{
			name: "pending only",
			hostStatus: HostStatus{
				Up:          false,
				Down:        false,
				Unreachable: false,
				Pending:     true,
			},
			want: "pending",
		},
		{
			name: "all",
			hostStatus: HostStatus{
				Up:          true,
				Down:        true,
				Unreachable: true,
				Pending:     true,
			},
			want: "up down unreachable pending",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hostStatus.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneralHostRequest_build(t *testing.T) {
	type args struct {
		query             string
		includeStartCount bool
	}
	tests := []struct {
		name               string
		generalHostRequest GeneralHostRequest
		args               args
		want               Query
	}{
		{
			name:               "blank",
			generalHostRequest: GeneralHostRequest{},
			args: args{
				query:             "",
				includeStartCount: false,
			},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{},
			},
		},
		{
			name:               "blank with Start and Count",
			generalHostRequest: GeneralHostRequest{Count: 1},
			args: args{
				query:             "",
				includeStartCount: true,
			},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"start": []string{"0"},
					"count": []string{"1"},
				},
			},
		},
		{
			name:               "blank with EndTime",
			generalHostRequest: GeneralHostRequest{EndTime: 1},
			args: args{
				query:             "",
				includeStartCount: false,
			},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"endtime": []string{"1"},
				},
			},
		},
		{
			name:               "blank with ShowDetails",
			generalHostRequest: GeneralHostRequest{ShowDetails: true},
			args: args{
				query:             "",
				includeStartCount: false,
			},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"details": []string{"true"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generalHostRequest.build(tt.args.query, tt.args.includeStartCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHostCountRequest_Build(t *testing.T) {
	tests := []struct {
		name               string
		GeneralHostRequest GeneralHostRequest
		want               Query
	}{
		{
			name:               "blank",
			GeneralHostRequest: GeneralHostRequest{},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"query": []string{"hostcount"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := HostCountRequest{
				GeneralHostRequest: tt.GeneralHostRequest,
			}
			if got := a.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHostListRequest_Build(t *testing.T) {
	tests := []struct {
		name               string
		GeneralHostRequest GeneralHostRequest
		want               Query
	}{
		{
			name:               "blank",
			GeneralHostRequest: GeneralHostRequest{},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"query": []string{"hostlist"},
					"start": []string{"0"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := HostListRequest{
				GeneralHostRequest: tt.GeneralHostRequest,
			}
			if got := a.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceStatus_String(t *testing.T) {
	tests := []struct {
		name          string
		serviceStatus ServiceStatus
		want          string
	}{
		{
			name:          "none",
			serviceStatus: ServiceStatus{},
			want:          "",
		},
		{
			name: "ok only",
			serviceStatus: ServiceStatus{
				Ok:       true,
				Warning:  false,
				Critical: false,
				Unknown:  false,
				Pending:  false,
			},
			want: "ok",
		},
		{
			name: "warning only",
			serviceStatus: ServiceStatus{
				Ok:       false,
				Warning:  true,
				Critical: false,
				Unknown:  false,
				Pending:  false,
			},
			want: "warning",
		},
		{
			name: "critical only",
			serviceStatus: ServiceStatus{
				Ok:       false,
				Warning:  false,
				Critical: true,
				Unknown:  false,
				Pending:  false,
			},
			want: "critical",
		},
		{
			name: "unknown only",
			serviceStatus: ServiceStatus{
				Ok:       false,
				Warning:  false,
				Critical: false,
				Unknown:  true,
				Pending:  false,
			},
			want: "unknown",
		},
		{
			name: "pending only",
			serviceStatus: ServiceStatus{
				Ok:       false,
				Warning:  false,
				Critical: false,
				Unknown:  false,
				Pending:  true,
			},
			want: "pending",
		},
		{
			name: "all",
			serviceStatus: ServiceStatus{
				Ok:       true,
				Warning:  true,
				Critical: true,
				Unknown:  true,
				Pending:  true,
			},
			want: "ok warning critical unknown pending",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.serviceStatus.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneralServiceRequest_build(t *testing.T) {
	type args struct {
		query             string
		includeStartCount bool
	}
	tests := []struct {
		name                  string
		generalServiceRequest GeneralServiceRequest
		args                  args
		want                  Query
	}{
		{
			name:                  "blank",
			generalServiceRequest: GeneralServiceRequest{},
			args: args{
				query:             "",
				includeStartCount: false,
			},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{},
			},
		},
		{
			name:                  "blank with Start and Count",
			generalServiceRequest: GeneralServiceRequest{Count: 1},
			args: args{
				query:             "",
				includeStartCount: true,
			},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"start": []string{"0"},
					"count": []string{"1"},
				},
			},
		},
		{
			name:                  "blank with EndTime",
			generalServiceRequest: GeneralServiceRequest{EndTime: 1},
			args: args{
				query:             "",
				includeStartCount: false,
			},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"endtime": []string{"1"},
				},
			},
		},
		{
			name:                  "blank with ShowDetails",
			generalServiceRequest: GeneralServiceRequest{ShowDetails: true},
			args: args{
				query:             "",
				includeStartCount: false,
			},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"details": []string{"true"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generalServiceRequest.build(tt.args.query, tt.args.includeStartCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceCountRequest_Build(t *testing.T) {
	tests := []struct {
		name                  string
		GeneralServiceRequest GeneralServiceRequest
		want                  Query
	}{
		{
			name:                  "blank",
			GeneralServiceRequest: GeneralServiceRequest{},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"query": []string{"servicecount"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ServiceCountRequest{
				GeneralServiceRequest: tt.GeneralServiceRequest,
			}
			if got := a.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceListRequest_Build(t *testing.T) {
	tests := []struct {
		name                  string
		GeneralServiceRequest GeneralServiceRequest
		want                  Query
	}{
		{
			name:                  "blank",
			GeneralServiceRequest: GeneralServiceRequest{},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"query": []string{"servicelist"},
					"start": []string{"0"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := ServiceListRequest{
				GeneralServiceRequest: tt.GeneralServiceRequest,
			}
			if got := a.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHostRequest_build(t *testing.T) {
	tests := []struct {
		name        string
		hostRequest HostRequest
		want        Query
	}{
		{
			name:        "blank",
			hostRequest: HostRequest{HostName: "localhost"},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"hostname": []string{"localhost"},
					"query":    []string{"host"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hostRequest.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceRequest_build(t *testing.T) {
	tests := []struct {
		name           string
		serviceRequest ServiceRequest
		want           Query
	}{
		{
			name: "blank",
			serviceRequest: ServiceRequest{
				HostName:           "localhost",
				ServiceDescription: "http",
			},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"hostname":           []string{"localhost"},
					"servicedescription": []string{"http"},
					"query":              []string{"service"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.serviceRequest.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPerformanceDataRequest_Build(t *testing.T) {
	tests := []struct {
		name                   string
		performanceDataRequest PerformanceDataRequest
		want                   Query
	}{
		{
			name:                   "blank",
			performanceDataRequest: PerformanceDataRequest{},
			want: Query{
				Endpoint: statusEndpoint,
				URLQuery: url.Values{
					"query": []string{"performancedata"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.performanceDataRequest.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("build() = %v, want %v", got, tt.want)
			}
		})
	}
}
