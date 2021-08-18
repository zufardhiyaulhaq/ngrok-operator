package handler

import "net/http"

type HTTPStatusHandler struct {
}

func (h HTTPStatusHandler) Running(domain string) (bool, error) {
	response, err := http.Get(domain)
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
