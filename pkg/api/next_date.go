package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// afterNow returns true if date is after now (both normalized to day)
func afterNow(date, now time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.After(now)
}

// NextDate calculates the next date for a repeating task
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(Layout, dstart)
	if err != nil {
		return "", errors.New("error in date parsing")
	}
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}
	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		return nextDay(now, date, parts)
	case "y":
		return nextYear(now, date, parts)
	case "w":
		return nextWeek(now, date, parts)
	case "m":
		return nextMonth(now, date, parts)
	default:
		return "", fmt.Errorf("wrong repeat date format: %s", parts[0])
	}
}

// nextDateHandler is an HTTP handler for calculating the next date
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nowStr := r.FormValue("now")
	var now time.Time
	if nowStr == "" {
		now = time.Now().UTC().Truncate(24 * time.Hour)
	} else {
		var err error
		now, err = time.Parse(Layout, nowStr)
		if err != nil {
			http.Error(w, "incorrect parameter 'now'", http.StatusBadRequest)
			return
		}
		now = now.Truncate(24 * time.Hour)
	}
	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	if _, err := fmt.Fprint(w, nextDate); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}
