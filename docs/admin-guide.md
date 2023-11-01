## Admin Guide

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Webhooks](#web-hooks)
- [Slash Commands](#slash-commands)
- [Onboard Users](#onboard-users)
- [FAQ](#faq)
- [Get Help](#get-help)

## Prerequsites

## Installation

### Marketplace Installation

1. Go to **Main Menu > Plugin Marketplace** in Mattermost.
2. Search for "Nagios" or find the plugin from the list.
3. Select **Install**.
4. When the plugin has downloaded and been installed, select **Configure**.

### Manual Installation

If your server doesn't have access to the internet, you can download the latest [plugin binary release](https://github.com/mattermost/mattermost-plugin-nagios/releases) and upload it to your server via **System Console > Plugin Management**. The releases on this page are the same used by the Marketplace. To learn more about how to upload a plugin, see [the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

### Cloud

### Upgrade

Ensure you configure the plugin before you enable it.

## Configuration

### Enter the URL for your Nagios instance

1. In Mattermost, go to **System Console > Plugins > Nagios**.
2. Set the **Nagios URL**. Remember to add `http://` or `https://` at the beginning (e.g.: `https://nagios.fedoraproject.org`).
3. Set **Nagios Username** to the Nagios user used to authentiacte.
4. Sett **Nagios Password** to the password of that user.
5. Select **Save**.
6. Go to **Plugins > Nagios**, then select **Enable**.
7. The plugin is now ready to use.

### Configure the configuration files watcher

*This step is optional, although highly recommended.*

Regenerate the token for the configuration files watcher:
1. In Mattermost, go to **System Console > Plugins > Nagios**.
2. Select **Regenerate** to regenerate the token.
3. Copy the token - you're going to use it later.
4. Select **Save**.
5. Switch to the machine where Nagios is running:
  - Download the latest stable version of the watcher from the [releases page](https://github.com/mattermost/mattermost-plugin-nagios/releases)
  - Move the watcher: `chmod +x watcher1.1.0.linux-amd64 && sudo mv watcher1.1.0.linux-amd64 /usr/local/bin/watcher`.

### Run the watcher as a systemd service

#### Prepare the systemd service unit file

Prepare the systemd service unit file by adjusting `dir` (default if not set: `/usr/local/nagios/etc/`), `url`, and `token` flags to your setup.

```sh
sudo bash -c 'cat << EOF > /etc/systemd/system/mattermost-plugin-nagios-watcher.service
[Unit]
Description=Nagios configuration files monitoring service
After=network.target

[Service]
Restart=on-failure
ExecStart=/usr/local/bin/watcher -dir /nagios/configuration/files/directory -url https://mattermost.server.address/plugins/nagios -token TheTokenFromStep1

[Install]
WantedBy=multi-user.target
EOF'
```

#### Start the watcher

```sh
systemctl daemon-reload
systemctl enable mattermost-plugin-nagios-watcher.service
systemctl start  mattermost-plugin-nagios-watcher.service
```

### Get logs
![logs](/docs/images/logs.png)

### Receive system monitoring reports
![reports](/docs/images/reports.png)

### Receive notifications about changes to the configuration
![changes](/docs/images/changes.png)

Ultimately, this will make you or your team more productive and make the experience with Nagios smoother.


## Web Hooks

## Slash Commands 

- `nagios`
    - `get-logs <alerts|notifications>`
        - `[host|service <host name|service description>]`
    - `set-logs-limit <count>`
    - `set-logs-start-time <seconds>`
    - `subscribe <report|configuration-changes>`
    - `unsubscribe <report|configuration-changes>`
    - `set-report-frequency <minutes>`

#### nagios

`nagios`

This is the root command.

##### get-logs

`get-logs <alerts|notifications> [<host|service> <host name|service description>]`

This action allows you to get alerts or notifications.

Example: `/nagios get-logs alerts`

###### host

`get-logs <alerts|notifications> host <host name>`

This optional parameter allows you to get alerts or notifications from a specific host.

Example: `/nagios get-logs alerts host bvmhost-p09-02.iad2.fedoraproject.org`

###### service

`get-logs <alerts|notifications> service <service description>`

This optional parameter allows you to get alerts or notifications from a specific service.

Example: `/nagios get-logs alerts service Swap-Is-Low`

##### set-logs-limit

`set-logs-limit <count>`

This action allows you to limit the number of logs `get-logs` fetches.

Example: `/nagios set-logs-limit 10`

##### set-logs-start-time

`set-logs-start-time <seconds>`

This action allows you to specify the age of the oldest log `get-logs` fetches.

Example: `/nagios set-logs-start-time 3600`

##### subscribe

`subscribe <report|configuration-changes>`

This action allows you to subscribe to system monitoring reports or configuration changes on the current channel.

Example: `/nagios subscribe report`

##### unsubscribe

`unsubscribe <report|configuration-changes>`

This action allows you to unsubscribe from system monitoring reports or configuration changes on the current channel.

Example: `/nagios unsubscribe configuration-changes`

##### set-report-frequency

`set-report-frequency <minutes>`

This action allows you to set the frequency of system monitoring reports.

Example: `/nagios set-report-frequency 60`

## Onboard Users

When you’ve tested the plugin and confirmed it’s working, notify your team so they can connect their Nagios account to Mattermost and get started. Copy and paste the text below, edit it to suit your requirements, and send it out.

> Hi team,
> 
> We've set up the Mattermost Nagios plugin, so you can get notifications from Nagios in Mattermost. 
> To get started, take a look at the slash commands section ([link](https://mattermost.gitbook.io/nagios-plugin/user-guide/slash-commands).

## FAQ

## Get Help