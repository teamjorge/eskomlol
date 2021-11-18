package eskomlol

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var testItems = []rawItem{
	{date: "Fri, 29 Oct", time: "04:00 - 06:30"},
	{date: "Sat, 30 Oct", time: "12:00 - 14:30"},
	{date: "Sun, 31 Oct", time: "20:00 - 22:30"},
	{date: "Mon, 01 Nov", time: "18:00 - 20:30"},
	{date: "Wed, 03 Nov", time: "02:00 - 04:30"},
	{date: "Thu, 04 Nov", time: "10:00 - 12:30"},
	{date: "Fri, 05 Nov", time: "16:00 - 18:30"},
	{date: "Sun, 07 Nov", time: "00:00 - 02:30"},
	{date: "Mon, 08 Nov", time: "08:00 - 10:30"},
	{date: "Tue, 09 Nov", time: "14:00 - 16:30"},
	{date: "Wed, 10 Nov", time: "22:00 - 00:30"},
	{date: "Fri, 12 Nov", time: "06:00 - 08:30"},
	{date: "Sat, 13 Nov", time: "12:00 - 14:30"},
	{date: "Sun, 14 Nov", time: "20:00 - 22:30"},
	{date: "Tue, 16 Nov", time: "04:00 - 06:30"},
	{date: "Wed, 17 Nov", time: "10:00 - 12:30"},
	{date: "Thu, 18 Nov", time: "18:00 - 20:30"},
	{date: "Sat, 20 Nov", time: "02:00 - 04:30"},
	{date: "Sun, 21 Nov", time: "08:00 - 10:30"},
	{date: "Mon, 22 Nov", time: "16:00 - 18:30"},
	{date: "Wed, 24 Nov", time: "00:00 - 02:30"},
	{date: "Thu, 25 Nov", time: "06:00 - 08:30"},
}

func TestParseScheduleHTML(t *testing.T) {
	testFile, err := os.Open("./test_data/schedule.html")
	if err != nil {
		t.Errorf("unexpected error opening test file: %v", err)
		return
	}

	testData, err := ioutil.ReadAll(testFile)
	if err != nil {
		t.Errorf("unexpected error reading test file: %v", err)
		return
	}

	rawItems, err := parseScheduleHTML(testData)
	if err != nil {
		t.Errorf("unexpected error parsing test file: %v", err)
		return
	}

	for index, rawItem := range rawItems {
		if testItems[index].date != rawItem.date {
			t.Errorf("expected item %d to date to be %s, got %s", index, testItems[index].date, rawItem.date)
		}
		if testItems[index].time != rawItem.time {
			t.Errorf("expected item %d to time to be %s, got %s", index, testItems[index].time, rawItem.time)
		}
	}
}

func TestMakeScheduleItemsNormal(t *testing.T) {
	loc, err := time.LoadLocation("Africa/Johannesburg")
	if err != nil {
		t.Errorf("unexpected error loading tz data: %v", err)
		return
	}
	now := time.Date(2021, 10, 27, 18, 00, 00, 0, loc)
	scheduleItems, err := makeScheduleItems(testItems, now)

	if err != nil {
		t.Errorf("unexpected error making schedule items: %v", err)
		return
	}

	expected := []ScheduleItem{
		{Start: time.Date(2021, 10, 29, 04, 0, 0, 0, loc), End: time.Date(2021, 10, 29, 06, 30, 00, 0, loc)},
		{Start: time.Date(2021, 10, 30, 12, 0, 0, 0, loc), End: time.Date(2021, 10, 30, 14, 30, 00, 0, loc)},
		{Start: time.Date(2021, 10, 31, 20, 0, 0, 0, loc), End: time.Date(2021, 10, 31, 22, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 01, 18, 0, 0, 0, loc), End: time.Date(2021, 11, 01, 20, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 03, 02, 0, 0, 0, loc), End: time.Date(2021, 11, 03, 04, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 04, 10, 0, 0, 0, loc), End: time.Date(2021, 11, 04, 12, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 05, 16, 0, 0, 0, loc), End: time.Date(2021, 11, 05, 18, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 07, 00, 0, 0, 0, loc), End: time.Date(2021, 11, 07, 02, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 8, 8, 0, 0, 0, loc), End: time.Date(2021, 11, 8, 10, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 9, 14, 0, 0, 0, loc), End: time.Date(2021, 11, 9, 16, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 10, 22, 0, 0, 0, loc), End: time.Date(2021, 11, 10, 00, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 12, 06, 0, 0, 0, loc), End: time.Date(2021, 11, 12, 8, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 13, 12, 0, 0, 0, loc), End: time.Date(2021, 11, 13, 14, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 14, 20, 0, 0, 0, loc), End: time.Date(2021, 11, 14, 22, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 16, 04, 0, 0, 0, loc), End: time.Date(2021, 11, 16, 06, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 17, 10, 0, 0, 0, loc), End: time.Date(2021, 11, 17, 12, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 18, 18, 0, 0, 0, loc), End: time.Date(2021, 11, 18, 20, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 20, 02, 0, 0, 0, loc), End: time.Date(2021, 11, 20, 04, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 21, 8, 0, 0, 0, loc), End: time.Date(2021, 11, 21, 10, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 22, 16, 0, 0, 0, loc), End: time.Date(2021, 11, 22, 18, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 24, 00, 0, 0, 0, loc), End: time.Date(2021, 11, 24, 02, 30, 00, 0, loc)},
		{Start: time.Date(2021, 11, 25, 06, 0, 0, 0, loc), End: time.Date(2021, 11, 25, 8, 30, 00, 0, loc)},
	}

	for index, scheduleItem := range scheduleItems {
		if !expected[index].Start.Equal(scheduleItem.Start) {
			t.Errorf("expected item %d to start to be %s, got %s", index, expected[index].Start, scheduleItem.Start)
		}
		if !expected[index].End.Equal(scheduleItem.End) {
			t.Errorf("expected item %d to end to be %s, got %s", index, expected[index].End, scheduleItem.End)
		}
	}
}

