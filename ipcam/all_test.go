package ipcam

import "testing"

const (
	//TODO
	testConfig  = "C:/github/IPCam/config.json"
	testCamName = "Woonkamer"
)

var (
	c   *Cams
	cam *Cam
)

func checkEnd(t *testing.T, checkcam bool) {
	if c == nil {
		t.FailNow()
	}
	if checkcam && cam == nil {
		t.FailNow()
	}
}

func TestConfig(t *testing.T) {

	var err error

	c, err = NewCams(testConfig)
	if err != nil {
		t.Fatal(err)
	}

}

func TestCamByName(t *testing.T) {
	checkEnd(t, false)

	cam = c.CamByName(testCamName)
	if cam == nil {
		t.Fatalf("Cannot find cam %q", testCamName)
	}

}

func TestSnapShot(t *testing.T) {
	checkEnd(t, true)

	snap, err := cam.TakeSnapshot()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(snap) //TODO
}
