package eskomlol

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

// Schedule represents a loadshedding schedule for specific stage.
type Schedule struct {
	Stage
	Times []ScheduleItem
}

// ScheduleItem represents a single instance of loadshedding.
type ScheduleItem struct {
	Start time.Time
	End   time.Time
}

// rawItem is used to parse the raw values from the Eskom page.
type rawItem struct {
	date, time string
}

// parseScheduleHTML iterates the schedule HTML page and parses any date and time combinations.
//
// Any non-ErrElementNotFound will be added to the returned error object
func parseScheduleHTML(data []byte) ([]rawItem, error) {
	errs := make([]string, 0)
	res := make([]rawItem, 0)
	document := soup.HTMLParse(string(data))
	days := document.FindAll("div", "class", "scheduleDay")
	for _, day := range days {
		dateDiv := day.Find("div", "class", "dayMonth")
		if dateDiv.Error != nil {
			if dateDiv.Error.(soup.Error).Type != soup.ErrElementNotFound {
				errs = append(errs, dateDiv.Error.Error())
			}
			continue
		}
		date := strings.TrimSpace(dateDiv.Text())
		timesA := day.Find("a")
		if timesA.Error != nil {
			if timesA.Error.(soup.Error).Type != soup.ErrElementNotFound {
				errs = append(errs, timesA.Error.Error())
			}
			continue
		}
		time := strings.TrimSpace(timesA.Text())
		res = append(res, rawItem{date: date, time: time})
	}

	var err error
	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, "; "))
	}

	return res, err
}

// makeScheduleItems parses the given rawItems into ScheduleItems.
//
// All times returned will be in SAST.
// Since the raw data has no value for the year, logic is applied to ensure
// that the year increments correctly for cases such as:
// A schedule going from 1 December 2021 to 31 January 2022.
//
// Any errors occurred while parsing the rawItems will be added to the returned
// error object
func makeScheduleItems(rawItems []rawItem, now time.Time) ([]ScheduleItem, error) {
	res := make([]ScheduleItem, 0)
	errs := make([]string, 0)
	currentYear := now.Year()

	var previousMonth string

	parseFormat := "Mon, 02 Jan 2006 15:04 MST"

	for _, rawItem := range rawItems {
		month := rawItem.date[len(rawItem.date)-3:]
		if month != "" && month != previousMonth && month == "Jan" {
			currentYear++
		}
		previousMonth = month
		rawTimeParts := strings.Split(rawItem.time, " - ")
		lowerTime, upperTime := rawTimeParts[0], rawTimeParts[1]

		startTime, err := time.Parse(parseFormat, fmt.Sprintf("%s %d %s SAST", rawItem.date, currentYear, lowerTime))
		if err != nil {
			errs = append(errs, err.Error())
		}
		endTime, err := time.Parse(parseFormat, fmt.Sprintf("%s %d %s SAST", rawItem.date, currentYear, upperTime))
		if err != nil {
			errs = append(errs, err.Error())
		}

		res = append(res, ScheduleItem{Start: startTime, End: endTime})
	}

	var err error
	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, "; "))
	}

	return res, err
}

// scheduleFromHTML parses a Schedule from the given HTML page contents.
func scheduleFromHTML(data []byte, stage Stage, now time.Time) (Schedule, error) {
	res, err := parseScheduleHTML(data)
	if err != nil {
		return Schedule{}, err
	}

	scheduleItems, err := makeScheduleItems(res, now)
	if err != nil {
		return Schedule{}, err
	}

	return Schedule{
		Stage: stage,
		Times: scheduleItems,
	}, nil
}
