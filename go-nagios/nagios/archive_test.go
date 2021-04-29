package nagios

import (
	"net/url"
	"reflect"
	"testing"
)

func TestObjectTypes_String(t *testing.T) {
	tests := []struct {
		name        string
		objectTypes ObjectTypes
		want        string
	}{
		{
			name:        "none",
			objectTypes: ObjectTypes{},
			want:        "",
		},
		{
			name: "host only",
			objectTypes: ObjectTypes{
				Host:    true,
				Service: false,
			},
			want: "host",
		},
		{
			name: "service only",
			objectTypes: ObjectTypes{
				Host:    false,
				Service: true,
			},
			want: "service",
		},
		{
			name: "all",
			objectTypes: ObjectTypes{
				Host:    true,
				Service: true,
			},
			want: "host service",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.objectTypes.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateTypes_String(t *testing.T) {
	tests := []struct {
		name       string
		stateTypes StateTypes
		want       string
	}{
		{
			name:       "none",
			stateTypes: StateTypes{},
			want:       "",
		},
		{
			name: "soft only",
			stateTypes: StateTypes{
				Soft: true,
				Hard: false,
			},
			want: "soft",
		},
		{
			name: "hard only",
			stateTypes: StateTypes{
				Soft: false,
				Hard: true,
			},
			want: "hard",
		},
		{
			name: "all",
			stateTypes: StateTypes{
				Soft: true,
				Hard: true,
			},
			want: "soft hard",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.stateTypes.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHostStates_String(t *testing.T) {
	tests := []struct {
		name       string
		hostStates HostStates
		want       string
	}{
		{
			name:       "none",
			hostStates: HostStates{},
			want:       "",
		},
		// TODO(amwolff): add test cases.
		{
			name: "all",
			hostStates: HostStates{
				Up:          true,
				Down:        true,
				Unreachable: true,
			},
			want: "up down unreachable",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hostStates.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceStates_String(t *testing.T) {
	tests := []struct {
		name          string
		serviceStates ServiceStates
		want          string
	}{
		{
			name:          "none",
			serviceStates: ServiceStates{},
			want:          "",
		},
		// TODO(amwolff): add test cases.
		{
			name: "all",
			serviceStates: ServiceStates{
				Ok:       true,
				Warning:  true,
				Critical: true,
				Unknown:  true,
			},
			want: "ok warning critical unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.serviceStates.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneralAlertRequest_build(t *testing.T) {
	type args struct {
		query             string
		includeStartCount bool
	}

	tests := []struct {
		name                string
		generalAlertRequest GeneralAlertRequest
		args                args
		want                Query
	}{
		{
			name:                "blank",
			generalAlertRequest: GeneralAlertRequest{},
			args: args{
				query:             "",
				includeStartCount: false,
			},
			want: Query{
				Endpoint: archiveEndpoint,
				URLQuery: url.Values{
					"starttime": []string{"0"},
					"endtime":   []string{"0"},
				},
			},
		},
		{
			name:                "blank with Start and Count",
			generalAlertRequest: GeneralAlertRequest{Count: 1},
			args: args{
				query:             "",
				includeStartCount: true,
			},
			want: Query{
				Endpoint: archiveEndpoint,
				URLQuery: url.Values{
					"start":     []string{"0"},
					"count":     []string{"1"},
					"starttime": []string{"0"},
					"endtime":   []string{"0"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generalAlertRequest.build(tt.args.query, tt.args.includeStartCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlertCountRequest_Build(t *testing.T) {
	tests := []struct {
		name                string
		GeneralAlertRequest GeneralAlertRequest
		want                Query
	}{
		{
			name:                "blank",
			GeneralAlertRequest: GeneralAlertRequest{},
			want: Query{
				Endpoint: archiveEndpoint,
				URLQuery: url.Values{
					"query":     []string{"alertcount"},
					"starttime": []string{"0"},
					"endtime":   []string{"0"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AlertCountRequest{
				GeneralAlertRequest: tt.GeneralAlertRequest,
			}
			if got := a.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlertListRequest_Build(t *testing.T) {
	tests := []struct {
		name                string
		GeneralAlertRequest GeneralAlertRequest
		want                Query
	}{
		{
			name:                "blank",
			GeneralAlertRequest: GeneralAlertRequest{},
			want: Query{
				Endpoint: archiveEndpoint,
				URLQuery: url.Values{
					"query":     []string{"alertlist"},
					"start":     []string{"0"},
					"starttime": []string{"0"},
					"endtime":   []string{"0"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AlertListRequest{
				GeneralAlertRequest: tt.GeneralAlertRequest,
			}
			if got := a.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHostNotificationTypes_String(t *testing.T) {
	tests := []struct {
		name                  string
		hostNotificationTypes HostNotificationTypes
		want                  string
	}{
		{
			name:                  "none",
			hostNotificationTypes: HostNotificationTypes{},
			want:                  "",
		},
		// TODO(amwolff): add test cases.
		{
			name: "all",
			hostNotificationTypes: HostNotificationTypes{
				NoData:        true,
				Down:          true,
				Unreachable:   true,
				Recovery:      true,
				HostCustom:    true,
				HostAck:       true,
				HostFlapStart: true,
				HostFlapStop:  true,
			},
			want: "nodata down unreachable recovery hostcustom hostack hostflapstart hostflapstop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hostNotificationTypes.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceNotificationTypes_String(t *testing.T) {
	tests := []struct {
		name                     string
		serviceNotificationTypes ServiceNotificationTypes
		want                     string
	}{
		{
			name:                     "none",
			serviceNotificationTypes: ServiceNotificationTypes{},
			want:                     "",
		},
		// TODO(amwolff): add test cases.
		{
			name: "all",
			serviceNotificationTypes: ServiceNotificationTypes{
				NoData:           true,
				Critical:         true,
				Warning:          true,
				Recovery:         true,
				Custom:           true,
				ServiceAck:       true,
				ServiceFlapStart: true,
				ServiceFlapStop:  true,
				Unknown:          true,
			},
			want: "nodata critical warning recovery custom serviceack serviceflapstart serviceflapstop unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.serviceNotificationTypes.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneralNotificationRequest_build(t *testing.T) {
	type args struct {
		query             string
		includeStartCount bool
	}

	tests := []struct {
		name                       string
		generalNotificationRequest GeneralNotificationRequest
		args                       args
		want                       Query
	}{
		{
			name:                       "blank",
			generalNotificationRequest: GeneralNotificationRequest{},
			args: args{
				query:             "",
				includeStartCount: false,
			},
			want: Query{
				Endpoint: archiveEndpoint,
				URLQuery: url.Values{
					"starttime": []string{"0"},
					"endtime":   []string{"0"},
				},
			},
		},
		{
			name:                       "blank with Start and Count",
			generalNotificationRequest: GeneralNotificationRequest{Count: 1},
			args: args{
				query:             "",
				includeStartCount: true,
			},
			want: Query{
				Endpoint: archiveEndpoint,
				URLQuery: url.Values{
					"start":     []string{"0"},
					"count":     []string{"1"},
					"starttime": []string{"0"},
					"endtime":   []string{"0"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generalNotificationRequest.build(tt.args.query, tt.args.includeStartCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotificationCountRequest_Build(t *testing.T) {
	tests := []struct {
		name                       string
		GeneralNotificationRequest GeneralNotificationRequest
		want                       Query
	}{
		{
			name:                       "blank",
			GeneralNotificationRequest: GeneralNotificationRequest{},
			want: Query{
				Endpoint: archiveEndpoint,
				URLQuery: url.Values{
					"query":     []string{"notificationcount"},
					"starttime": []string{"0"},
					"endtime":   []string{"0"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NotificationCountRequest{
				GeneralNotificationRequest: tt.GeneralNotificationRequest,
			}
			if got := n.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotificationListRequest_Build(t *testing.T) {
	tests := []struct {
		name                       string
		GeneralNotificationRequest GeneralNotificationRequest
		want                       Query
	}{
		{
			name:                       "blank",
			GeneralNotificationRequest: GeneralNotificationRequest{},
			want: Query{
				Endpoint: archiveEndpoint,
				URLQuery: url.Values{
					"query":     []string{"notificationlist"},
					"start":     []string{"0"},
					"starttime": []string{"0"},
					"endtime":   []string{"0"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NotificationListRequest{
				GeneralNotificationRequest: tt.GeneralNotificationRequest,
			}
			if got := n.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
