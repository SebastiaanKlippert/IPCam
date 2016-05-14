package ipcam

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type config struct {
	globalconfig
	cams []camconfig
}

type globalconfig struct {
	Timeout int
	start   time.Time
}

type camconfig struct {
	Name string
	Host string
	Port string
	User string
	Pass string
}

func readConfig(configFile string) (*config, error) {

	if filepath.Dir(configFile) == "." {
		configFile = filepath.Join(filepath.Dir(os.Args[0]), configFile)
	}

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	c := new(config)
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	c.start = time.Now()
	return c, nil
}
