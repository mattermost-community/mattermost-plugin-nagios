## Admin Guide

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Webhooks](#web-hooks)
- [Slash Commands](#slash-commands)
- [Onboard Users](#onboard-users)
- [FAQ](#faq)
- [Get Help](#get-help)

## Installation

### Manual Installation

You can download the latest [plugin binary release](https://github.com/mattermost/mattermost-plugin-nagios/releases) and upload it to your server via **System Console > Plugin Management**.

## Configuration

### Enter the URL for your Nagios instance

1. In Mattermost, go to **System Console > Plugins > Nagios**.
2. Set the **Nagios URL**. Remember to add `http://` or `https://` at the beginning (e.g.: `https://nagios.fedoraproject.org`).
3. Set **Nagios Username** to the Nagios user used to authentiacte.
4. Sett **Nagios Password** to the password of that user.
5. Select **Save**.
6. Go to **Plugins > Nagios, then select \*\*Enable**.
7. The plugin is now ready to use.

### Configure the configuration files watcher

This step is optional, although highly recommended.

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
systemctl start mattermost-plugin-nagios-watcher.service
```
### Get logs

![logs](/docs/images/logs.png)

### Receive system monitoring reports

![reports](/docs/images/reports.png)

### Receive notifications about changes to the configuration

![changes](/docs/images/changes.png)

Ultimately, this will make you or your team more productive and make the experience with Nagios smoother.

## [Slash Commands](../README.md/#slash-commands-overview)

## Onboard Users

When you’ve tested the plugin and confirmed it’s working, notify your team so they can connect their Nagios account to Mattermost and get started. Copy and paste the text below, edit it to suit your requirements, and send it out.

> Hi team,
>
> We've set up the Mattermost Nagios plugin, so you can get notifications from Nagios in Mattermost.
> To get started, take a look at the slash commands section ([link](https://mattermost.gitbook.io/nagios-plugin/user-guide/slash-commands)).
