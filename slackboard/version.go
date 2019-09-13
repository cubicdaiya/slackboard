package slackboard

import (
	"fmt"
	"runtime"
)

func PrintVersion() {
	fmt.Printf(`slackboard %s
Compiler: %s %s
Copyright (C) 2014-2019 Tatsuhiko Kubo <cubicdaiya@gmail.com>
`,
		Version,
		runtime.Compiler,
		runtime.Version())
}
