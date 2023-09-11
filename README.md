# Mattermost Nagios Plugin

**Maintainers**: [@amwolff](https://github.com/amwolff) & [@DanielSz50](https://github.com/DanielSz50)

A Mattermost Plugin for Nagios to get logs, alerts, and notifications in Mattermost. Supports Nagios Core >= 4.4.x.

This plugin allows you to:

* Get logs from specific systems without leaving Mattermost.
* Get alerts and notifications resembling the `showlog.cgi` UI instantly delivered.
* Receive system monitoring reports on a subscribed channel.
* Be frequently informed which hosts and services have an abnormal state.
* Receive notifications about changes to the configuration on a subscribed channel.
* Receive a diff between the old and the new version anytime a change has been made to Nagios configuration.

This README provides guidance on installation, configuration, and usage.

## Installation

### Marketplace Installation

1. Go to **Main Menu > Plugin Marketplace** in Mattermost.
2. Search for "Nagios" or find the plugin from the list.
3. Select **Install**.
4. When the plugin has downloaded and been installed, select **Configure**.

### Manual Installation

If your server doesn't have access to the internet, you can download the latest [plugin binary release](https://github.com/mattermost/mattermost-plugin-nagios/releases) and upload it to your server via **System Console > Plugin Management**. The releases on this page are the same used by the Marketplace. To learn more about how to upload a plugin, see [the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

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
![logs](https://github.com/mattermost/mattermost-plugin-nagios/assets/74422101/7701b754-dc79-4199-9779-a52c173fbbba)

### Receive system monitoring reports
![reports](https://github.com/mattermost/mattermost-plugin-nagios/assets/74422101/c4369e89-a518-498a-841e-a608dacce40c)

### Receive notifications about changes to the configuration
![changes](https://github.com/mattermost/mattermost-plugin-nagios/assets/74422101/d8920978-c184-45bd-b443-cdae31a979e4)

Ultimately, this will make you or your team more productive and make the experience with Nagios smoother.

### Onboard your users

When you’ve tested the plugin and confirmed it’s working, notify your team so they can connect their Nagios account to Mattermost and get started. Copy and paste the text below, edit it to suit your requirements, and send it out.

> Hi team,
> 
> We've set up the Mattermost Nagios plugin, so you can get notifications from Nagios in Mattermost. 
> To get started, take a look at the slash commands section ([link](https://mattermost.gitbook.io/nagios-plugin/user-guide/slash-commands).

## Development

This repository uses the [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template). Therefore, developing this plugin is roughly the same as it is with every plugin using the template. All the necessary steps to develop are in the template's repository.

If you are a Nagios admin/user and think there is something this plugin lacks or something that it does could be done another way, let us know! We are trying to develop this plugin based on users' needs. If there is a certain feature you or your team needs, open up an issue, and explain your needs. We will be happy to help.

This plugin only contains a server portion. Read our documentation about the [Developer Workflow](https://developers.mattermost.com/extend/plugins/developer-workflow/) and [Developer Setup](https://developers.mattermost.com/extend/plugins/developer-setup/) for more information about developing and extending plugins.

### Running a Nagios server with Docker

There is a [docker-compose.yml](https://github.com/mattermost/mattermost-plugin-naguis/blob/master/dev/docker-compose.yml) in the `dev` folder of the repository, configured to run a Nagios server for development. You can run `make nagios` in the root of the repository to spin up the Nagios server. The Nagios web application will be served at http://localhost:8080.

You can login with these credentials:

- Username: `nagiosadmin`
- Password: `nagios`

## Develop the watcher

To build the watcher, you can use the following command:

```sh
env GOOS=linux GOARCH=amd64 go build -o dist/watcherX.Y.Z.linux-amd64 -a -v cmd/watcher/main.go
```

Of course, you can build the watcher for other operating systems and architectures too.

### Deploy with Local Mode

If your Mattermost server is running locally, you can enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode) to streamline deploying your plugin. After configuring it, just run:

```sh
make deploy
```

### Deploy with credentials

Alternatively, you can authenticate with the server's API with a [personal access token](https://docs.mattermost.com/developer/personal-access-tokens.html):

```sh
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
```

## Usage

### Slash Commands Overview

- `nagios`
    - `get-logs <alerts|notifications>`
        - `[host|service <host name|service description>]`
    - `set-logs-limit <count>`
    - `set-logs-start-time <seconds>`
    - `subscribe <report|configuration-changes>`
    - `unsubscribe <report|configuration-changes>`
    - `set-report-frequency <minutes>`

### Slash commands

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

## Contribute

See the Development guide section for information about how to contribute to this plugin.

## I saw a bug, I have a feature request or a suggestion

Please file a [GitHub issue](https://github.com/mattermost/mattermost-plugin-nagios/issues), it will be very useful!

Pull Requests are welcome! You can contact us on the [Mattermost Community ~Plugin: Nagios channel](https://community.mattermost.com/core/channels/plugin-nagios).

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options.

## Share feedback

Feel free to create a [GitHub Issue](https://github.com/mattermost/mattermost-plugin-nagios/issues) or join the [Nagios Plugin channel](https://community.mattermost.com/core/channels/plugin-nagios) on the Mattermost Community server to discuss.

## Security vulnerability disclosure

Please report any security vulnerability to [https://mattermost.com/security-vulnerability-report/](https://mattermost.com/security-vulnerability-report/).
