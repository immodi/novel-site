package search

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseTotalChapterRange parses strings like:
//
//	""         -> no filter (nil, nil)
//	"0"        -> no filter (nil, nil)   // <- important for your "All" select option
//	"100"      -> exact (100,100)
//	"1,49"     -> range (1,49)
//	"1000,500" -> swapped -> (500,1000)
func ParseTotalChapterRange(param string) (min, max *int, err error) {
	param = strings.TrimSpace(param)
	if param == "" || param == "0" {
		// "0" is your "All" option so treat it as no filter
		return nil, nil, nil
	}

	parts := strings.Split(param, ",")
	switch len(parts) {
	case 1:
		v, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, nil, fmt.Errorf("invalid totalchapter value: %w", err)
		}
		// exact match -> min == max
		minVal, maxVal := v, v
		return &minVal, &maxVal, nil

	case 2:
		a, errA := strconv.Atoi(strings.TrimSpace(parts[0]))
		b, errB := strconv.Atoi(strings.TrimSpace(parts[1]))
		if errA != nil || errB != nil {
			return nil, nil, fmt.Errorf("invalid totalchapter range")
		}
		if a > b {
			a, b = b, a
		}
		minVal, maxVal := a, b
		return &minVal, &maxVal, nil

	default:
		return nil, nil, fmt.Errorf("invalid totalchapter format")
	}
}