func TestMakeScheduleItemsYearRollover(t *testing.T) {
	loc, err := time.LoadLocation("Africa/Johannesburg")
	if err != nil {
		t.Errorf("unexpected error loading tz data: %v", err)
		return
	}

	var testItems = []rawItem{
		{date: "Fri, 26 Dec", time: "04:00 - 06:30"},
		{date: "Sat, 27 Dec", time: "12:00 - 14:30"},
		{date: "Sun, 28 Dec", time: "20:00 - 22:30"},
		{date: "Mon, 29 Dec", time: "18:00 - 20:30"},
		{date: "Wed, 30 Dec", time: "02:00 - 04:30"},
		{date: "Thu, 31 Dec", time: "10:00 - 12:30"},
		{date: "Fri, 01 Jan", time: "16:00 - 18:30"},
		{date: "Sun, 02 Jan", time: "00:00 - 02:30"},
		{date: "Mon, 03 Jan", time: "08:00 - 10:30"},
		{date: "Tue, 04 Jan", time: "14:00 - 16:30"},
		{date: "Wed, 05 Jan", time: "22:00 - 00:30"},
		{date: "Fri, 06 Jan", time: "06:00 - 08:30"},
		{date: "Sat, 07 Jan", time: "12:00 - 14:30"},
		{date: "Sun, 08 Jan", time: "20:00 - 22:30"},
		{date: "Tue, 09 Jan", time: "04:00 - 06:30"},
		{date: "Wed, 10 Jan", time: "10:00 - 12:30"},
		{date: "Thu, 11 Jan", time: "18:00 - 20:30"},
		{date: "Sat, 12 Jan", time: "02:00 - 04:30"},
		{date: "Sun, 13 Jan", time: "08:00 - 10:30"},
		{date: "Mon, 14 Jan", time: "16:00 - 18:30"},
		{date: "Wed, 15 Jan", time: "00:00 - 02:30"},
		{date: "Thu, 16 Jan", time: "06:00 - 08:30"},
	}

	now := time.Date(2021, 12, 25, 18, 00, 00, 0, loc)
	scheduleItems, err := makeScheduleItems(testItems, now)
	if err != nil {
		t.Errorf("unexpected error making schedule items: %v", err)
		return
	}

	expected := []ScheduleItem{
		{Start: time.Date(2021, 12, 26, 04, 00, 00, 0, loc), End: time.Date(2021, 12, 26, 06, 30, 00, 0, loc)},
		{Start: time.Date(2021, 12, 27, 12, 00, 00, 0, loc), End: time.Date(2021, 12, 27, 14, 30, 00, 0, loc)},
		{Start: time.Date(2021, 12, 28, 20, 00, 00, 0, loc), End: time.Date(2021, 12, 28, 22, 30, 00, 0, loc)},
		{Start: time.Date(2021, 12, 29, 18, 00, 00, 0, loc), End: time.Date(2021, 12, 29, 20, 30, 00, 0, loc)},
		{Start: time.Date(2021, 12, 30, 02, 00, 00, 0, loc), End: time.Date(2021, 12, 30, 04, 30, 00, 0, loc)},
		{Start: time.Date(2021, 12, 31, 10, 00, 00, 0, loc), End: time.Date(2021, 12, 31, 12, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 01, 16, 00, 00, 0, loc), End: time.Date(2022, 01, 01, 18, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 02, 00, 00, 00, 0, loc), End: time.Date(2022, 01, 02, 02, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 03, 8, 00, 00, 0, loc), End: time.Date(2022, 01, 03, 10, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 04, 14, 00, 00, 0, loc), End: time.Date(2022, 01, 04, 16, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 05, 22, 00, 00, 0, loc), End: time.Date(2022, 01, 05, 00, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 06, 06, 00, 00, 0, loc), End: time.Date(2022, 01, 06, 8, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 07, 12, 00, 00, 0, loc), End: time.Date(2022, 01, 07, 14, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 8, 20, 00, 00, 0, loc), End: time.Date(2022, 01, 8, 22, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 9, 04, 00, 00, 0, loc), End: time.Date(2022, 01, 9, 06, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 10, 10, 00, 00, 0, loc), End: time.Date(2022, 01, 10, 12, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 11, 18, 00, 00, 0, loc), End: time.Date(2022, 01, 11, 20, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 12, 02, 00, 00, 0, loc), End: time.Date(2022, 01, 12, 04, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 13, 8, 00, 00, 0, loc), End: time.Date(2022, 01, 13, 10, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 14, 16, 00, 00, 0, loc), End: time.Date(2022, 01, 14, 18, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 15, 00, 00, 00, 0, loc), End: time.Date(2022, 01, 15, 02, 30, 00, 0, loc)},
		{Start: time.Date(2022, 01, 16, 06, 00, 00, 0, loc), End: time.Date(2022, 01, 16, 8, 30, 00, 0, loc)},
	}

	for index, scheduleItem := range scheduleItems {
		if !expected[index].Start.Equal(scheduleItem.Start) {
			t.Errorf("expected item %d to start to be %s, got %s", index, expected[index].Start, scheduleItem.Start)
		}
		if !expected[index].End.Equal(scheduleItem.End) {
			t.Errorf("expected item %d to end to be %s, got %s", index, expected[index].End, scheduleItem.End)
		}
	}
}

