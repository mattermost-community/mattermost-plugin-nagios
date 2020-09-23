package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/ulumuri/go-nagios/nagios"
)

type commandHandlerFunc func(client *nagios.Client, parameters []string) string

var commandHandlers = map[string]commandHandlerFunc{
	"get-logs":             getLogs,
	"set-logs-limit":       nil,
	"set-logs-start-time":  nil,
	"help":                 nil,
	"set-report-frequency": nil,
}

const (
	resultTypeTextSuccess = "Success"
	defaultLogsFrom       = 24 * time.Hour
	defaultMaxLogs        = 100
)

func errorMessage(message interface{}) string {
	return fmt.Sprintf("Getting logs unsuccessful: %v.", message)
}

func unknownParameterMessage(parameter string) string {
	return fmt.Sprintf("Unknown parameter (%s).", parameter)
}

func formatNagiosTimestamp(t int64) string {
	return time.Unix(t/1000, 0).String()
}

func extractHostName(e nagios.NotificationListEntry) string {
	if len(e.HostName) == 0 {
		return e.Name
	}
	return e.HostName
}

// TODO(amwolff, DanielSz50): rewrite formatAlertListEntry (mimic showlog.cgi).
func formatAlertListEntry(e nagios.AlertListEntry) string {
	return fmt.Sprintf("[%s] [%s] [%s] [%s] [%s] [%s] %s",
		formatNagiosTimestamp(e.Timestamp),
		e.ObjectType,
		e.HostName,
		e.Description,
		e.StateType,
		e.State,
		e.PluginOutput)
}

func formatAlerts(alerts nagios.AlertList) string {
	if alerts.Result.TypeText != resultTypeTextSuccess {
		return errorMessage(alerts.Result.Message)
	}

	if len(alerts.Data.AlertList) == 0 {
		return "No alerts."
	}

	var b strings.Builder

	for i, v := range alerts.Data.AlertList {
		if i > 0 {
			b.WriteRune('\n')
		}
		b.WriteString(formatAlertListEntry(v))
	}

	return b.String()
}

// TODO(amwolff, DanielSz50): rewrite formatNotificationListEntry (mimic showlog.cgi).
func formatNotificationListEntry(e nagios.NotificationListEntry) string {
	return fmt.Sprintf("[%s] [%s] [%s] [%s] [%s] [%s] [%s] %s",
		formatNagiosTimestamp(e.Timestamp),
		e.ObjectType,
		extractHostName(e),
		e.Description,
		e.Contact,
		e.NotificationType,
		e.Method,
		e.Message)
}

func formatNotifications(notifications nagios.NotificationList) string {
	if notifications.Result.TypeText != resultTypeTextSuccess {
		return errorMessage(notifications.Result.Message)
	}

	if len(notifications.Data.NotificationList) == 0 {
		return "No notifications."
	}

	var b strings.Builder

	for i, v := range notifications.Data.NotificationList {
		if i > 0 {
			b.WriteRune('\n')
		}
		b.WriteString(formatNotificationListEntry(v))
	}

	return b.String()
}

// get-log alerts        <host>    <URL>
// get-log alerts        <service> <SVC>
// get-log notifications <host>    <URL>
// get-log notifications <service> <SVC>

func getLogsSpecific(parameters []string) (hostName, serviceDescription, message string, ok bool) {
	if len(parameters) == 0 {
		return "", "", "", true
	}

	switch parameters[0] {
	case "host":
		if len(parameters) < 2 {
			return "", "", "You must supply host name.", false
		}
		return parameters[1], "", "", true
	case "service":
		if len(parameters) < 2 {
			return "", "", "You must supply service description.", false
		}
		return "", parameters[1], "", true
	default:
		return "", "", unknownParameterMessage(parameters[0]), false
	}
}

func getLogs(client *nagios.Client, parameters []string) string {
	if len(parameters) == 0 {
		return "You must supply at least one parameter (alerts|notifications)."
	}

	hostName, serviceDescription, message, ok := getLogsSpecific(parameters[1:])
	if !ok {
		return message
	}

	now := time.Now()
	then := now.Add(-defaultLogsFrom)

	switch parameters[0] {
	case "alerts":
		q := nagios.AlertListRequest{
			GeneralAlertRequest: nagios.GeneralAlertRequest{
				FormatOptions: nagios.FormatOptions{
					Enumerate: true,
				},
				Count:              defaultMaxLogs,
				HostName:           hostName,
				ServiceDescription: serviceDescription,
				StartTime:          then.Unix(),
				EndTime:            now.Unix(),
			},
		}
		var alerts nagios.AlertList
		if err := client.Query(q, &alerts); err != nil {
			return errorMessage(err)
		}
		return formatAlerts(alerts)
	case "notifications":
		q := nagios.NotificationListRequest{
			GeneralNotificationRequest: nagios.GeneralNotificationRequest{
				FormatOptions: nagios.FormatOptions{
					Enumerate: true,
				},
				Count:              defaultMaxLogs,
				HostName:           hostName,
				ServiceDescription: serviceDescription,
				StartTime:          then.Unix(),
				EndTime:            now.Unix(),
			},
		}
		var notifications nagios.NotificationList
		if err := client.Query(q, &notifications); err != nil {
			return errorMessage(err)
		}
		return formatNotifications(notifications)
	default:
		return unknownParameterMessage(parameters[0])
	}
}
