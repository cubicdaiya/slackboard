# Slackboard

Slackboard is a proxy server for Slack.

## Status

Slackboard is production ready.

## Features

 * `slackboard`
  * A proxy server for slack
 * `slackboard-cli`
  * A client for `slackboard`
 * `slackboard-log`
  * A client like [`cronlog`](https://github.com/kazuho/kaztools/blob/master/cronlog) for `slackboard`

## Build

```
make gom
make bundle
make
```

## Configuration

See [CONFIGURATION.md](https://github.com/cubicdaiya/slackboard/blob/master/CONFIGURATION.md) about details.

## Specification

See [SPEC.md](https://github.com/cubicdaiya/slackboard/blob/master/SPEC.md) about details.

## Run

```
slackboard -c conf/slackboard.toml
```

## Client for Slackboard

`slackboard-cli` is a client for `slackboard`. It reads `stdin` and sends a message to `slackboard`.

```
echo message | slackboard-cli -t test -s slackboard-host:29800
```

### Synchronous notification with slackboard-cli

From v0.3.0 `slackboard-cli` sends a notification-request to `slackboard` asynchronously by default.
If you want to send a notification-request to `slackboard` synchronously, you may add the option `-sync` to `slackboard-cli`.

```
echo message | slackboard-cli -t test -s slackboard-host:29800 -sync
```

## Notification with slackboard-log

`slackboard-log` is a client for `slackboard` also. `slackboard-log` is an utility like [`cronlog`](https://github.com/kazuho/kaztools/blob/master/cronlog).
It sends a notification to `slackboard` when the command after `--` failed.

```
slackboard-log -s 127.0.0.1:29800 -t test -- some-command
```

## License

Copyright 2014-2015 Tatsuhiko Kubo


Licensed under the MIT License.
