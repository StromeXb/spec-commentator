package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"spec-commentor/pkg/utils"
)

func (cfg Config) formatURL(operation string) string {
	serviceURL := &url.URL{
		Scheme: cfg.Scheme,
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   cfg.URL,
	}

	return serviceURL.String() + operation
}

func (cfg Config) MakeRequestURLQuery(operation string, headers map[string]string, data map[string]interface{}) (body []byte, err error) {
	reqData := url.Values{}
	for k, i := range data {
		reqData.Set(k, fmt.Sprintf("%v", i))
	}

	client := &http.Client{}
	timeout := time.Second * TimeoutSecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	reqDataEncoded := reqData.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.formatURL(operation), strings.NewReader(reqDataEncoded))
	if err != nil {
		return nil, err
	}

	for header, value := range headers {
		r.Header.Add(header, value)
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if (resp.StatusCode >= http.StatusOK && resp.StatusCode <= 299) || resp.StatusCode == http.StatusBadRequest {
		if utils.IsJSON(string(body)) {
			fmt.Println("Response error")
		} else {
			fmt.Println("Response error")
		}

		return body, nil
	}
	if utils.IsJSON(string(body)) {
		fmt.Println("Response error")
	} else {
		fmt.Println("Response error")
	}

	return nil, NewError(fmt.Sprintf("%s response error", cfg.formatURL(operation)))
}
