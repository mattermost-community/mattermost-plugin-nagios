package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mattermost/mattermost-plugin-nagios/go-nagios/nagios"
)

const (
	gettingReportUnsuccessful = "Getting system monitoring report unsuccessful"
	maximumReportLength       = 50
)

func gettingReportUnsuccessfulMessage(reportPart, message string) string {
	return fmt.Sprintf("%s (%s): %s", gettingReportUnsuccessful, reportPart, message)
}

func reportPreamble(t time.Time) string {
	return fmt.Sprintf("#### %s System monitoring report (%s)\n\n", barChartEmoji, t.Format(time.UnixDate))
}

func formatHostCount(count nagios.HostCount) string {
	if count.Result.TypeText != resultTypeTextSuccess {
		return gettingReportUnsuccessfulMessage("host summary", count.Result.TypeText)
	}

	var b strings.Builder

	b.WriteString("##### HOST SUMMARY\n\n")
	b.WriteString(fmt.Sprintf("%s Up: **%d**", upEmoji, count.Data.Count.Up))
	b.WriteString(fmt.Sprintf("  %s Down: **%d**", smallRedTriangleDownEmoji, count.Data.Count.Down))
	b.WriteString(fmt.Sprintf("  %s Unreachable: **%d**", mailboxWithNoMailEmoji, count.Data.Count.Unreachable))
	b.WriteString(fmt.Sprintf("  %s Pending: **%d**", hourglassFlowingSandEmoji, count.Data.Count.Pending))

	return b.String()
}

type extractedObject struct {
	name, state string
}

type extractedObjects []extractedObject

func (e extractedObjects) Len() int {
	return len(e)
}

func (e extractedObjects) Less(i, j int) bool {
	return e[i].name < e[j].name
}

func (e extractedObjects) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// extractHosts returns extractedObjects, extracted from hostListData. It
// returns unknownState state for every host it failed to extract the state.
func extractHosts(hostListData nagios.HostListData) extractedObjects {
	var hosts extractedObjects

	for k, v := range hostListData.HostList {
		var state string

		if err := json.Unmarshal(v, &state); err != nil {
			state = unknownState
		}

		hosts = append(hosts, extractedObject{name: k, state: state})
	}

	return hosts
}

func formatHostList(list nagios.HostList) string {
	if list.Result.TypeText != resultTypeTextSuccess {
		return gettingReportUnsuccessfulMessage("host list", list.Result.TypeText)
	}

	hosts := extractHosts(list.Data)

	sort.Sort(hosts)

	var b strings.Builder

	b.WriteString("##### HOST LIST\n\n")

	var abnormalOnly bool

	if len(hosts) > maximumReportLength {
		abnormalOnly = true

		b.WriteString("**Too many hosts. Showing only abnormal state hosts.**\n\n")
	}

	var linesWritten int

	for _, h := range hosts {
		if linesWritten == maximumReportLength {
			b.WriteString("\n\n**Skipped the rest of the hosts.**")
			return b.String()
		}

		if abnormalOnly && h.state == upState {
			continue
		}

		if linesWritten > 0 {
			b.WriteRune('\n')
		}

		b.WriteString(fmt.Sprintf("%s `%s` %s", emoji(h.state), h.name, strings.ToUpper(h.state)))
		linesWritten++
	}

	if linesWritten == 0 {
		b.WriteString("No hosts to show.")
	}

	return b.String()
}

func formatServiceCount(count nagios.ServiceCount) string {
	if count.Result.TypeText != resultTypeTextSuccess {
		return gettingReportUnsuccessfulMessage("service summary", count.Result.TypeText)
	}

	var b strings.Builder

	b.WriteString("##### SERVICE SUMMARY\n\n")
	b.WriteString(fmt.Sprintf("%s OK: **%d**", whiteCheckMarkEmoji, count.Data.Count.Ok))
	b.WriteString(fmt.Sprintf("  %s Warning: **%d**", warningEmoji, count.Data.Count.Warning))
	b.WriteString(fmt.Sprintf("  %s Critical: **%d**", bangBangEmoji, count.Data.Count.Critical))
	b.WriteString(fmt.Sprintf("  %s Unknown: **%d**", questionEmoji, count.Data.Count.Unknown))
	b.WriteString(fmt.Sprintf("  %s Pending: **%d**", hourglassFlowingSandEmoji, count.Data.Count.Pending))

	return b.String()
}

// extractServices returns extractedObjects, extracted from rawMessage. It will
// return a single-element slice initialized to unknownState state if it fails
// to process rawMessage.
func extractServices(rawMessage json.RawMessage) extractedObjects {
	var rawStates map[string]json.RawMessage

	if err := json.Unmarshal(rawMessage, &rawStates); err != nil {
		return extractedObjects{{state: unknownState}}
	}

	var services extractedObjects

	for k, v := range rawStates {
		var state string

		if err := json.Unmarshal(v, &state); err != nil {
			state = unknownState
		}

		services = append(services, extractedObject{name: k, state: state})
	}

	return services
}

func formatServiceList(list nagios.ServiceList) string {
	if list.Result.TypeText != resultTypeTextSuccess {
		return gettingReportUnsuccessfulMessage("service list", list.Result.TypeText)
	}

	var (
		reportLength int
		hosts        []string
	)

	hostToServices := make(map[string]extractedObjects)

	for k, v := range list.Data.ServiceList {
		services := extractServices(v)

		reportLength += len(services) + 1 // add 1 for a line with hostname.

		hosts, hostToServices[k] = append(hosts, k), services
	}

	sort.Strings(hosts)

	var b strings.Builder

	b.WriteString("##### SERVICE LIST\n\n")

	var abnormalOnly bool

	if reportLength > maximumReportLength {
		abnormalOnly = true

		b.WriteString("**Too many services. Showing only abnormal state services.**\n\n")
	}

	var linesWritten int

	for _, h := range hosts {
		var hostWritten bool

		// Sorting services here guarantees we sort only as much as we need to.
		sort.Sort(hostToServices[h])

		for _, s := range hostToServices[h] {
			if linesWritten == maximumReportLength {
				b.WriteString("\n\n**Skipped the rest of the services.**")
				return b.String()
			}

			if abnormalOnly && s.state == okState {
				continue
			}

			if !hostWritten {
				hostWritten = true

				if linesWritten > 0 {
					b.WriteString("\n\n")
				}

				b.WriteString(fmt.Sprintf("`%s`:", h))
				linesWritten++

				if linesWritten == maximumReportLength {
					continue
				}

				// This additional newline produces better Markdown, but we want
				// to write it only if theEnd hasn't been written.
				b.WriteRune('\n')
			}

			b.WriteString(fmt.Sprintf("\n- %s `%s` %s", emoji(s.state), s.name, strings.ToUpper(s.state)))
			linesWritten++
		}
	}

	if linesWritten == 0 {
		b.WriteString("No services to show.")
	}

	return b.String()
}
