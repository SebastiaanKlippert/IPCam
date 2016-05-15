package ipcam

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Snapshot struct {
	Buf      []byte
	DateTime time.Time
}

//SaveFile saves the snapshot to disk in folder under filename.jpg
//If folder is empty it uses the current dir
//If filename is empty it uses the timestamp of the snapshot
func (s *Snapshot) SaveFile(folder, filename string) (string, error) {

	if folder == "" {
		folder = filepath.Dir(os.Args[0])
	}
	if filename == "" {
		filename = s.DateTime.Format("2006-01-02T150405.000.jpg")
	}
	if filepath.Ext(filename) == "" {
		filename += ".jpg"
	}

	f := filepath.Join(folder, filename)

	return f, ioutil.WriteFile(f, s.Buf, os.ModeExclusive)
}
