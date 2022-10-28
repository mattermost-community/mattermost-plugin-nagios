# Development

This repository uses the [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template). Therefore, developing this plugin is roughly the same as it is with every plugin using the template. All the necessary steps to develop are in the template's repository.

If you are a Nagios admin/user and think there is something this plugin lacks or something that it does could be done another way, let us know! We are trying to develop this plugin based on users' needs. If there is a certain feature you or your team needs, open up an issue, and explain your needs. We will be happy to help.

This plugin only contains a server portion. Read our documentation about the [Developer Workflow](https://developers.mattermost.com/extend/plugins/developer-workflow/) and [Developer Setup](https://developers.mattermost.com/extend/plugins/developer-setup/) for more information about developing and extending plugins.

## Developing the watcher

To build the watcher, you can use the following command:

```sh
env GOOS=linux GOARCH=amd64 go build -o dist/watcherX.Y.Z.linux-amd64 -a -v cmd/watcher/main.go
```

Of course, you can build the watcher for other operating systems and architectures too.

## I saw a bug, I have a feature request or a suggestion

Please file a [GitHub issue](https://github.com/mattermost/mattermost-plugin-nagios/issues), it will be very useful!

Pull Requests are welcome! You can contact us on the [Mattermost Community ~Plugin: Nagios channel](https://community.mattermost.com/core/channels/plugin-nagios).

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options.

### Deploying with Local Mode

If your Mattermost server is running locally, you can enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode) to streamline deploying your plugin. After configuring it, just run:

```sh
make deploy
```

### Deploying with credentials

Alternatively, you can authenticate with the server's API with a [personal access token](https://docs.mattermost.com/developer/personal-access-tokens.html):

```sh
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
```
