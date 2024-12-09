package parseurl

import (
	"net/url"
	"strconv"
	"time"
)

func ParseString(queryValuer url.Values, key, defaultValue string) string {
	value := queryValuer.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func ParseInt(queryValuer url.Values, key string, defaultValue int) int {
	valueString := queryValuer.Get(key)

	if valueString == "" {
		return defaultValue
	}

	if value, err := strconv.Atoi(valueString); err == nil {
		return value
	}
	return defaultValue
}

func ParseTime(queryValuer url.Values, key string, defaultValue time.Time) time.Time {
	value := queryValuer.Get(key)
	if value == "" {
		return defaultValue
	}
	res, err := time.Parse("02.01.2006", value)
	if err != nil {
		return defaultValue
	}
	return res
}
