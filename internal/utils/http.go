package utils

import (
	"net/http"
	"strconv"
)

func GetQueryInt(r *http.Request, key string, defaultValue int32) int32 {
	values := r.URL.Query().Get(key)

	if values == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(values)
	if err != nil {
		return defaultValue
	}

	return int32(intValue)
}
