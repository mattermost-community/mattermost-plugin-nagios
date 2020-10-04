package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/ulumuri/go-nagios/nagios"
)

const (
	gettingReportUnsuccessful = "Getting monitoring report unsuccessful"
	maximumReportLength       = 50
)

func gettingReportUnsuccessfulMessage(reportPart, message string) string {
	return fmt.Sprintf("%s (%s): %s", gettingReportUnsuccessful, reportPart, message)
}

func reportPreamble(t time.Time) string {
	return fmt.Sprintf("#### :bar_chart: System monitoring report (%s)\n\n", t.Format(time.UnixDate))
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

func formatHostList(list nagios.HostList) string {
	if list.Result.TypeText != resultTypeTextSuccess {
		return gettingReportUnsuccessfulMessage("host list", list.Result.TypeText)
	}

	var b strings.Builder

	b.WriteString("##### HOST LIST\n\n")

	var abnormalOnly bool
	if len(list.Data.HostList) > maximumReportLength {
		abnormalOnly = true
		b.WriteString("**Too many hosts. Showing only abnormal state hosts.**\n\n")
	}

	var linesWritten int
	for k, v := range list.Data.HostList {
		if linesWritten == maximumReportLength {
			b.WriteString("\n\n**Skipped the rest of the hosts.**")
			break
		}

		var state string

		if err := json.Unmarshal(v, &state); err != nil {
			state = unknownState
		}

		if state == upState && abnormalOnly {
			continue
		}

		if linesWritten > 0 {
			b.WriteRune('\n')
		}
		b.WriteString(fmt.Sprintf("%s `%s` %s", emoji(state), k, strings.ToUpper(state)))
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

// extractServices returns a slice of extractedService, extracted from
// rawMessage. It returns a slice with single element initialized to a
// unknownState state if it fails to process rawMessage.
func extractServices(rawMessage json.RawMessage) []extractedService {
	var rawStates map[string]json.RawMessage

	if err := json.Unmarshal(rawMessage, &rawStates); err != nil {
		return []extractedService{{state: unknownState}}
	}

	var services []extractedService

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

	var reportLength int

	hostToServices := make(map[string][]extractedService)

	for k, v := range list.Data.ServiceList {
		services := extractServices(v)

		reportLength += len(services) + 1 // add 1 for a line with hostname.

		hostToServices[k] = services
	}

	var b strings.Builder

	b.WriteString("##### SERVICE LIST\n\n")

	var abnormalOnly bool

	if reportLength > maximumReportLength {
		abnormalOnly = true
		b.WriteString("**Too many services. Showing only abnormal state services.**\n\n")
	}

	var linesWritten int

	const theEnd = "\n\n**Skipped the rest of the services.**"

	for host, services := range hostToServices {
		for i, s := range services {
			if s.state == okState && abnormalOnly {
				continue
			}

			if i == 0 {
				if linesWritten > 0 {
					b.WriteRune('\n')
				}
				b.WriteString(fmt.Sprintf("`%s`:", host))
				linesWritten++

				if linesWritten == maximumReportLength {
					b.WriteString(theEnd)
					goto end
				}
			}

			b.WriteString(fmt.Sprintf("\n\t%s `%s` %s", emoji(s.state), s.name, strings.ToUpper(s.state)))
			linesWritten++

			if linesWritten == maximumReportLength {
				b.WriteString(theEnd)
				goto end
			}
		}
	}

	if linesWritten == 0 {
		b.WriteString("No services to show.")
	}

end: // Dijkstra probably hates me.
	return b.String()
}

func (p *Plugin) sendMessages(channelID string, messages ...string) error {
	for _, m := range messages {
		post := &model.Post{
			UserId:    p.botUserID,
			ChannelId: channelID,
			Message:   m,
		}
		if _, err := p.API.CreatePost(post); err != nil {
			return fmt.Errorf("p.API.CreatePost: %w", err)
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
		return fmt.Errorf("clent.Query: %w", err)
	}

	hostListReq := nagios.HostListRequest{
		GeneralHostRequest: nagios.GeneralHostRequest{
			FormatOptions: nagios.FormatOptions{
				Enumerate: true,
			},
		}}

	var hostList nagios.HostList

	if err := p.client.Query(hostListReq, &hostList); err != nil {
		return fmt.Errorf("clent.Query: %w", err)
	}

	serviceCountReq := nagios.ServiceCountRequest{
		GeneralServiceRequest: nagios.GeneralServiceRequest{
			FormatOptions: nagios.FormatOptions{
				Enumerate: true,
			},
		}}

	var serviceCount nagios.ServiceCount

	if err := p.client.Query(serviceCountReq, &serviceCount); err != nil {
		return fmt.Errorf("clent.Query: %w", err)
	}

	serviceListReq := nagios.ServiceListRequest{
		GeneralServiceRequest: nagios.GeneralServiceRequest{
			FormatOptions: nagios.FormatOptions{
				Enumerate: true,
			},
		}}

	var serviceList nagios.ServiceList

	if err := p.client.Query(serviceListReq, &serviceList); err != nil {
		return fmt.Errorf("clent.Query: %w", err)
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

func (p *Plugin) addMonitoringReport(channelID string, stop <-chan bool) {
	for {
		d, err := getReportFrequency(p.API)
		if err != nil {
			p.API.LogError("getReportFrequency", logErrorKey, err)
		}
		select {
		case <-stop:
			return
		case <-time.NewTimer(d).C:
			if err := p.sendMonitoringReport(channelID); err != nil {
				p.API.LogError("sendMonitoringReport", logErrorKey, err)
			}
		}
	}
}

func (p *Plugin) subscribe(channelID string, parameters []string) string {
	if len(parameters) > 0 {
		return "subscribe does not take any additional parameters."
	}

	stop := make(chan bool, 1)

	go p.addMonitoringReport(channelID, stop)

	p.subscriptionStop = stop

	return "Subscribed successfully."
}

func subscribe(p *Plugin, channelID string, parameters []string) string {
	return p.subscribe(channelID, parameters)
}

func (p *Plugin) unsubscribe(parameters []string) string {
	if len(parameters) > 0 {
		return "unsubscribe does not take any additional parameters."
	}

	p.subscriptionStop <- true

	return "Unsubscribed successfully."
}

func unsubscribe(p *Plugin, channelID string, parameters []string) string {
	return p.unsubscribe(parameters)
}
