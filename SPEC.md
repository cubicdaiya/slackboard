# Specification for Slackboard

Slackboard is a proxy server for Slack. It accepts a HTTP request.

## API

Slackboard has some APIs.

 * [POST /notify](#post-notify)
 * [POST /notify-directly](#post-notify-directly)
 * [GET /stat/go](#get-statgo)
 * [GET /app/config](#get-appconfig)

### POST /notify

Accepts a HTTP request for notification to Slack.

The JSON below is a request-body example.

```json
{
    "tag": "random",
    "host": "localhost",
    "text": "notification text",
    "sync": false
}
```

|name |type  |description               |required|note            |
|-----|------|--------------------------|--------|----------------|
|tag  |string|tag for selecting channel |o       |                |
|host |string|hostname(client)          |-       |                |
|text |string|notification text         |o       |                |
|level|string|message/info/warn/crit    |-       |default: message|
|title|string|message title             |-       |                |
|sync |bool  |synchronous notification  |-       |default: false  |



The JSON below is a response-body example from Slackboard. In this case, a status is 200(OK).

```json
{
    "message" : "ok",
}
```

When Slackboard receives an invalid request(for example, malformed body is included), a status of response it returns is 400(Bad Request).

### POST /notify-directly

Accepts a HTTP request for notification to Slack.

The JSON below is a request-body example.

```json
{
    "payload": {
        "channel": "random",
        "username": "slackboard",
        "icon_emoji": ":clipboard:",
        "text": "notification text",
        "parse": "full",
        "attachments": [...]
    },
    "sync": false
}
```

|name               |type  |description             |required|note                                       |
|-------------------|------|------------------------|--------|-------------------------------------------|
|payload            |object|payload object          |o       |                                           |
|payload.channel    |string|channel name            |o       |                                           |
|payload.username   |string|user name               |-       |default: slackboard                        |
|payload.icon_emoji |string|icon emoji              |-       |default: :clipboard:                       |
|payload.text       |string|notification text       |o       |                                           |
|payload.parse      |string|parsing mode            |-       |default: full                              |
|payload.attachments|object|attachments             |-       |see https://api.slack.com/docs/attachments |
|sync               |bool  |synchronous notification|-       |default: false                             |

The JSON below is a response-body example from Slackboard. In this case, a status is 200(OK).

```json
{
    "message" : "ok",
}
```

When Slackboard receives an invalid request(for example, malformed body is included), a status of response it returns is 400(Bad Request).

### GET /stat/go

Returns a statictics for golang-runtime. See [golang-stats-api-handler](https://github.com/fukata/golang-stats-api-handler) about details.

### GET /app/config

Returns a current configuration for Slackboard.
