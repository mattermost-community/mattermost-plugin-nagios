package main

import (
	"testing"

	"github.com/mattermost/mattermost-plugin-nagios/go-nagios/nagios"
)

func Test_formatHostName(t *testing.T) {
	type args struct {
		name string
		alt  string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "use name",
			args: args{
				name: "test",
				alt:  "test alt",
			},
			want: "test",
		},
		{
			name: "use alt",
			args: args{
				name: "",
				alt:  "test alt",
			},
			want: "test alt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatHostName(tt.args.name, tt.args.alt); got != tt.want {
				t.Errorf("formatHostName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gettingLogsUnsuccessfulMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "basic",
			message: "test",
			want:    gettingLogsUnsuccessful + ": test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gettingLogsUnsuccessfulMessage(tt.message); got != tt.want {
				t.Errorf("gettingLogsUnsuccessfulMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unknownParameterMessage(t *testing.T) {
	tests := []struct {
		name      string
		parameter string
		want      string
	}{
		{
			name:      "basic",
			parameter: "test",
			want:      "Unknown parameter (test).",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unknownParameterMessage(tt.parameter); got != tt.want {
				t.Errorf("unknownParameterMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatAlerts(t *testing.T) {
	tests := []struct {
		name   string
		alerts nagios.AlertList
		want   string
	}{
		{
			name: "empty",
			alerts: nagios.AlertList{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
			},
			want: "No alerts.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatAlerts(tt.alerts); got != tt.want {
				t.Errorf("formatAlerts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatNotifications(t *testing.T) {
	tests := []struct {
		name          string
		notifications nagios.NotificationList
		want          string
	}{
		{
			name: "empty",
			notifications: nagios.NotificationList{
				Result: nagios.Result{
					TypeText: resultTypeTextSuccess,
				},
			},
			want: "No notifications.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatNotifications(tt.notifications); got != tt.want {
				t.Errorf("formatNotifications() = %v, want %v", got, tt.want)
			}
		})
	}
}
