package main

import (
	"testing"
)

func Test_configuration_isValid(t *testing.T) {
	type fields struct {
		NagiosURL              string
		Token                  string
		InitialLogsLimit       int
		InitialLogsStartTime   int
		InitialReportFrequency int
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "URL only",
			fields: fields{
				NagiosURL: "https://example.com",
			},
			wantErr: true,
		},
		{
			name: "URL and token only",
			fields: fields{
				NagiosURL: "https://example.com",
				Token:     "test",
			},
			wantErr: true,
		},
		{
			name: "URL and logs limit only",
			fields: fields{
				NagiosURL:        "https://example.com",
				InitialLogsLimit: 10,
			},
			wantErr: true,
		},
		{
			name: "URL, token, logs limit and logs start time only",
			fields: fields{
				NagiosURL:            "https://example.com",
				Token:                "test",
				InitialLogsLimit:     100,
				InitialLogsStartTime: 1000,
			},
			wantErr: true,
		},
		{
			name: "valid configuration (without token)",
			fields: fields{
				NagiosURL:              "https://example.com",
				InitialLogsLimit:       5,
				InitialLogsStartTime:   25,
				InitialReportFrequency: 50,
			},
			wantErr: false,
		},
		{
			name: "valid configuration 1 (with token)",
			fields: fields{
				NagiosURL:              "https://example.com",
				Token:                  "test",
				InitialLogsLimit:       5,
				InitialLogsStartTime:   25,
				InitialReportFrequency: 50,
			},
			wantErr: false,
		},
		{
			name: "invalid configuration (URL)",
			fields: fields{
				NagiosURL:              "",
				InitialLogsLimit:       5,
				InitialLogsStartTime:   25,
				InitialReportFrequency: 50,
			},
			wantErr: true,
		},
		{
			name: "invalid configuration (logs limit)",
			fields: fields{
				NagiosURL:              "https://example.com",
				Token:                  "test",
				InitialLogsLimit:       0,
				InitialLogsStartTime:   25,
				InitialReportFrequency: 50,
			},
			wantErr: true,
		},
		{
			name: "invalid configuration (logs start time)",
			fields: fields{
				NagiosURL:              "https://example.com",
				InitialLogsLimit:       5,
				InitialLogsStartTime:   -0xC00FEE,
				InitialReportFrequency: 50,
			},
			wantErr: true,
		},
		{
			name: "invalid configuration (report frequency)",
			fields: fields{
				NagiosURL:              "https://example.com",
				Token:                  "test",
				InitialLogsLimit:       5,
				InitialLogsStartTime:   25,
				InitialReportFrequency: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &configuration{
				NagiosURL:              tt.fields.NagiosURL,
				Token:                  tt.fields.Token,
				InitialLogsLimit:       tt.fields.InitialLogsLimit,
				InitialLogsStartTime:   tt.fields.InitialLogsStartTime,
				InitialReportFrequency: tt.fields.InitialReportFrequency,
			}
			if err := c.isValid(); (err != nil) != tt.wantErr {
				t.Errorf("configuration.isValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
