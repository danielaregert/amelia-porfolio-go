package sanitize

import "github.com/microcosm-cc/bluemonday"

var policy = bluemonday.UGCPolicy()

func HTML(input string) string {
	return policy.Sanitize(input)
}
