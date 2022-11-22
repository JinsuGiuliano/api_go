package version

import (
	"fmt"
	"runtime"
	"strings"
)

const poweredBy = "Powered by Gustavo Giuliano & Co"

var versionBase = "1.6.0" // nolint
var buildHash string      // nolint
var buildBranch string    // nolint
var builtTime string      // nolint
var environment string    // nolint

func Base() string {
	return versionBase
}

func BuildHash() string {
	return buildHash
}

func BuildBranch() string {
	return buildBranch
}

func BuiltTime() string {
	return builtTime
}

func Environment() string {
	return environment
}

func Name() string {
	bh := buildHash
	if len(buildHash) >= 8 {
		bh = buildHash[0:8]
	}
	return versionBase + "-" + bh
}

func Info() string {
	bt := strings.Replace(builtTime, "-", " ", -1)
	bh := buildHash
	if len(buildHash) >= 8 {
		bh = buildHash[0:8]
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf(" Giuliano v%s-%s \n", versionBase, bh))
	sb.WriteString(fmt.Sprintf(" Build v%s-%s-%s, built on %s \n", versionBase, buildBranch, bh, bt))
	sb.WriteString(fmt.Sprintln("", poweredBy, ""))
	sb.WriteString(fmt.Sprintf(" Go %s \n", runtime.Version()))
	sb.WriteString(fmt.Sprintf(" Environment %s \n", environment))

	return sb.String()
}
