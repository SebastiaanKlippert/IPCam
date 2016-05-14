package ipcam

import "time"

type Snapshot struct {
	Buf      []byte
	DateTime time.Time
}

//SaveFile saves the snapshot to disk in folder under filename.jpg
//If filename is empty it uses the timestamp of the snapshot
func (s *Snapshot) SaveFile(folder, filename string) (string, error) {
	return "", nil
}
