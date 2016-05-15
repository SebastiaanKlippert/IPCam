package ipcam

import (
	"fmt"
	"log"
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

	for _, c := range cfg.Cams {

		cam := Cam{
			cconfig: &c,
			gconfig: &cfg.globalconfig,
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
	gconfig *globalconfig
	cconfig *camconfig
}

//Name returns the name of this cam
func (c *Cam) Name() string {
	return c.cconfig.Name
}

//TakeSnapshot takes a single snapshot
func (c *Cam) TakeSnapshot() (*Snapshot, error) {
	buf, err := snapshot(c.gconfig, c.cconfig)
	if err != nil {
		return nil, err
	}
	s := &Snapshot{
		Buf:      buf,
		DateTime: time.Now().Local(),
	}
	return s, nil
}

//TakeSnapshots takes one snapshot every interval for duration and saves them in folder
func (c *Cam) TakeSnapshots(interval, duration time.Duration, folder string) error {

	if duration.Nanoseconds() <= interval.Nanoseconds() {
		return fmt.Errorf("duration must be longer than interval")
	}

	tick := time.NewTicker(interval)
	timer := time.NewTimer(duration)

	for {
		select {
		case <-tick.C:
			go func() {
				s, err := c.TakeSnapshot()
				if err != nil {
					log.Println(err)
				}
				_, err = s.SaveFile(folder, "")
				if err != nil {
					log.Println(err)
				}
			}()
		case <-timer.C:
			tick.Stop()
			return nil
		}
	}

	return nil
}
