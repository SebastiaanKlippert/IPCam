package ipcam

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func snapshot(globalcfg *globalconfig, camcfg *camconfig) ([]byte, error) {

	client := http.Client{
		Timeout:   time.Duration(globalcfg.Timeout) * time.Second,
		Transport: http.DefaultTransport,
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

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//always read body
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Statuscode %d, response:\n%s\n", resp.StatusCode, string(buf))
	}

	return buf, nil
}
