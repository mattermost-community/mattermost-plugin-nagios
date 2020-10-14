# Mattermost Nagios Plugin

<!-- TODO(amwolff): add CI badges and stuff. -->

**Maintainers**: [@amwolff](https://github.com/amwolff) & [@DanielSz50](https://github.com/DanielSz50)

A Nagios plugin for Mattermost. Supports Nagios Core >= 4.4.x.

## Table of contents
- [About](https://github.com/ulumuri/mattermost-plugin-nagios/#about)
    - [Audience](https://github.com/ulumuri/mattermost-plugin-nagios/#audience)
    - [Important notice](https://github.com/ulumuri/mattermost-plugin-nagios/#important-notice)
- [Installing the plugin](https://github.com/ulumuri/mattermost-plugin-nagios/#installing-the-plugin)
- [Configuring the plugin](https://github.com/ulumuri/mattermost-plugin-nagios/#configuring-the-plugin)
    - [Configuring the configuration files watcher](https://github.com/ulumuri/mattermost-plugin-nagios/#configuring-the-configuration-files-watcher)
- [Updating the plugin](https://github.com/ulumuri/mattermost-plugin-nagios/#updating-the-plugin)
- [Using the plugin](https://github.com/ulumuri/mattermost-plugin-nagios/#using-the-plugin)
    - [Slash commands overview](https://github.com/ulumuri/mattermost-plugin-nagios/#slash-commands-overview)
    - [Slash commands documentation](https://github.com/ulumuri/mattermost-plugin-nagios/#slash-commands-documentation)
        - [nagios](https://github.com/ulumuri/mattermost-plugin-nagios/#nagios)
            - [get-logs](https://github.com/ulumuri/mattermost-plugin-nagios/#get-logs)
                - [host](https://github.com/ulumuri/mattermost-plugin-nagios/#host)
                - [service](https://github.com/ulumuri/mattermost-plugin-nagios/#service)
        - [set-logs-limit](https://github.com/ulumuri/mattermost-plugin-nagios/#set-logs-limit)
        - [set-logs-start-time](https://github.com/ulumuri/mattermost-plugin-nagios/#set-logs-start-time)
        - [subscribe](https://github.com/ulumuri/mattermost-plugin-nagios/#subscribe)
        - [unsubscribe](https://github.com/ulumuri/mattermost-plugin-nagios/#unsubscribe)
        - [set-report-frequency](https://github.com/ulumuri/mattermost-plugin-nagios/#set-report-frequency)
- [Contributing](https://github.com/ulumuri/mattermost-plugin-nagios/#contributing)

## About

This plugin allows you to

- [x] get logs from specific systems without leaving the Mattermost
    - get alerts and notifications instantly delivered, resembling the `showlog.cgi` UI
- [x] receive the system monitoring reports on a subscribed channel
    - be frequently informed which hosts and/or services have an abnormal state
- [ ] (in progress) receive notifications about changes to the configuration on a subscribed channel
    - anytime a change has been made to Nagios configuration, receive a diff between the old and the new version

Ultimately, this will make you or your team more productive and make the experience with Nagios smoother.

### Audience

This guide is for Mattermost System Admins setting up the Nagios plugin and Mattermost users who want information about the plugin functionality.

### Important notice

If you are a Nagios admin/user and think there is something this plugin lacks or something that it does could be done the other way around, let us know!
We are trying to develop this plugin based on users' needs.
If there is a certain feature you or your team needs, open up an issue and explain your needs.
We will be happy to help.

## Installing the plugin

1. Download the latest version of the plugin from the [releases page](https://github.com/ulumuri/mattermost-plugin-nagios/releases)
2. In Mattermost, go to **System Console → Plugins → Management**
3. Upload the plugin in the **Upload Plugin** section
4. Configure the plugin before you enable it :arrow_down:

## Configuring the plugin

1. Enter the URL for your Nagios instance
    1. In Mattermost, go to **System Console → Plugins → Nagios**
    2. Set the **Nagios URL**
        1. Remember to add `http://` or `https://` at the beginning!
        2. Example: `https://nagios.fedoraproject.org`
2. Click *Save* to save the settings
3. In Mattermost, go to **System Console → Plugins → Management** and click *Enable* underneath the Nagios plugin
4. The plugin is now ready to use! :congratulations:

### Configuring the configuration files watcher

*This step is optional, although highly recommended.*

1. Regenerate the token for the configuration files watcher
    1. In Mattermost, go to **System Console → Plugins → Nagios**
    2. Click *Regenerate* to regenerate the token
    3. Copy the token (you are going to use it later)
2. Click *Save* to save the settings

`TODO(amwolff): add the rest of the setup instructions for the configuration files watcher.`

## Updating the plugin

To update the plugin repeat the [Installation](https://github.com/ulumuri/mattermost-plugin-nagios/#installation) step.

## Using the plugin

Interaction with the plugin involves using the slash commands.

### Slash commands overview

- `nagios`
    - `get-logs <alerts|notifications>`
        - `[host|service <host name|service description>]`
    - `set-logs-limit <count>`
    - `set-logs-start-time <seconds>`
    - `subscribe`
    - `unsubscribe`
    - `set-report-frequency <minutes>`

### Slash commands documentation

#### nagios

`nagios`

This is the root command.

##### get-logs

`get-logs <alerts|notifications> [<host|service> <host name|service description>]`

This action allows you to get alerts or notifications.

Example: `/nagios get-logs alerts`

###### host

`get-logs <alerts|notifications> host <host>`

This optional parameter allows you to get alerts or notifications from a specific host.

Example: `/nagios get-logs alerts host bvmhost-p09-02.iad2.fedoraproject.org`

###### service

`get-logs <alerts|notifications> service <service>`

This optional parameter allows you to get alerts or notifications from a specific service.

Example: `/nagios get-logs alerts service Swap-Is-Low`

##### set-logs-limit

`set-logs-limit <count>`

This action allows you to limit the number of logs fetched by `get-logs`.

Example: `/nagios set-logs-limit 100`

##### set-logs-start-time

`set-logs-start-time <seconds>`

This action allows you to specify the age of the oldest log `get-logs` fetches.

Example: `/nagios set-logs-start-time 3600`

##### subscribe

`subscribe`

This action allows you to subscribe to the system monitoring reports on the current channel.

Example: `/nagios subscribe`

##### unsubscribe

`unsubscribe`

This action allows you to unsubscribe from the system monitoring reports.

Example: `/nagios unsubscribe`

##### set-report-frequency

`set-report-frequency <minutes>`

This action allows you to set the frequency of the system monitoring reports.

Example: `/nagios set-report-frequency 15`

## Contributing

<!-- TODO(amwolff): Write more about contributing to the plugin. Add CONTRIBUTING.md? -->

This repository uses the [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template).
Therefore, developing this plugin is roughly the same as it is with every plugin using the template.
All the necessary steps to develop are in the template's repository.
