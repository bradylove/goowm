package config

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseHexColor(h string) int {
	h = strings.TrimPrefix(h, "#")
	c64, err := strconv.ParseInt(h, 16, 32)
	if err != nil {
		panic(fmt.Errorf("Failed to parse color #%s: %s", h, err))
	}

	return int(c64)
}
