package conv

import (
	"fmt"
	"strings"
)

func TypeOf(v interface{}, trimPackage ...bool) string {
	typText := fmt.Sprintf("%T", v)
	if trimPackage != nil && trimPackage[0] {
		split := strings.Split(typText, ".")
		if strings.HasPrefix(typText, "*") {
			return fmt.Sprintf("*%s", split[len(split)-1])
		}
		return split[len(split)-1]
	}
	return typText
}
