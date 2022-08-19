package gonmap

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type match struct {
	//match <Service> <pattern> <patternopt> [<Versioninfo>]
	Soft          bool
	Service       string
	Pattern       string
	PatternRegexp *regexp.Regexp
	Versioninfo   TcpFinger
}

var MATCH_LOAD_REGEXPS = []*regexp.Regexp{
	regexp.MustCompile("^([a-zA-Z0-9-_./]+) m\\|([^|]+)\\|([is]{0,2})(?: (.*))?$"),
	regexp.MustCompile("^([a-zA-Z0-9-_./]+) m=([^=]+)=([is]{0,2})(?: (.*))?$"),
	regexp.MustCompile("^([a-zA-Z0-9-_./]+) m%([^%]+)%([is]{0,2})(?: (.*))?$"),
	regexp.MustCompile("^([a-zA-Z0-9-_./]+) m@([^@]+)@([is]{0,2})(?: (.*))?$"),
}

var MATCH_VARSIONINFO_REGEXPS = map[string]*regexp.Regexp{
	"PRODUCTNAME": regexp.MustCompile("p/([^/]+)/"),
	"VERSION":     regexp.MustCompile("v/([^/]+)/"),
	"INFO":        regexp.MustCompile("i/([^/]+)/"),
	"HOSTNAME":    regexp.MustCompile("h/([^/]+)/"),
	"OS":          regexp.MustCompile("o/([^/]+)/"),
	"DEVICE":      regexp.MustCompile("d/([^/]+)/"),
}

var MATCH_VERSIONINFO_HELPER_P_REGEXP = regexp.MustCompile(`\$P\((\d)\)`)
var MATCH_VERSIONINFO_HELPER_REGEXP = regexp.MustCompile(`\$(\d)`)

func newMatch() *match {
	return &match{
		Soft:        false,
		Service:     "",
		Pattern:     "",
		Versioninfo: newFinger(),
	}
}

func (m *match) load(s string, soft bool) bool {
	var MATCH_LOAD_REGEXP *regexp.Regexp
	for _, r := range MATCH_LOAD_REGEXPS {
		if r.MatchString(s) {
			MATCH_LOAD_REGEXP = r
		}
	}
	if MATCH_LOAD_REGEXP == nil {
		return false
	}
	args := MATCH_LOAD_REGEXP.FindStringSubmatch(s)
	m.Soft = soft
	m.Service = args[1]
	m.Pattern = args[2]
	m.Versioninfo.Service = m.Service
	m.Versioninfo.ProductName = m.getVersionInfo(s, "PRODUCTNAME")
	m.Versioninfo.Version = m.getVersionInfo(s, "VERSION")
	m.Versioninfo.Info = m.getVersionInfo(s, "INFO")
	m.Versioninfo.Hostname = m.getVersionInfo(s, "HOSTNAME")
	m.Versioninfo.OperatingSystem = m.getVersionInfo(s, "OS")
	m.Versioninfo.DeviceType = m.getVersionInfo(s, "INFO")

	m.PatternRegexp = m.getPatternRegexp(m.Pattern, args[3])
	return true
}

func (m *match) getPatternRegexp(pattern string, opt string) *regexp.Regexp {
	pattern = strings.ReplaceAll(pattern, `\0`, `\x00`)
	if opt != "" {
		if pattern[:1] == "^" {
			pattern = fmt.Sprintf("^(?%s:%s", opt, pattern[1:])
		} else {
			pattern = fmt.Sprintf("(?%s:%s", opt, pattern)
		}
		if pattern[len(pattern)-1:] == "$" {
			pattern = fmt.Sprintf("%s)$", pattern[:len(pattern)-1])
		} else {
			pattern = fmt.Sprintf("%s)", pattern)
		}
	}
	//pattern = regexp.MustCompile(`\\x[89a-f][0-9a-f]`).ReplaceAllString(pattern,".")
	return regexp.MustCompile(pattern)
}

func (m *match) getVersionInfo(s string, regID string) string {
	if MATCH_VARSIONINFO_REGEXPS[regID].MatchString(s) {
		return MATCH_VARSIONINFO_REGEXPS[regID].FindStringSubmatch(s)[1]
	} else {
		return ""
	}
}

func (m *match) makeVersionInfo(s string) TcpFinger {
	f := newFinger()
	//fmt.Println(s)
	f.Info = m.makeVersionInfoSubHelper(s, m.Versioninfo.Info)
	f.DeviceType = m.makeVersionInfoSubHelper(s, m.Versioninfo.DeviceType)
	f.Hostname = m.makeVersionInfoSubHelper(s, m.Versioninfo.Hostname)
	f.OperatingSystem = m.makeVersionInfoSubHelper(s, m.Versioninfo.OperatingSystem)
	f.ProductName = m.makeVersionInfoSubHelper(s, m.Versioninfo.ProductName)
	f.Version = m.makeVersionInfoSubHelper(s, m.Versioninfo.Version)
	f.Service = m.makeVersionInfoSubHelper(s, m.Versioninfo.Service)
	return f
}

func (m *match) makeVersionInfoSubHelper(s string, pattern string) string {
	if len(m.PatternRegexp.FindStringSubmatch(s)) == 1 {
		return pattern
	}
	if pattern == "" {
		return pattern
	}
	sArr := m.PatternRegexp.FindStringSubmatch(s)

	if MATCH_VERSIONINFO_HELPER_P_REGEXP.MatchString(pattern) {
		pattern = MATCH_VERSIONINFO_HELPER_P_REGEXP.ReplaceAllStringFunc(pattern, func(repl string) string {
			a := MATCH_VERSIONINFO_HELPER_P_REGEXP.FindStringSubmatch(repl)[1]
			return "$" + a
		})
	}

	if MATCH_VERSIONINFO_HELPER_REGEXP.MatchString(pattern) {
		pattern = MATCH_VERSIONINFO_HELPER_REGEXP.ReplaceAllStringFunc(pattern, func(repl string) string {
			i, _ := strconv.Atoi(MATCH_VERSIONINFO_HELPER_REGEXP.FindStringSubmatch(repl)[1])
			return sArr[i]
		})
	}
	pattern = strings.ReplaceAll(pattern, "\n", "")
	pattern = strings.ReplaceAll(pattern, "\r", "")
	return pattern
}
