package pkg

import (
	"net"
	"net/http"
	"time"
)

// KtCustomClient returns a custom http client for ktCloud API
// It has a timeout of 5 seconds and transport with a timeout of 3 seconds
func KtCustomClient() *http.Client {
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}

	return &client
}
