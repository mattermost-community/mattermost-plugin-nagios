package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/ulumuri/go-nagios/nagios"
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

type extractedHost struct {
	name, state string
}

type extractedHosts []extractedHost

func (e extractedHosts) Len() int {
	return len(e)
}

func (e extractedHosts) Less(i, j int) bool {
	return e[i].name < e[j].name
}

func (e extractedHosts) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// extractHosts returns extractedHosts, extracted from hostListData. It returns
// unknownState state for every host it failed to extract the state.
func extractHosts(hostListData nagios.HostListData) extractedHosts {
	var hosts extractedHosts

	for k, v := range hostListData.HostList {
		var state string

		if err := json.Unmarshal(v, &state); err != nil {
			state = unknownState
		}

		hosts = append(hosts, extractedHost{name: k, state: state})
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

type extractedService struct {
	name, state string
}

type extractedServices []extractedService

func (e extractedServices) Len() int {
	return len(e)
}

func (e extractedServices) Less(i, j int) bool {
	return e[i].name < e[j].name
}

func (e extractedServices) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// extractServices returns extractedServices, extracted from rawMessage. It will
// return a single-element slice initialized to unknownState state if it fails
// to process rawMessage.
func extractServices(rawMessage json.RawMessage) extractedServices {
	var rawStates map[string]json.RawMessage

	if err := json.Unmarshal(rawMessage, &rawStates); err != nil {
		return extractedServices{{state: unknownState}}
	}

	var services extractedServices

	for k, v := range rawStates {
		var state string

		if err := json.Unmarshal(v, &state); err != nil {
			state = unknownState
		}

		services = append(services, extractedService{name: k, state: state})
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

	hostToServices := make(map[string]extractedServices)

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

func (p *Plugin) sendMessages(channelID string, messages ...string) error {
	var firstID string

	for i, m := range messages {
		post := &model.Post{
			UserId:    p.botUserID,
			ChannelId: channelID,
			RootId:    firstID,
			Message:   m,
		}

		created, err := p.API.CreatePost(post)
		if err != nil {
			return fmt.Errorf("p.API.CreatePost: %w", err)
		}

		if i == 0 {
			firstID = created.Id
		}
	}

	return nil
}

func (p *Plugin) sendMonitoringReport(channelID string) error {
	hostCountReq := nagios.HostCountRequest{
		GeneralHostRequest: nagios.GeneralHostRequest{
			FormatOptions: nagios.FormatOptions{
				Enumerate: true,
			},
		}}

	var hostCount nagios.HostCount

	if err := p.client.Query(hostCountReq, &hostCount); err != nil {
		return fmt.Errorf("client.Query: %w", err)
	}

	hostListReq := nagios.HostListRequest{
		GeneralHostRequest: nagios.GeneralHostRequest{
			FormatOptions: nagios.FormatOptions{
				Enumerate: true,
			},
		}}

	var hostList nagios.HostList

	if err := p.client.Query(hostListReq, &hostList); err != nil {
		return fmt.Errorf("client.Query: %w", err)
	}

	serviceCountReq := nagios.ServiceCountRequest{
		GeneralServiceRequest: nagios.GeneralServiceRequest{
			FormatOptions: nagios.FormatOptions{
				Enumerate: true,
			},
		}}

	var serviceCount nagios.ServiceCount

	if err := p.client.Query(serviceCountReq, &serviceCount); err != nil {
		return fmt.Errorf("client.Query: %w", err)
	}

	serviceListReq := nagios.ServiceListRequest{
		GeneralServiceRequest: nagios.GeneralServiceRequest{
			FormatOptions: nagios.FormatOptions{
				Enumerate: true,
			},
		}}

	var serviceList nagios.ServiceList

	if err := p.client.Query(serviceListReq, &serviceList); err != nil {
		return fmt.Errorf("client.Query: %w", err)
	}

	if err := p.sendMessages(
		channelID,
		reportPreamble(time.Now()),
		formatHostCount(hostCount),
		formatHostList(hostList),
		formatServiceCount(serviceCount),
		formatServiceList(serviceList)); err != nil {
		return fmt.Errorf("p.sendMessages: %w", err)
	}

	return nil
}

const reportLockKey = "report-lock"

var reportLock = []byte{0b10, 0b1, 0b11, 0b111}

func (p *Plugin) monitoringReportLoop() {
	for {
		time.Sleep(time.Minute)

		d, err := getReportFrequency(p.API)
		if err != nil {
			p.API.LogError("getReportFrequency", logErrorKey, err)
			continue
		}

		opts := model.PluginKVSetOptions{
			Atomic:          true,
			OldValue:        nil,
			ExpireInSeconds: int64(d.Seconds()),
		}

		ok, appErr := p.API.KVSetWithOptions(reportLockKey, reportLock, opts)
		if !ok {
			if appErr != nil {
				p.API.LogError("KVSetWithOptions", logErrorKey, err)
			}

			continue
		}

		p.API.LogDebug("monitoringReportLoop: acquired lock", "id", p.API.GetDiagnosticId())

		c, err := getReportChannel(p.API)
		if err != nil {
			p.API.LogError("getReportChannel", logErrorKey, err)
			continue
		}

		if c == "" { // fast path, there is no subscription.
			p.API.LogDebug("monitoringReportLoop: no subscription")
			continue
		}

		if err := p.sendMonitoringReport(c); err != nil {
			p.API.LogError("sendMonitoringReport", logErrorKey, err)
		}
	}
}
