package handler

import (
	"net/http"
	"time"
)

type HTTPStatusHandler struct {
}

func (h HTTPStatusHandler) Running(domain string) (bool, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Get(domain)
	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusPaymentRequired {
		return false, nil
	}

	return true, nil
}

func NewHHTPStatusHandler() StatusHandler {
	return HTTPStatusHandler{}
}
