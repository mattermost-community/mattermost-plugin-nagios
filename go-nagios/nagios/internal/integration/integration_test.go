package integration

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ulumuri/go-nagios/nagios"
)

const (
	success      = "Success"
	dumpResponse = false
)

func addr(t *testing.T, address string) string {
	if len(address) == 0 {
		if address = os.Getenv("TEST_INSTANCE_ADDRESS"); len(address) == 0 {
			t.Skip()
		}
	}
	return address
}

func TestArchive(t *testing.T) {
	c, err := nagios.NewClient(http.DefaultClient, addr(t, testInstanceAddress))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	now := time.Now()
	then := now.Add(-24 * time.Hour)

	t.Run("blank alert count", func(t *testing.T) {
		req := nagios.AlertCountRequest{
			GeneralAlertRequest: nagios.GeneralAlertRequest{
				FormatOptions: nagios.FormatOptions{
					Enumerate: true,
				},
				StartTime: then.Unix(),
				EndTime:   now.Unix(),
			},
		}

		var count nagios.AlertCount

		if err := c.Query(req, &count); err != nil {
			t.Errorf("Query: %v", err)
		}

		if count.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(count)
		}
	})

	t.Run("blank alert count with options switched", func(t *testing.T) {
		req := nagios.AlertCountRequest{
			GeneralAlertRequest: nagios.GeneralAlertRequest{
				FormatOptions: nagios.FormatOptions{
					Whitespace: true,
					Enumerate:  true,
					Bitmask:    true,
					Duration:   true,
				},
				ObjectTypes: nagios.ObjectTypes{
					Host:    true,
					Service: true,
				},
				StateTypes: nagios.StateTypes{
					Soft: true,
					Hard: true,
				},
				HostStates: nagios.HostStates{
					Up:          true,
					Down:        true,
					Unreachable: true,
				},
				ServiceStates: nagios.ServiceStates{
					Ok:       true,
					Warning:  true,
					Critical: true,
					Unknown:  true,
				},
				StartTime: then.Unix(),
				EndTime:   now.Unix(),
			},
		}

		var count nagios.AlertCount

		if err := c.Query(req, &count); err != nil {
			t.Errorf("Query: %v", err)
		}

		if count.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(count)
		}
	})

	t.Run("blank alert list", func(t *testing.T) {
		req := nagios.AlertListRequest{
			GeneralAlertRequest: nagios.GeneralAlertRequest{
				FormatOptions: nagios.FormatOptions{
					Enumerate: true,
				},
				StartTime: then.Unix(),
				EndTime:   now.Unix(),
			},
		}

		var list nagios.AlertList

		if err := c.Query(req, &list); err != nil {
			t.Errorf("Query: %v", err)
		}

		if list.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(list)
		}
	})

	t.Run("blank alert list with options switched", func(t *testing.T) {
		req := nagios.AlertListRequest{
			GeneralAlertRequest: nagios.GeneralAlertRequest{
				FormatOptions: nagios.FormatOptions{
					Whitespace: true,
					Enumerate:  true,
					Bitmask:    true,
					Duration:   true,
				},
				ObjectTypes: nagios.ObjectTypes{
					Host:    true,
					Service: true,
				},
				StateTypes: nagios.StateTypes{
					Soft: true,
					Hard: true,
				},
				HostStates: nagios.HostStates{
					Up:          true,
					Down:        true,
					Unreachable: true,
				},
				ServiceStates: nagios.ServiceStates{
					Ok:       true,
					Warning:  true,
					Critical: true,
					Unknown:  true,
				},
				StartTime: then.Unix(),
				EndTime:   now.Unix(),
			},
		}

		var list nagios.AlertList

		if err := c.Query(req, &list); err != nil {
			t.Errorf("Query: %v", err)
		}

		if list.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(list)
		}
	})

	t.Run("blank notification count", func(t *testing.T) {
		req := nagios.NotificationCountRequest{
			GeneralNotificationRequest: nagios.GeneralNotificationRequest{
				FormatOptions: nagios.FormatOptions{
					Enumerate: true,
				},
				StartTime: then.Unix(),
				EndTime:   now.Unix(),
			},
		}

		var count nagios.NotificationCount

		if err := c.Query(req, &count); err != nil {
			t.Errorf("Query: %v", err)
		}

		if count.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(count)
		}
	})

	t.Run("blank notification count with options switched", func(t *testing.T) {
		req := nagios.NotificationCountRequest{
			GeneralNotificationRequest: nagios.GeneralNotificationRequest{
				FormatOptions: nagios.FormatOptions{
					Whitespace: true,
					Enumerate:  true,
					Bitmask:    true,
					Duration:   true,
				},
				ObjectTypes: nagios.ObjectTypes{
					Host:    true,
					Service: true,
				},
				HostNotificationTypes: nagios.HostNotificationTypes{
					NoData:        true,
					Down:          true,
					Unreachable:   true,
					Recovery:      true,
					HostCustom:    true,
					HostAck:       true,
					HostFlapStart: true,
					HostFlapStop:  true,
				},
				ServiceNotificationTypes: nagios.ServiceNotificationTypes{
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
				StartTime: then.Unix(),
				EndTime:   now.Unix(),
			},
		}

		var count nagios.NotificationCount

		if err := c.Query(req, &count); err != nil {
			t.Errorf("Query: %v", err)
		}

		if count.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(count)
		}
	})

	t.Run("blank notification list", func(t *testing.T) {
		req := nagios.NotificationListRequest{
			GeneralNotificationRequest: nagios.GeneralNotificationRequest{
				FormatOptions: nagios.FormatOptions{
					Enumerate: true,
				},
				StartTime: then.Unix(),
				EndTime:   now.Unix(),
			},
		}

		var list nagios.NotificationList

		if err := c.Query(req, &list); err != nil {
			t.Errorf("Query: %v", err)
		}

		if list.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(list)
		}
	})

	t.Run("blank notification list with options switched", func(t *testing.T) {
		req := nagios.NotificationListRequest{
			GeneralNotificationRequest: nagios.GeneralNotificationRequest{
				FormatOptions: nagios.FormatOptions{
					Whitespace: true,
					Enumerate:  true,
					Bitmask:    true,
					Duration:   true,
				},
				ObjectTypes: nagios.ObjectTypes{
					Host:    true,
					Service: true,
				},
				HostNotificationTypes: nagios.HostNotificationTypes{
					NoData:        true,
					Down:          true,
					Unreachable:   true,
					Recovery:      true,
					HostCustom:    true,
					HostAck:       true,
					HostFlapStart: true,
					HostFlapStop:  true,
				},
				ServiceNotificationTypes: nagios.ServiceNotificationTypes{
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
				StartTime: then.Unix(),
				EndTime:   now.Unix(),
			},
		}

		var list nagios.NotificationList

		if err := c.Query(req, &list); err != nil {
			t.Errorf("Query: %v", err)
		}

		if list.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(list)
		}
	})
}

