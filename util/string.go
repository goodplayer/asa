package util

import "strings"

func SplitHost(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

func SplitHostAndPort(hostport string) []string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return []string{hostport, ""}
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		portStart := i + 2
		if len(hostport) <= portStart {
			return []string{strings.TrimPrefix(hostport[:i], "["), ""}
		}
		return []string{strings.TrimPrefix(hostport[:i], "["), hostport[portStart:]}
	}
	if len(hostport) <= colon+1 {
		return []string{hostport[:colon], ""}
	}
	return []string{hostport[:colon], hostport[colon+1:]}
}
