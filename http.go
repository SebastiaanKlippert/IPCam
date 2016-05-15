package ipcam

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
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

func stream(fw *HttpFlushWriter, globalcfg *globalconfig, camcfg *camconfig) error {

	uv := make(url.Values)
	uv.Add("user", camcfg.User)
	uv.Add("pwd", camcfg.Pass)

	u := url.URL{
		Scheme:   "http",
		Host:     fmt.Sprintf("%s:%s", camcfg.Host, camcfg.Port),
		Path:     "videostream.cgi",
		RawQuery: uv.Encode(),
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	mr := multipart.NewReader(resp.Body, "ipcamera")

	mw := multipart.NewWriter(fw.w)
	defer mw.Close()

	mw.SetBoundary("ipcamera")

	for {
		part, err := mr.NextPart()
		if err != nil {
			return err
		}

		tw, err := mw.CreatePart(part.Header)
		if err != nil {
			return err
		}

		_, err = io.Copy(tw, part)
		if err != nil {
			return err
		}

		fw.f.Flush()
	}

	return nil

}

type HttpFlushWriter struct {
	f http.Flusher
	w http.ResponseWriter
}

func (fw *HttpFlushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}

func HttpRwToFw(w http.ResponseWriter) *HttpFlushWriter {
	fw := HttpFlushWriter{w: w}
	if f, ok := w.(http.Flusher); ok {
		fw.f = f
	}
	return &fw
}
