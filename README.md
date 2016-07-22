# hue-alert: Hue light alert

Simple and secure notification system for
[Philips Hue](http://www2.meethue.com/en-us/) lights. Supports Gmail and Slack
using Oauth for authentication. All authentication credentials are stored
locally and authentication is done using Oauth without any third-party servers.

## Create Oauth API keys

Google and Slack Oauth API keys will be needed to authenticate with Oauth. This
is done by creating an Oauth application using the links below.

#### Google setup

[Google API Credentials](https://console.developers.google.com/apis/credentials)

For Google select create *Oauth client ID* then select *Web application*. Set
the name to `Hue Alert`. Then enter `http://localhost:9300` for the *Authorized
javascript origins* and `http://localhost:9300/callback/google` for the
*Authorized redirect URIs*

#### Slack setup

[Slack Create App](https://api.slack.com/apps/new)

For Slack set the name and descriptions to `Hue Alert`. Then enter
`http://localhost:9300/callback/slack` for the *Redirect URI(s)*.

## Hue setup

First connect to the Hue bridge using the command below. Then run the next
command to select which lights to use for notifications. The `solid` mode will
keep the light on until the notification is read, the other modes will flash
the light. The hostname of the Hue bridge can be found in the Network settings
by disabling DHCP. The IP address will then be shown.

```
$ hue-alert hue-setup
Enter Hue Bridge hostname: 10.0.0.161
Press the link button on top of the Hue bridge then press enter...
Hue successfully linked
$ hue-alert hue-lights
Add Bloom (Color light)? [y/N] y
Light 'Bloom' has been added...
Enter alert light brightness: [1-254] 254
Enter alert light mode: [solid,slow,medium,fast] solid
```

## Adding accounts

Run the commands below to generate an Oauth url. Then open the url in your
browser to authorize the Oauth application. Accounts can be listed by running
`hue-alert accounts`. Accounts can be removed by running
`hue-alert account-remove google username@gmail.com` or
`hue-alert account-remove slack username`.

```
$ hue-alert google-add
Open URL below to authenticate Google account:
https://accounts.google.com/o/oauth2/auth
[INFO] ▶ server: Starting oauth server ◆ address=":9300"
Google account successfully authenticated
$ hue-alert slack-add
Open URL below to authenticate Slack account:
https://slack.com/oauth/authorize
[INFO] ▶ server: Starting oauth server ◆ address=":9300"
Slack account successfully authenticated
```

## Test

Running the command `hue-alert hue-color f00` should set the color of the
selected lights to red.

## Start hue-alert

Run `hue-alert start` to start monitoring for notifications. Notifications will
be triggered for unread email in the Gmail inbox and unread messages in any
Slack channel.
