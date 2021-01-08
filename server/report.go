package main

import (
	"fmt"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/ulumuri/go-nagios/nagios"
)

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

var reportLock = make([]byte, 0)

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
