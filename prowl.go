package prowl

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	serverURL  = "https://api.prowlapp.com"
	addPath    = "/publicapi/add"
	verifyPath = "/publicapi/verify"
	addURL     = serverURL + addPath
	verifyURL  = serverURL + verifyPath
)

// Service represents the prowl service
type Service struct {
	APIKey string
}

// Notification represents a prowl message
type Notification struct {
	AppName   string
	EventName string
	Message   string
	Priority  string
	URL       string
	// Provider string // Not gonna get one of those keys.
}

// PublishMsg publishes a simplified message
func (p *Service) PublishMsg(message string) (err error) {
	notification := Notification{
		AppName:   "Go Prowl",
		EventName: "Message",
		Message:   message,
	}

	return p.Publish(notification)
}

// Publish a notification to prowl
func (p *Service) Publish(notification Notification) (err error) {
	formParams := url.Values{
		"apikey":      []string{p.APIKey},
		"application": []string{notification.AppName},
		"description": []string{notification.Message},
		"event":       []string{notification.EventName},
		//		"priority":    []string{n.Priority},
	}

	if notification.URL != "" {
		formParams["url"] = []string{notification.URL}
	} else if strings.HasPrefix(notification.Message, "http") && strings.Contains(notification.Message, " ") == false {
		formParams["url"] = []string{notification.Message}
	}

	response, err := http.PostForm(addURL, formParams)

	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		err = decodeError(response.Status, response.Body)
		// } else {
		//	decodeResponse(response.Status, response.Body)
	}
	return
}

func decodeResponse(status string, bodyReader io.Reader) string {
	xmlResponse := successResponse{}
	if xml.NewDecoder(bodyReader).Decode(&xmlResponse) != nil {
		return status
	}
	return string(xmlResponse.Succcess.Remaining)
}

func decodeError(status string, bodyReader io.Reader) (err error) {
	xmlResponse := errorResponse{}
	if xml.NewDecoder(bodyReader).Decode(&xmlResponse) != nil {
		err = errors.New(status)
	} else {
		err = errors.New(xmlResponse.Error.Message)
	}
	return
}

// GG. Must be a way to do this in a dynamic and more elegant manner.
type errorResponse struct {
	Error struct {
		Code    int    `xml:"code,attr"`
		Message string `xml:",chardata"`
	} `xml:"error"`
}

type successResponse struct {
	Succcess struct {
		Code      int `xml:"code,attr"`
		ResetDate int `xml:"resetdate,attr"`
		Remaining int `xml:"remaining,attr"`
	} `xml:"success"`
}
