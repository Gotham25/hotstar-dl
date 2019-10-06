package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//CopyMap creates a copy of given map of string.
func CopyMap(m map[string]string) map[string]string {
	cp := make(map[string]string)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

func padZeroRight(num int64) int64 {
	tmp := fmt.Sprintf("%-13d", num)
	tmp = strings.Replace(tmp, " ", "0", -1)
	paddedNum, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		panic(err)
	}
	return paddedNum
}

func countDigits(i int64) (count int64) {
	for i != 0 {
		i /= 10
		count = count + 1
	}
	return count
}

//GetDateStr parses given time in milliseconds to human readable Date string.
func GetDateStr(timeFloat64 float64) string {
	timeMillis := int64(timeFloat64)
	if countDigits(timeMillis) != 13 {
		timeMillis = padZeroRight(timeMillis)
	}
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		panic(err)
	}
	return time.Unix(0, timeMillis*int64(time.Millisecond)).In(location).String()
}

//MakeRange generates a range of numbers for the given range
func MakeRange(min, max int) []int {
	contents := make([]int, max-min+1)
	for index := range contents {
		contents[index] = min + index
	}
	return contents
}

//ReSubMatchMap returns the regex submatches as a map
func ReSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}

	return subMatchMap
}
