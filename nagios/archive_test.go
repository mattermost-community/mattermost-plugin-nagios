package nagios

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_buildOptions(t *testing.T) {
	tests := []struct {
		name    string
		options []string
		want    string
	}{
		{
			name:    "none",
			options: nil,
			want:    "",
		},
		{
			name:    "some",
			options: []string{"abc", "def"},
			want:    "abc def",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildOptions(tt.options); got != tt.want {
				t.Errorf("buildOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatOptions_String(t *testing.T) {
	tests := []struct {
		name          string
		formatOptions FormatOptions
		want          string
	}{
		{
			name:          "none",
			formatOptions: FormatOptions{},
			want:          "",
		},
		{
			name: "few",
			formatOptions: FormatOptions{
				Whitespace: true,
				Enumerate:  false,
				Bitmask:    true,
				Duration:   false,
			},
			want: "whitespace bitmask",
		},
		{
			name: "all",
			formatOptions: FormatOptions{
				Whitespace: true,
				Enumerate:  true,
				Bitmask:    true,
				Duration:   true,
			},
			want: "whitespace enumerate bitmask duration",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.formatOptions.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			name: "few",
			objectTypes: ObjectTypes{
				Host:    true,
				Service: false,
			},
			want: "host",
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
			name: "few",
			stateTypes: StateTypes{
				Soft: true,
				Hard: false,
			},
			want: "soft",
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
		{
			name: "few",
			hostStates: HostStates{
				Up:          true,
				Down:        false,
				Unreachable: true,
			},
			want: "up unreachable",
		},
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
		{
			name: "few",
			serviceStates: ServiceStates{
				Ok:       true,
				Warning:  false,
				Critical: true,
				Unknown:  false,
			},
			want: "ok critical",
		},
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

func Test_alertRequest_build(t *testing.T) {
	tests := []struct {
		name              string
		alertRequest      alertRequest
		includeStartCount bool
		want              Query
	}{
		{
			name:              "none",
			alertRequest:      alertRequest{},
			includeStartCount: false,
			want: Query{
				Endpoint: archiveEndpoint,
				URLQuery: url.Values{
					"starttime": []string{"0"},
					"endtime":   []string{"0"},
				},
			},
		},
		{
			name:              "none with Start and Count",
			alertRequest:      alertRequest{},
			includeStartCount: true,
			want: Query{
				Endpoint: archiveEndpoint,
				URLQuery: url.Values{
					"start":     []string{"0"},
					"count":     []string{"0"},
					"starttime": []string{"0"},
					"endtime":   []string{"0"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.alertRequest.build(tt.includeStartCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("build() = %v, want %v", got, tt.want)
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
		{
			name: "few",
			hostNotificationTypes: HostNotificationTypes{
				NoData:        true,
				Down:          false,
				Unreachable:   true,
				Recovery:      false,
				HostCustom:    true,
				HostAck:       false,
				HostFlapStart: true,
				HostFlapStop:  false,
			},
			want: "nodata unreachable hostcustom hostflapstart",
		},
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
		{
			name: "few",
			serviceNotificationTypes: ServiceNotificationTypes{
				NoData:           true,
				Critical:         false,
				Warning:          true,
				Recovery:         false,
				Custom:           true,
				ServiceAck:       false,
				ServiceFlapStart: true,
				ServiceFlapStop:  false,
				Unknown:          true,
			},
			want: "nodata warning custom serviceflapstart unknown",
		},
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

func Test_notificationRequest_build(t *testing.T) {
	tests := []struct {
		name                string
		notificationRequest notificationRequest
		includeStartCount   bool
		want                Query
	}{
		{
			name:                "none",
			notificationRequest: notificationRequest{},
			includeStartCount:   false,
			want: Query{
				Endpoint: archiveEndpoint,
				URLQuery: url.Values{
					"starttime": []string{"0"},
					"endtime":   []string{"0"},
				},
			},
		},
		{
			name:                "none with Start and Count",
			notificationRequest: notificationRequest{},
			includeStartCount:   true,
			want: Query{
				Endpoint: archiveEndpoint,
				URLQuery: url.Values{
					"start":     []string{"0"},
					"count":     []string{"0"},
					"starttime": []string{"0"},
					"endtime":   []string{"0"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.notificationRequest.build(tt.includeStartCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("build() = %v, want %v", got, tt.want)
			}
		})
	}
}
