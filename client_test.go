package eskomlol

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

type clientMockHTTPClient struct {
	StatusResponse         []byte
	MunicipalitiesResponse []byte
	SuburbsResponse        []byte
	SearchSuburbsResponse  []byte
	ScheduleResponse       []byte
}

func (m *clientMockHTTPClient) mapResponses(url string) []byte {
	return map[string][]byte{
		baseURL + "/GetStatus":               m.StatusResponse,
		baseURL + "/GetMunicipalities/?Id=1": m.MunicipalitiesResponse,
		baseURL + "/GetSurburbData/?pageSize=100&pageNum=1&searchTerm=bryanston&id=1": m.SuburbsResponse,
		baseURL + "/FindSuburbs?searchText=bryanston&maxResults=300":                  m.SearchSuburbsResponse,
		baseURL + "/GetScheduleM/1/1/_/1":                                             m.ScheduleResponse,
	}[url]
}

func (m *clientMockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	responseData := m.mapResponses(req.URL.String())
	res := http.Response{}

	buf := bytes.NewBuffer(responseData)

	res.Body = io.NopCloser(buf)

	return &res, nil
}

func TestNewClient(t *testing.T) {
	c := New()
	if c.timeout != time.Duration(30*time.Second) {
		t.Error("expected default timeout to be 30 seconds")
	}
	nowFuncDate := c.nowFunc()
	if !nowFuncDate.After(time.Time{}) {
		t.Errorf("expected default nowFunc to return a date after nil value. got: %v", nowFuncDate)
	}

	// With Opts
	c = New(WithTimeout(40 * time.Second))
	if c.timeout != time.Duration(40*time.Second) {
		t.Error("expected configured timeout to be 40 seconds")
	}
}

func TestStatus(t *testing.T) {
	c := New(withHTTPClient(&clientMockHTTPClient{
		StatusResponse: []byte("2"),
	}))
	stage, err := c.Status(context.Background())
	if err != nil {
		t.Errorf("did not expect an error when calling Status, got: %v", err)
	}

	if stage != 1 {
		t.Errorf("expected stage value to be 1, got %d", stage)
	}
}

func TestStatusError(t *testing.T) {
	c := New(withHTTPClient(&clientMockHTTPClient{
		StatusResponse: []byte("asd"),
	}))
	stage, err := c.Status(context.Background())
	if err == nil {
		t.Error("expected the error to not be nil")
	}

	if stage != -1 {
		t.Errorf("expected stage to be -1, got %d", stage)
	}

	expectedErr := "strconv.Atoi: parsing \"asd\": invalid syntax"

	if err.Error() != expectedErr {
		t.Errorf("expected err value to be %s, got %s", expectedErr, err.Error())
	}
}

func TestMunicipalities(t *testing.T) {
	c := New(withHTTPClient(&clientMockHTTPClient{
		MunicipalitiesResponse: []byte(`
		[
			{
				"Value": "1",
				"Text": "Blah",
				"Disabled": false,
				"Selected": false,
				"Group": "1"
			}
		]
		`),
	}))

	municipalities, err := c.Municipalities(context.Background(), 1)
	if err != nil {
		t.Errorf("did not expect an error when calling Municipalities, got: %v", err)
	}

	if len(municipalities) != 1 {
		t.Error("expected only 1 municipality")
	}

	if municipalities[0].ID != "1" {
		t.Errorf("expected id to be 1, got %s", municipalities[0].ID)
	}
	if municipalities[0].Name != "Blah" {
		t.Errorf("expected id to be Blah, got %s", municipalities[0].Name)
	}
}

func TestSuburbs(t *testing.T) {
	c := New(withHTTPClient(&clientMockHTTPClient{
		SuburbsResponse: []byte(`
		{
			"Results": [
				{
					"id": "1",
					"text": "Bryanston1",
					"Tot": 4
				},
				{
					"id": "2",
					"text": "Bryanston2",
					"Tot": 6
				}
			],
			"Total": 2
		}
		`),
	}))

	suburbResult, err := c.Suburbs(context.Background(), "1", "bryanston", 1)
	if err != nil {
		t.Errorf("did not expect an error when calling Suburbs, got: %v", err)
	}

	if len(suburbResult.Results) != 2 {
		t.Error("expected only 2 suburbs")
	}

	if suburbResult.Results[0].ID != "1" {
		t.Errorf("expected the first suburb to have ID 1, got %s", suburbResult.Results[0].ID)
	}

	if suburbResult.Results[1].Name != "Bryanston2" {
		t.Errorf("expected the second suburb to have name Bryanston2, got %s", suburbResult.Results[1].Name)
	}
}

func TestSearchSuburbs(t *testing.T) {
	c := New(withHTTPClient(&clientMockHTTPClient{
		SearchSuburbsResponse: []byte(`
		[
			{
				"MunicipalityName": "City Power",
				"ProvinceName": "Gauteng",
				"Name": "Bryanston1",
				"ID": 992,
				"Total": 555
			},
			{
				"MunicipalityName": "City Power",
				"ProvinceName": "Gauteng",
				"Name": "Bryanston2",
				"ID": 993,
				"Total": 556
			}
		]
		`),
	}))

	suburbs, err := c.SearchSuburbs(context.Background(), "bryanston", nil)
	if err != nil {
		t.Errorf("did not expect an error when calling SearchSuburbs, got: %v", err)
	}

	if len(suburbs) != 2 {
		t.Error("expected only 2 suburbs")
	}

	if suburbs[0].ID != 992 {
		t.Errorf("expected the first suburb to have ID 992, got %d", suburbs[0].ID)
	}
	if suburbs[1].Name != "Bryanston2" {
		t.Errorf("expected the second suburb to have name Bryanston2, got %s", suburbs[1].Name)
	}
	if suburbs[0].ProvinceName != "Gauteng" {
		t.Errorf("expected first suburb to have Province Gauteng, got %s", suburbs[0].ProvinceName)
	}
	if suburbs[1].MunicipalityName != "City Power" {
		t.Errorf("expected second suburb to have MunicipalityName 'City Power', got %s", suburbs[1].MunicipalityName)
	}
}

func TestSchedule(t *testing.T) {
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
	loc, err := time.LoadLocation("Africa/Johannesburg")
	if err != nil {
		t.Errorf("unexpected error loading tz data: %v", err)
		return
	}

	c := New(withHTTPClient(&clientMockHTTPClient{
		ScheduleResponse: []byte(testData),
	}), withNowFunc(func() time.Time {
		return time.Date(2021, 10, 27, 18, 00, 00, 0, loc)
	}))

	schedule, err := c.Schedule(context.Background(), "1", 1)
	if err != nil {
		t.Errorf("did not expect an error when calling Schedule, got: %v", err)
	}

	expectedResult := map[Stage]Schedule{
		1: {
			Stage: 1,
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
		},
	}

	if len(schedule[1].Times) != 22 {
		t.Errorf("expected 22 schedule times, only got %d", len(schedule[1].Times))
	}

	for stage, val := range schedule {
		for index, scheduleItem := range val.Times {
			if !expectedResult[stage].Times[index].Start.Equal(scheduleItem.Start) {
				t.Errorf("expected item %d to start to be %s, got %s", index, expectedResult[stage].Times[index].Start, scheduleItem.Start)
			}
			if !expectedResult[stage].Times[index].End.Equal(scheduleItem.End) {
				t.Errorf("expected item %d to end to be %s, got %s", index, expectedResult[stage].Times[index].End, scheduleItem.End)
			}
		}
	}
}
