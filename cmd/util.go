package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var wfRegex = regexp.MustCompile("^([0-9]+)(s|m|h|d)$")

func parseWatchFrequency(wf string) (*time.Duration, error) {
	s := wfRegex.FindStringSubmatch(wf)
	if len(s) != 3 {
		return nil, fmt.Errorf("invalid watch frequency %q", wf)
	}
	v, err := strconv.Atoi(s[1])
	if err != nil {
		return nil, err
	}

	var d time.Duration
	switch s[2] {
	case "s":
		d = time.Duration(v) * time.Second
	case "m":
		d = time.Duration(v) * time.Minute
	case "h":
		d = time.Duration(v) * time.Hour
		return &d, nil
	case "d":
		d = time.Duration(v) * 24 * time.Hour
	default:
		return nil, fmt.Errorf("couldn't find any matching frequencies for %q", wf)
	}

	return &d, nil
}
