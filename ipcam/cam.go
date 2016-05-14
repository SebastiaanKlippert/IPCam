package ipcam

import (
	"fmt"
	"time"
)

//NewCams returns the Cams object for a config file
func NewCams(configFile string) (*Cams, error) {

	cfg, err := readConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading config file %s: %s", configFile, err)
	}

	cams := new(Cams)
	cams.Config = cfg.globalconfig

	for _, c := range cfg.cams {

		cam := Cam{
			config: c,
		}

		cams.Cams = append(cams.Cams, cam)
	}

	return cams, nil
}

//Cams contains all IP cameras, one Cam objct for each camera in config
type Cams struct {
	Config globalconfig
	Cams   []Cam
}

//CamByName returns a Cam object by its name
func (c *Cams) CamByName(name string) *Cam {
	for _, cam := range c.Cams {
		if cam.Name() == name {
			return &cam
		}
	}
	return nil
}

//Cam is a single IP camera
type Cam struct {
	config camconfig
}

//Name returns the name of this cam
func (c *Cam) Name() string {
	return c.config.Name
}

//TakeSnapshot takes a single snapshot
func (c *Cam) TakeSnapshot() (*Snapshot, error) {
	return nil, nil
}

//TakeSnapshots takes one snapshot every interval and saves them in folder
func (c *Cam) TakeSnapshots(interval time.Duration, folder string) (*Snapshot, error) {
	return nil, nil
}
