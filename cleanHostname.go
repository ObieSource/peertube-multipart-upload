package main

import (
	"fmt"
	"strings"
)

func CleanHostname(raw string) string {
	/*
		Depending on how the hostname is entered
		(likely via copy-paste), there may be a
		trailing / to signal a path of "/".

		This function removes that if it is present,
		as the root / is added by all urls that
		are constructed by the program.

		Note that this function does not do any
		validation of the url at all.
	*/
	out := strings.ToLower(strings.TrimRight(raw, "/"))
	if !strings.HasPrefix(out, "https://") && !strings.HasPrefix(out, "http://") && out != "" {
		out = fmt.Sprintf("https://%s", out)

	}
	return out
}