func TestScheduleFromHTML(t *testing.T) {
	loc, err := time.LoadLocation("Africa/Johannesburg")
	if err != nil {
		t.Errorf("unexpected error loading tz data: %v", err)
		return
	}
	now := time.Date(2021, 12, 25, 18, 00, 00, 0, loc)
	testFile, err := os.Open("./test_data/schedule.html")
	if err != nil {
		t.Errorf("unexpected error opening test file: %v", err)
		return
	}

	testData, err := ioutil.ReadAll(testFile)
	if err != nil {
		t.Errorf("unexpected error reading test file: %v", err)
		return
	}

	res, err := scheduleFromHTML(testData, 1, now)
	if err != nil {
		t.Errorf("expected err to be nil, got: %v", err)
	}
	expected := Schedule{
		Stage: Stage(1),
		Times: []ScheduleItem{
			{Start: time.Date(2021, 10, 29, 04, 0, 0, 0, loc), End: time.Date(2021, 10, 29, 06, 30, 00, 0, loc)},
			{Start: time.Date(2021, 10, 30, 12, 0, 0, 0, loc), End: time.Date(2021, 10, 30, 14, 30, 00, 0, loc)},
			{Start: time.Date(2021, 10, 31, 20, 0, 0, 0, loc), End: time.Date(2021, 10, 31, 22, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 01, 18, 0, 0, 0, loc), End: time.Date(2021, 11, 01, 20, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 03, 02, 0, 0, 0, loc), End: time.Date(2021, 11, 03, 04, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 04, 10, 0, 0, 0, loc), End: time.Date(2021, 11, 04, 12, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 05, 16, 0, 0, 0, loc), End: time.Date(2021, 11, 05, 18, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 07, 00, 0, 0, 0, loc), End: time.Date(2021, 11, 07, 02, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 8, 8, 0, 0, 0, loc), End: time.Date(2021, 11, 8, 10, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 9, 14, 0, 0, 0, loc), End: time.Date(2021, 11, 9, 16, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 10, 22, 0, 0, 0, loc), End: time.Date(2021, 11, 10, 00, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 12, 06, 0, 0, 0, loc), End: time.Date(2021, 11, 12, 8, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 13, 12, 0, 0, 0, loc), End: time.Date(2021, 11, 13, 14, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 14, 20, 0, 0, 0, loc), End: time.Date(2021, 11, 14, 22, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 16, 04, 0, 0, 0, loc), End: time.Date(2021, 11, 16, 06, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 17, 10, 0, 0, 0, loc), End: time.Date(2021, 11, 17, 12, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 18, 18, 0, 0, 0, loc), End: time.Date(2021, 11, 18, 20, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 20, 02, 0, 0, 0, loc), End: time.Date(2021, 11, 20, 04, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 21, 8, 0, 0, 0, loc), End: time.Date(2021, 11, 21, 10, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 22, 16, 0, 0, 0, loc), End: time.Date(2021, 11, 22, 18, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 24, 00, 0, 0, 0, loc), End: time.Date(2021, 11, 24, 02, 30, 00, 0, loc)},
			{Start: time.Date(2021, 11, 25, 06, 0, 0, 0, loc), End: time.Date(2021, 11, 25, 8, 30, 00, 0, loc)},
		},
	}
	for index, scheduleItem := range res.Times {
		if !expected.Times[index].Start.Equal(scheduleItem.Start) {
			t.Errorf("expected item %d to start to be %s, got %s", index, expected.Times[index].Start, scheduleItem.Start)
		}
		if !expected.Times[index].End.Equal(scheduleItem.End) {
			t.Errorf("expected item %d to end to be %s, got %s", index, expected.Times[index].End, scheduleItem.End)
		}
	}
}
