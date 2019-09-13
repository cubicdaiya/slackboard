# Configuration for Slackboard

A configuration file format for Slackboard is [TOML](https://github.com/toml-lang/toml).

A configuration for Slackboard has some sections. A example is [here](conf/slackboard.toml).

 * [Core Section](#core-section)
 * [Tag Section](#tag-section)
 * [Log Section](#log-section)

## Core Section

|name               |type  |description                    |default|note                                 |
|-------------------|------|-------------------------------|-------|-------------------------------------|
|port               |string|port number or unix socket path|29800  |e.g.)29800, unix:/tmp/slackboard.sock|
|slack_url          |string|Incomming Webhook url for slack|       |If slack_token is specified, this field won't be used.|
|slack_token        |string|Slack token to post messages   |       |token with permission to use [chat.postMessage](https://api.slack.com/methods/chat.postMessage)|
|qps                |int   |Queries to slack Per Second    |0      |0 means unlimited. See also [slack api docs](https://api.slack.com/docs/rate-limits) |
|max_delay_duration |int   |Allowable delay message seconds|       |must be specified when using qps > 0 and planning to process requests with sync=false |

## Tag Section

|name      |type  |description              |default    |note                    |
|----------|------|-------------------------|-----------|------------------------|
|tag       |string|tag for selecting channel|           |empty tag is not allowed|
|channel   |string|channel name             |#random    |                        |
|username  |string|user name                |slackboard |                        |
|icon_emoji|string|emoji icon               |:clipboard:|                        |
|parse     |string|slack parse mode         |full       |                        |

## Log Section

|name      |type  |description    |default|note                             |
|----------|------|---------------|-------|---------------------------------|
|access_log|string|access log path|stdout |                                 |
|error_log |string|error log path |stderr |                                 |
|level     |string|log level      |error  |panic,fatal,error,warn,info,debug|

`access_log` and `error_log` are allowed to give not only file-path but `stdout` and `stderr`.
