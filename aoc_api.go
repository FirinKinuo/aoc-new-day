package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

var aocUrl = url.URL{
	Scheme: "https",
	Host:   "adventofcode.com",
}

type AOCDay struct {
	Title        string
	Description  string
	TestInput    string
	ProblemInput string
}

type AOC struct {
	client *http.Client
}

func NewAOCApi(session string) *AOC {
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(&aocUrl, []*http.Cookie{
		{
			Name:  "session",
			Value: session,
		},
	},
	)
	client := &http.Client{
		Jar:     jar,
		Timeout: 5 * time.Second,
	}
	return &AOC{client: client}
}

func (a *AOC) GetDayInfo(day, year int) (AOCDay, error) {
	responseDayData, err := a.requestDayData(day, year)
	if err != nil {
		return AOCDay{}, fmt.Errorf("request for day data: %w", err)
	}

	if responseDayData.StatusCode != http.StatusOK {
		return AOCDay{}, fmt.Errorf("request for day data: status %d", responseDayData.StatusCode)

	}

	aocDay, err := a.parseDayDataResponse(responseDayData)
	if err != nil {
		return AOCDay{}, fmt.Errorf("parse day data response: %w", err)
	}

	responseProblemInput, err := a.requestProblemInput(day, year)
	if err != nil {
		return AOCDay{}, fmt.Errorf("request for day problem input: %w", err)
	}

	switch responseProblemInput.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
		log.Warnf(
			"Request for day problem input: status %d: empty or expired session",
			responseProblemInput.StatusCode,
		)
	default:
		return AOCDay{}, fmt.Errorf("request for day problem input: status %d", responseDayData.StatusCode)
	}

	aocDay.ProblemInput, err = a.readInputProblemResponse(responseProblemInput)
	if err != nil {
		return AOCDay{}, fmt.Errorf("read problem input: %w", err)
	}

	return aocDay, nil
}

func (a *AOC) requestDayData(day, year int) (*http.Response, error) {
	requestPath := fmt.Sprintf("%s/%d/day/%d", aocUrl.String(), year, day)

	response, err := a.client.Get(requestPath)
	if err != nil {
		return nil, fmt.Errorf("do request to advent of code: %w", err)
	}

	return response, nil
}

func (a *AOC) requestProblemInput(day, year int) (*http.Response, error) {
	requestPath := fmt.Sprintf("%s/%d/day/%d/input", aocUrl.String(), year, day)

	response, err := a.client.Get(requestPath)
	if err != nil {
		return nil, fmt.Errorf("do request to advent of code: %w", err)
	}

	return response, nil
}

func (a *AOC) parseDayDataResponse(response *http.Response) (AOCDay, error) {
	var day AOCDay

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return AOCDay{}, fmt.Errorf("new document from response: %w", err)
	}

	doc.Find(".day-desc").Each(a.parseDayDescription(&day))

	day.TestInput = doc.Find("pre code").First().Text()

	return day, nil
}

func (a *AOC) parseDayDescription(day *AOCDay) func(int, *goquery.Selection) {
	return func(_ int, selection *goquery.Selection) {
		description, _ := selection.Html()
		day.Description += description

		title := selection.Find("h2").Text()

		match := regexp.MustCompile(`---\s*Day\s*\d+:\s*(.*?)\s*---`).FindStringSubmatch(title)
		if len(match) > 1 {
			day.Title = match[1]
		}
	}
}

func (a *AOC) readInputProblemResponse(response *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", fmt.Errorf("new document from response: %w", err)
	}

	return doc.Text(), nil
}
