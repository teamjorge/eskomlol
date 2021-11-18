package eskomlol

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client is the structure for performing requests to the Eskom API.
type Client struct {
	timeout    time.Duration
	httpClient HttpClient
	nowFunc    func() time.Time
}

// New creates an instance of the Client with the given options.
//
// The timeout is set by default to 30 seconds and can be overridden with the withTimeout option.
func New(opts ...ClientOpt) *Client {
	c := new(Client)
	c.timeout = 30 * time.Second
	c.nowFunc = time.Now
	c.httpClient = nil

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Status retrieves the current Loadshedding stage.
//
// Values of -1 and 0 indicate no loadshedding currently.
func (c *Client) Status(ctx context.Context) (Stage, error) {
	h := getClient(c)

	data, err := doRequest(ctx, h, "/GetStatus", nil)
	if err != nil {
		return -1, err
	}

	status, err := strconv.Atoi(string(data))
	if err != nil {
		return -1, err
	}

	return stageMap[status], nil
}

// Municipalities returns a list of municipalities that Eskom supplies to.
func (c *Client) Municipalities(ctx context.Context, province Province) (Municipalities, error) {
	h := getClient(c)

	requestURL := fmt.Sprintf("/GetMunicipalities/?Id=%d", province)
	var municipalities Municipalities

	err := doRequestJSON(ctx, h, requestURL, nil, &municipalities)

	return municipalities, err
}

// Suburbs returns a list of suburbs from the given municipality and search term.
//
// The responses are paginated and can be iterated by using the page parameter. The result
// object contains a Total field to indicate the total number of results.
func (c *Client) Suburbs(ctx context.Context, municipalityID string, searchTerm string, page int) (SuburbResult, error) {
	h := getClient(c)
	if page < 1 {
		page = 1
	}
	requestURL := fmt.Sprintf(
		"/GetSurburbData/?pageSize=100&pageNum=%d&searchTerm=%s&id=%s",
		page, url.QueryEscape(searchTerm), url.QueryEscape(municipalityID),
	)
	var suburbResult SuburbResult
	err := doRequestJSON(ctx, h, requestURL, nil, &suburbResult)

	return suburbResult, err
}

// SearchSuburbs returns all suburbs that match the given searchTerm.
//
// maxResults can be omitted (nil) if the default of 300 is acceptable.
func (c *Client) SearchSuburbs(ctx context.Context, searchTerm string, maxResults *int) (SearchSuburbs, error) {
	h := getClient(c)
	maxRes := 300
	if maxResults != nil {
		maxRes = *maxResults
	}

	requestURL := fmt.Sprintf(
		"/FindSuburbs?searchText=%s&maxResults=%d",
		url.QueryEscape(searchTerm), maxRes,
	)
	var searchSuburbs SearchSuburbs
	err := doRequestJSON(ctx, h, requestURL, nil, &searchSuburbs)

	return searchSuburbs, err
}

// Schedule returns the loadshedding schedule for the given suburb and stage(s).
func (c *Client) Schedule(ctx context.Context, suburbID string, stages ...Stage) (map[Stage]Schedule, error) {
	h := getClient(c)
	errs := make([]string, 0)
	res := make(map[Stage]Schedule)

	for _, stage := range stages {
		if !stage.Valid() {
			errs = append(errs, fmt.Sprintf("%d is not a valid stage", stage))
			continue
		}
		if stage < 1 {
			errs = append(errs, "only Stages 1 - 8 are valid for schedules")
			continue
		}
		requestURL := fmt.Sprintf(`/GetScheduleM/%s/%d/_/1`, suburbID, stage)
		data, err := doRequest(ctx, h, requestURL, nil)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		s, err := scheduleFromHTML(data, stage, c.nowFunc())
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		res[stage] = s
	}
	var err error
	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, "; "))
	}
	return res, err
}
