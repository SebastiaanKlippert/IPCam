package ipcam

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func snapshot(globalcfg *globalconfig, camcfg *camconfig) ([]byte, error) {

	client := http.Client{
		Timeout: time.Duration(globalcfg.Timeout) * time.Second,
	}

	uv := make(url.Values)
	uv.Add("user", camcfg.User)
	uv.Add("pwd", camcfg.Pass)

	u := url.URL{
		Scheme:   "http",
		Host:     fmt.Sprintf("%s:%s", camcfg.Host, camcfg.Port),
		Path:     "snapshot.cgi",
		RawQuery: uv.Encode(),
	}

	fmt.Println(u.String())

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(req.Method) //TODO

	client.Do(req)

	return nil, nil
}
