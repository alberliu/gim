package util

import (
	"strconv"
	"strings"
)

func In(ids []int64) string {
	build := strings.Builder{}
	build.WriteString("(")
	for i := range ids {
		build.WriteString(strconv.FormatInt(ids[i], 10))
		if i != len(ids)-1 {
			build.WriteString(",")
		}
	}
	build.WriteString(")")
	return build.String()
}