func TestStatus(t *testing.T) {
	c, err := nagios.NewClient(http.DefaultClient, addr(t, testInstanceAddress))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	t.Run("blank host count", func(t *testing.T) {
		req := nagios.HostCountRequest{
			GeneralHostRequest: nagios.GeneralHostRequest{
				FormatOptions: nagios.FormatOptions{
					Enumerate: true,
				},
			},
		}

		var count nagios.HostCount

		if err := c.Query(req, &count); err != nil {
			t.Errorf("Query: %v", err)
		}

		if count.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(count)
		}
	})

	t.Run("blank host count with options switched", func(t *testing.T) {
		req := nagios.HostCountRequest{
			GeneralHostRequest: nagios.GeneralHostRequest{
				FormatOptions: nagios.FormatOptions{
					Whitespace: true,
					Enumerate:  true,
					Bitmask:    true,
					Duration:   true,
				},
				HostStatus: nagios.HostStatus{
					Up:          true,
					Down:        true,
					Unreachable: true,
					Pending:     true,
				},
			},
		}

		var count nagios.HostCount

		if err := c.Query(req, &count); err != nil {
			t.Errorf("Query: %v", err)
		}

		if count.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(count)
		}
	})

	t.Run("blank host list", func(t *testing.T) {
		req := nagios.HostListRequest{
			GeneralHostRequest: nagios.GeneralHostRequest{
				FormatOptions: nagios.FormatOptions{
					Enumerate: true,
				},
			},
		}

		var list nagios.HostList

		if err := c.Query(req, &list); err != nil {
			t.Errorf("Query: %v", err)
		}

		if list.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(list)
		}
	})

	t.Run("blank host list with options switched", func(t *testing.T) {
		req := nagios.HostListRequest{
			GeneralHostRequest: nagios.GeneralHostRequest{
				FormatOptions: nagios.FormatOptions{
					Whitespace: true,
					Enumerate:  true,
					Bitmask:    true,
					Duration:   true,
				},
				HostStatus: nagios.HostStatus{
					Up:          true,
					Down:        true,
					Unreachable: true,
					Pending:     true,
				},
				ShowDetails: true,
			},
		}

		var list nagios.HostList

		if err := c.Query(req, &list); err != nil {
			t.Errorf("Query: %v", err)
		}

		if list.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(list)
		}
	})

	t.Run("blank service count", func(t *testing.T) {
		req := nagios.ServiceCountRequest{
			GeneralServiceRequest: nagios.GeneralServiceRequest{
				FormatOptions: nagios.FormatOptions{
					Enumerate: true,
				},
			},
		}

		var count nagios.ServiceCount

		if err := c.Query(req, &count); err != nil {
			t.Errorf("Query: %v", err)
		}

		if count.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(count)
		}
	})

	t.Run("blank service count with options switched", func(t *testing.T) {
		req := nagios.ServiceCountRequest{
			GeneralServiceRequest: nagios.GeneralServiceRequest{
				FormatOptions: nagios.FormatOptions{
					Whitespace: true,
					Enumerate:  true,
					Bitmask:    true,
					Duration:   true,
				},
				ServiceStatus: nagios.ServiceStatus{
					Ok:       true,
					Warning:  true,
					Critical: true,
					Unknown:  true,
					Pending:  true,
				},
				HostStatus: nagios.HostStatus{
					Up:          true,
					Down:        true,
					Unreachable: true,
					Pending:     true,
				},
			},
		}

		var count nagios.ServiceCount

		if err := c.Query(req, &count); err != nil {
			t.Errorf("Query: %v", err)
		}

		if count.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(count)
		}
	})

	t.Run("blank service list", func(t *testing.T) {
		req := nagios.ServiceListRequest{
			GeneralServiceRequest: nagios.GeneralServiceRequest{
				FormatOptions: nagios.FormatOptions{
					Enumerate: true,
				},
			},
		}

		var list nagios.ServiceList

		if err := c.Query(req, &list); err != nil {
			t.Errorf("Query: %v", err)
		}

		if list.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(list)
		}
	})

	t.Run("blank service list with options switched", func(t *testing.T) {
		req := nagios.ServiceListRequest{
			GeneralServiceRequest: nagios.GeneralServiceRequest{
				FormatOptions: nagios.FormatOptions{
					Whitespace: true,
					Enumerate:  true,
					Bitmask:    true,
					Duration:   true,
				},
				ServiceStatus: nagios.ServiceStatus{
					Ok:       true,
					Warning:  true,
					Critical: true,
					Unknown:  true,
					Pending:  true,
				},
				HostStatus: nagios.HostStatus{
					Up:          true,
					Down:        true,
					Unreachable: true,
					Pending:     true,
				},
				ShowDetails: true,
			},
		}

		var list nagios.ServiceList

		if err := c.Query(req, &list); err != nil {
			t.Errorf("Query: %v", err)
		}

		if list.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(list)
		}
	})

	t.Run("blank host", func(t *testing.T) {
		req := nagios.HostRequest{
			FormatOptions: nagios.FormatOptions{
				Enumerate: true,
			},
			HostName: "localhost",
		}

		var host nagios.Host

		if err := c.Query(req, &host); err != nil {
			t.Errorf("Query: %v", err)
		}

		if host.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(host)
		}
	})

	t.Run("blank host with options switched", func(t *testing.T) {
		req := nagios.HostRequest{
			FormatOptions: nagios.FormatOptions{
				Whitespace: true,
				Enumerate:  true,
				Bitmask:    true,
				Duration:   true,
			},
			HostName: "localhost",
		}

		var host nagios.Host

		if err := c.Query(req, &host); err != nil {
			t.Errorf("Query: %v", err)
		}

		if host.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(host)
		}
	})

	t.Run("blank service", func(t *testing.T) {
		req := nagios.ServiceRequest{
			FormatOptions: nagios.FormatOptions{
				Enumerate: true,
			},
			HostName:           "localhost",
			ServiceDescription: "HTTP",
		}

		var service nagios.Service

		if err := c.Query(req, &service); err != nil {
			t.Errorf("Query: %v", err)
		}

		if service.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(service)
		}
	})

	t.Run("blank service with options switched", func(t *testing.T) {
		req := nagios.ServiceRequest{
			FormatOptions: nagios.FormatOptions{
				Whitespace: true,
				Enumerate:  true,
				Bitmask:    true,
				Duration:   true,
			},
			HostName:           "localhost",
			ServiceDescription: "HTTP",
		}

		var service nagios.Service

		if err := c.Query(req, &service); err != nil {
			t.Errorf("Query: %v", err)
		}

		if service.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(service)
		}
	})

	t.Run("blank performance data", func(t *testing.T) {
		req := nagios.PerformanceDataRequest{
			FormatOptions: nagios.FormatOptions{
				Enumerate: true,
			},
		}

		var performance nagios.Performance

		if err := c.Query(req, &performance); err != nil {
			t.Errorf("Query: %v", err)
		}

		if performance.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(performance)
		}
	})

	t.Run("blank performance data with options switched", func(t *testing.T) {
		req := nagios.PerformanceDataRequest{
			FormatOptions: nagios.FormatOptions{
				Whitespace: true,
				Enumerate:  true,
				Bitmask:    true,
				Duration:   true,
			},
		}

		var performance nagios.Performance

		if err := c.Query(req, &performance); err != nil {
			t.Errorf("Query: %v", err)
		}

		if performance.Result.TypeText != success {
			t.Errorf("TypeText != %s", success)
		}

		if dumpResponse {
			spew.Dump(performance)
		}
	})
}
