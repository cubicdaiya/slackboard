# Configuration for Slackboard

A configuration file format for Slackboard is [TOML](https://github.com/toml-lang/toml).

A configuration for Gaurun has some sections. A example is [here](conf/slackboard.toml).

 * [Core Section](#core-section)
 * [Tag Section](#tag-section)
 * [Log Section](#log-section)

## Core Section

|name     |type  |description                    |default|note                                 |
|---------|------|-------------------------------|-------|-------------------------------------|
|port     |string|port number or unix socket path|29800  |e.g.)29800, unix:/tmp/slackboard.sock|
|slack_url|string|web hook url for slack         |       |                                     |

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
