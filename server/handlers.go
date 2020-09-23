package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/ulumuri/go-nagios/nagios"
)

type commandHandlerFunc func(client *nagios.Client, parameters []string) string

var commandHandlers = map[string]commandHandlerFunc{
	"get-logs": getLogs,
}

func getLogsSpecific(parameters []string) (
	host, service bool,
	hostName, serviceDescription string) {

	if len(parameters) > 0 {
		switch parameters[0] {
		case "host":
			if len(parameters) > 1 {
				return true, false, parameters[1], ""
			}
		case "service":
			if len(parameters) > 1 {
				return false, true, "", parameters[1]
			}
		}
	}

	return false, false, "", ""
}

const (
	resultTypeTextSuccess = "Success"
	errorMessage          = "Getting logs unsuccessful: %v"
	logsFrom              = 24 * time.Hour
	maxLogs               = 100
)

func formatAlerts(alerts nagios.AlertList) string {
	if alerts.Result.TypeText != resultTypeTextSuccess {
		return fmt.Sprintf(errorMessage, alerts.Result.Message)
	}

	if len(alerts.Data.AlertList) == 0 {
		return "No alerts."
	}

	var b strings.Builder

	for _, a := range alerts.Data.AlertList {
		line := fmt.Sprintf("[%s] [%s] [%s] [%s] [%s] [%s] %s\n",
			time.Unix(a.Timestamp/1000, 0).String(),
			a.ObjectType,
			a.HostName,
			a.Description,
			a.StateType,
			a.State,
			a.PluginOutput)

		b.WriteString(line)
	}

	return b.String()
}

func extractHostName(e nagios.NotificationListEntry) string {
	if len(e.HostName) == 0 {
		return e.Name
	}
	return e.HostName
}

func formatNotifications(notifications nagios.NotificationList) string {
	if notifications.Result.TypeText != resultTypeTextSuccess {
		return fmt.Sprintf(errorMessage, notifications.Result.Message)
	}

	if len(notifications.Data.NotificationList) == 0 {
		return "No notifications."
	}

	var b strings.Builder

	for _, n := range notifications.Data.NotificationList {
		line := fmt.Sprintf("[%s] [%s] [%s] [%s] [%s] [%s] [%s] %s\n",
			time.Unix(n.Timestamp/1000, 0).String(),
			n.ObjectType,
			extractHostName(n),
			n.Description,
			n.Contact,
			n.NotificationType,
			n.Method,
			n.Message)

		b.WriteString(line)
	}

	return b.String()
}

// get-log alerts        <host>    <URL>
// get-log alerts        <service> <SVC>
// get-log notifications <host>    <URL>
// get-log notifications <service> <SVC>

func getLogs(client *nagios.Client, parameters []string) string {
	if len(parameters) == 0 {
		return "You must supply at least one parameter (alerts|notifications)"
	}

	now := time.Now()
	then := now.Add(-logsFrom)

	host, service, hostName, serviceDescription := getLogsSpecific(parameters[1:])

	switch parameters[0] {
	case "alerts":
		q := nagios.AlertListRequest{
			GeneralAlertRequest: nagios.GeneralAlertRequest{
				FormatOptions: nagios.FormatOptions{
					Enumerate: true,
				},
				Count: maxLogs,
				ObjectTypes: nagios.ObjectTypes{
					Host:    host,
					Service: service,
				},
				HostName:           hostName,
				ServiceDescription: serviceDescription,
				StartTime:          then.Unix(),
				EndTime:            now.Unix(),
			},
		}
		var alerts nagios.AlertList
		if err := client.Query(q, &alerts); err != nil {
			return fmt.Sprintf(errorMessage, err)
		}
		return formatAlerts(alerts)
	case "notifications":
		q := nagios.NotificationListRequest{
			GeneralNotificationRequest: nagios.GeneralNotificationRequest{
				FormatOptions: nagios.FormatOptions{
					Enumerate: true,
				},
				Count: maxLogs,
				ObjectTypes: nagios.ObjectTypes{
					Host:    host,
					Service: service,
				},
				HostName:           hostName,
				ServiceDescription: serviceDescription,
				StartTime:          then.Unix(),
				EndTime:            now.Unix(),
			},
		}
		var notifications nagios.NotificationList
		if err := client.Query(q, &notifications); err != nil {
			return fmt.Sprintf(errorMessage, err)
		}
		return formatNotifications(notifications)
	default:
		return fmt.Sprintf("Unknown parameter %s", parameters[0])
	}
}
