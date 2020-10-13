# Mattermost Nagios Plugin

<!-- TODO(amwolff): add CI badges and stuff. -->

**Maintainers**: @amwolff & @DanielSz50

A Nagios plugin for Mattermost. Supports Nagios Core >= 4.4.x.

## About

This plugin allows you to

- [x] get logs to specific systems without leaving the Mattermost
    - get alerts and notifications instantly delivered, resembling the `showlog.cgi` UI
- [x] receive system monitoring reports on a subscribed channel
    - you will be frequently informed which hosts and/or services have abnormal state
- [ ] (in progress) receive notifications about changes to configuration on a subscribed channel
    - anytime a change has been made to Nagios configuration, you will receive a diff between the old and the new version

Ultimately, this will make you or your team more productive and make the experience with Nagios smoother.

### Audience

This guide is for Mattermost System Admins setting up the Nagios plugin and Mattermost users who want information about the plugin functionality.

### Important notice

If you are a Nagios admin/user and think there is something this plugin lacks or something that it does could be done the other way around, let us know!
We are trying to develop this plugin basing on users' needs.
If there is a certain feature you or your team needs, open up an issue and explain your needs.
We will be happy to help.

## Installation

1. Download the latest version of the plugin from the [releases page](https://github.com/ulumuri/mattermost-plugin-nagios/releases)
2. In Mattermost, go to **System Console → Plugins → Management**
3. Upload the plugin in the **Upload Plugin** section
4. Configure the plugin before you enable it

## Configuration

1. Enter the URL for your Nagios instance
    1. In Mattermost, go to **System Console → Plugins → Nagios**
    2. Set the **Nagios URL**
        1. Remember to add `http://` or `https://` at the beginning!
        2. Example URL: `https://nagios.fedoraproject.org`
2. Click *Save* to save the settings
3. In Mattermost, go to **System Console → Plugins → Management** and click *Enable* underneath the Nagios plugin
4. The plugin is now ready to use! 🎉

<!-- TODO(amwolff): add setup instructions for the configuration files watcher. -->

## Using the plugin

### Slash commands

## Contributing
