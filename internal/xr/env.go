package xr

import (
	"regexp"
	"strings"
)

func expandVar(env []string, v string) string {
	prefix := v + "="
	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			return e[len(prefix):]
		}
	}
	return ""
}

var bashVarRE = regexp.MustCompile(`^\"?\$\{(\w+)\}\"?$`)

func expandArgs(env []string, args []string) []string {
	expanded := make([]string, len(args))
	for i, a := range args {
		m := bashVarRE.FindStringSubmatch(a)
		v := ""
		if 1 < len(m) {
			v = m[1]
		}
		if len(v) < 1 {
			expanded[i] = a
		} else {
			expanded[i] = expandVar(env, v)
		}
	}

	return expanded
}
