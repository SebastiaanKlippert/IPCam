package ipcam

import (
	"fmt"
	"log"
	"sync"
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

//TakeSnapshots takes one snapshot every interval for duration and adds them to the return channel
func (c *Cam) TakeSnapshots(interval, duration time.Duration) (chan *Snapshot, error) {

	if duration.Nanoseconds() <= interval.Nanoseconds() {
		return nil, fmt.Errorf("duration must be longer than interval")
	}

	ch := make(chan *Snapshot, 100)

	go func() {

		tick := time.NewTicker(interval)
		end := time.Now().Add(duration)
		wg := new(sync.WaitGroup)

		for time.Now().Before(end) {
			select {
			case <-tick.C:
				wg.Add(1)
				go func() {

					defer wg.Done()

					s, err := c.TakeSnapshot()
					if err != nil {
						log.Println(err)
						return
					}
					ch <- s
				}()
			}
		}

		//Stop ticker, wait until last snapshot is written and then close channel
		tick.Stop()
		wg.Wait()
		close(ch)

	}()

	return ch, nil
}

//TakeSnapshots takes one snapshot every interval for duration and saves them in folder, save errors are logged
func (c *Cam) TakeAndSaveSnapshots(interval, duration time.Duration, folder string) error {

	ch, err := c.TakeSnapshots(interval, duration)
	if err != nil {
		return err
	}

	for snap := range ch {
		_, err := snap.SaveFile(folder, "")
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

//Stream opens a video stream which is written to w.
//Data written is a multipart JFIF stream with --ipcamera borders
func (c *Cam) Stream(w *HttpFlushWriter) error {
	return stream(w, c.gconfig, c.cconfig)
}
