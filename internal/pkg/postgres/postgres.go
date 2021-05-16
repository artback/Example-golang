package postgres

import (
	"fmt"
	"strings"
)

func BuildValuesString(strFmt string, length int) string {
	var valueStrings []string
	for i := 0; i < length; i++ {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d)", i*2+1, i*2+2))
	}
	return fmt.Sprintf(strFmt, strings.Join(valueStrings, ","))
}
