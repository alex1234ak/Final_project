package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Реализация дневного ("d") правила повторения
func nextDay(now, date time.Time, parts []string) (string, error) {
	//Проверка на правильность правила
	if len(parts) != 2 {
		return "", errors.New("wrong rule for 'd' format")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", errors.New("interval must be a number")
	}

	if interval < 1 || interval > 400 {
		return "", errors.New("selected interval must be 1-400")
	}

	//Увеличиваем время, пока не будет больше текущего
	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			return date.Format(Layout), nil
		}
	}
}

// Реализация годового ("y") правила повторения
func nextYear(now, date time.Time, parts []string) (string, error) {
	//Проверка на правильность правила
	if len(parts) != 1 {
		return "", errors.New("wrong rule for 'y' format")
	}

	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			return date.Format(Layout), nil
		}
	}
}

// Реализация недельного ("w") правила повторения
func nextWeek(now, date time.Time, parts []string) (string, error) {
	//Проверка, что есть перечисление дней
	if len(parts) != 2 {
		return "", errors.New("wrong rule for 'w' format")
	}

	days, err := parseIntList(parts[1], 1, 7)
	if err != nil {
		return "", fmt.Errorf("invalid week days: %w", err)
	}
	if len(days) == 0 {
		return "", errors.New("empty week days list")
	}

	//Проверка удовлетворяет ли стартовая дата условию "позже чем сейчас"
	if afterNow(date, now) && correctWeekDay(date, days) {
		return date.Format(Layout), nil
	}

	//Перебор дней
	for i := 0; i < 1500; i++ {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) && correctWeekDay(date, days) {
			return date.Format(Layout), nil
		}
	}

	return "", errors.New("date not found within max search period")
}

// Реализация месячного ("m") правила повторения
func nextMonth(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("wrong rule for 'm' format")
	}

	days, months, err := parseMonthRule(parts)
	if err != nil {
		return "", fmt.Errorf("invalid month rule: %w", err)
	}

	//newDate := date

	//Проверка удовлетворяет ли стартовая дата условию "позже чем сейчас"
	if afterNow(date, now) && correctMonthDay(date, days, months) {
		return date.Format(Layout), nil
	}

	//Перебор дней
	for i := 0; i < 1500; i++ {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) && correctMonthDay(date, days, months) {
			return date.Format(Layout), nil
		}
	}

	return "", errors.New("date not found within max search period")
}

// Парсинг строки в диапазоне min - max для поиска целых чисел с разделителем ","
func parseIntList(s string, min, max int) ([]int, error) {
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			continue
		}

		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, errors.New("not a number")
		}

		if val < min || val > max {
			return nil, fmt.Errorf("value %d out of range [%d,%d]", val, min, max)
		}

		result = append(result, val)
	}
	return result, nil
}

// Функция проверяет дату на соответсвие дням недели переданным в массиве
func correctWeekDay(date time.Time, days []int) bool {
	wday := date.Weekday()
	wdayInt := int(wday)
	if wday == time.Sunday {
		wdayInt = 7
	}
	for _, d := range days {
		if d == wdayInt {
			return true
		}
	}
	return false
}

// Парсинг данных из правила "repeat" для разделения их на слайс дней и месяцев
func parseMonthRule(parts []string) (days []int, months []int, err error) {

	// Парсинг дней (-2..31)
	days, err = parseIntList(parts[1], -2, 31)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid days: %w", err)
	}

	// Парсинг месяцев
	if len(parts) > 2 {
		months, err = parseIntList(parts[2], 1, 12)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid months: %w", err)
		}
	}

	return days, months, nil
}

// Проверка дней на соответствие правилу повторения
func correctMonthDay(date time.Time, days, months []int) bool {

	if len(months) > 0 {
		found := false
		for _, m := range months {
			if m == int(date.Month()) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Проверка дня
	lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	currentDay := date.Day()

	//	Отработка последнего и предпоследнего дня
	for _, d := range days {
		switch {
		case d > 0 && currentDay == d:
			return true
		case d == -1 && currentDay == lastDay:
			return true
		case d == -2 && currentDay == lastDay-1:
			return true
		}
	}
	return false
}
