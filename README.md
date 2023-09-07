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

![Screenshot of getting logs in action](../screenshots/logs.png)

### Receive system monitoring reports

![Screenshot of receiving system monitoring reports in action](../screenshots/reports.png)

### Receive notifications about changes to the configuration

![Screenshot of receiving notifications about changes to the configuration in action](../screenshots/changes.png)

Ultimately, this will make you or your team more productive and make the experience with Nagios smoother.

### Onboard Your Users

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

### Slash Commands Documentation

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

## Nagios Core API Client

This is a Go client for the Nagios Core API that can communicate with the Nagios Core JSON CGIs. It's the programmatic equivalent of the <nagios>/nagios/jsonquery.html.

### Support

**Archive JSON CGI**

    alertcount - Yes
    alertlist - Yes
    notificationcount - Yes
    notificationlist - Yes
    statechangelist - No
    availability - No

**Object JSON CGI**

    hostcount - No
    hostlist - No
    host - No
    hostgroupcount - No
    hostgrouplist - No
    hostgroup - No
    servicecount - No
    servicelist - No
    service - No
    servicegroupcount - No
    servicegrouplist - No
    servicegroup - No
    contactcount - No
    contactlist - No
    contact - No
    contactgroupcount - No
    contactgrouplist - No
    contactgroup - No
    timeperiodcount - No
    timeperiodlist - No
    timeperiod - No
    commandcount - No
    commandlist - No
    command - No
    servicedependencycount - No
    servicedependencylist - No
    serviceescalationcount - No
    serviceescalationlist - No
    hostdependencycount - No
    hostdependencylist - No
    hostescalationcount - No
    hostescalationlist - No

**Status JSON CGI**

    hostcount - Yes
    hostlist - Yes
    host - Yes
    servicecount - Yes
    servicelist - Yes
    service - Yes
    commentcount - No
    commentlist - No
    comment - No
    downtimecount - No
    downtimelist - No
    downtime - No
    programstatus - No
    performancedata - Yes

## Sync

The sync tool is a proof-of-concept implementation of a tool for synchronizing mattermost plugin repositories with the mattermost-plugin-starter-template repo.

### Overview

At its core the tool is just a collection of checks and actions that are executed according to a synchronization plan (see [./build/sync/plan.yml](https://github.com/mattermost/mattermost-plugin-starter-template/blob/sync/build/sync/plan.yml) for an example). The plan defines a set of files and/or directories that need to be kept in sync between the plugin repository and the template (this repo).

For each set of paths, a set of actions to be performed is outlined. No more than one action of that set will be executed - the first one whose checks pass. Other actions are meant to act as fallbacks. The idea is to be able to e.g. overwrite a file if it has no local changes or apply a format-specific merge algorithm otherwise.

Before running each action, the tool will check if any checks are defined for that action. If there are any, they will be executed and their results examined. If all checks pass, the action will be executed. If there is a check failure, the tool will locate the next applicable action according to the plan and start over with it.

The synchronization plan can also run checks before running any actions, e.g. to check if the source and target worktrees are clean.

### Run Sync Tool

The tool can be executed from the root of this repository with a command:

```
$ go run ./build/sync/main.go ./build/sync/plan.yml ../mattermost-plugin-github
```

(assuming `mattermost-plugin-github` is the target repository we want to synchronize with the source).

### plan.yml

The `plan.yml` file (located in `build/sync/plan.yml`) consists of two parts:

- checks
- actions

The `checks` section defines tests to run before executing the plan itself. Currently the only available such check is `repo_is_clean defined` as:

```
type: repo_is_clean
params:
  repo: source
```

The `repo` parameter takes one of two values:
- `source` - the `mattermost-plugin-starter-template` repository.
- `target` - the repository of the plugin being updated.

The `actions` section defines actions to be run as part of the synchronization. Each entry in this section has the form:

```
paths:
  - path1
  - path2
actions:
  - type: action_type
    params:
      action_parameter: value
    conditions:
      - type: check_type
        params:
          check_parameter: value
```

`paths` is a list of file or directory paths (relative to the root of the repository) synchronization should be performed on.

Each action in the `actions` section is defined by its type. Currently supported action types are:
- `overwrite_file` - overwrite the specified file in the target repository with the file in the source repository.
- `overwrite_directory` - overwrite a directory.

Both actions accept a parameter called create which determines if the file or directory should be created if it does not exist in the target repository.

The `conditions` part of an action definition defines tests that need to pass for the action to be run. Available checks are:
- `exists`
- `file_unaltered`

The `exists` check takes a single parameter - `repo` (referencing either the source or target repository) and it passes only if the file or directory the action is about to be run on exists. If the repo parameter is not specified, it will default to `target`.

The `file_unaltered` check is only applicable to file paths. It passes if the file has not been altered - i.e. it is identical to some version of that same file in the reference repository (usually `source`). This check takes two parameters:
- `in` - repository to check the file in, default `target`.
- `compared-to` - repository to check the file against, default `source`.

When multiple actions are specified for a set of paths, the `sync` tool will only execute a single action for each path. The first action in the list, whose conditions are all satisfied will be executed.
If an action fails due to an error, the synchronization run will be aborted.

### Caveat emptor

This is a very basic proof-of-concept and there are many things that should be improved/implemented: (in no specific order)

1. Format-specific merge actions for go.mod, go.sum, webapp/package.json and other files should be implemented.
2. Better logging should be implemented.
3. Handling action dependencies should be investigated. e.g. if the `build` directory is overwritten, that will in some cases mean that the `go.mod` file also needs to be updated.
4. Storing the tree-hash of the template repository that the plugin was synchronized with would allow improving the performance of the tool by restricting the search space when examining if a file
has been altered in the plugin repository.

## Contribute

See the Development guide section for information about how to contribute to this plugin.

## I saw a bug, I have a feature request or a suggestion

Please file a [GitHub issue](https://github.com/mattermost/mattermost-plugin-nagios/issues), it will be very useful!

Pull Requests are welcome! You can contact us on the [Mattermost Community ~Plugin: Nagios channel](https://community.mattermost.com/core/channels/plugin-nagios).

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options.

## Share Feedback

Feel free to create a [GitHub Issue](https://github.com/mattermost/mattermost-plugin-nagios/issues) or join the [Nagios Plugin channel](https://community.mattermost.com/core/channels/plugin-nagios) on the Mattermost Community server to discuss.

## Security Vulnerability Disclosure

Please report any security vulnerability to [https://mattermost.com/security-vulnerability-report/](https://mattermost.com/security-vulnerability-report/).
