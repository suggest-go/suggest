package http

import (
	"errors"
	"net/http"
	"strconv"
)

// FormTopKValue returns the first value for the named component of the query and
// validates it with topK value restrictions
func FormTopKValue(r *http.Request, field string, defaultVal int) (int, error) {
	val, err := FormIntValue(r, field, defaultVal)

	if err != nil {
		return 0, err
	}

	if val < 0 {
		return 0, errors.New("topK should be positive integer")
	}

	return val, nil
}

// FormSimilarityValue returns the first value for the named component of the query and
// validates it with similarity value restrictions
func FormSimilarityValue(r *http.Request, field string, defaultVal float64) (float64, error) {
	val, err := FormFloatValue(r, field, defaultVal)

	if err != nil {
		return 0, err
	}

	if val < 0 || val > 1 {
		return 0, errors.New("similarity should be in [0, 1] range")
	}

	return val, nil
}

func FormIntValue(r *http.Request, field string, defaultVal int) (int, error) {
	val := r.FormValue(field)

	if val == "" {
		return defaultVal, nil
	}

	i64, err := strconv.ParseInt(val, 10, 0)

	if err != nil {
		return 0, err
	}

	return int(i64), nil
}

func FormFloatValue(r *http.Request, field string, defaultVal float64) (float64, error) {
	val := r.FormValue(field)

	if val == "" {
		return defaultVal, nil
	}

	f64, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return 0, err
	}

	return f64, nil
}

