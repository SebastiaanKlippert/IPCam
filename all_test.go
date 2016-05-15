package ipcam

import (
	"log"
	"os"
	"testing"
	"time"
)

const (
	//TODO
	testPath      = "C:/github/IPCam/"
	testConfig    = "config.json"
	testCamName   = "Woonkamer"
	testPathSnaps = testPath + "snaps"
)

var (
	c   *Cams
	cam *Cam
)

func init() {
	err := os.MkdirAll(testPathSnaps, os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}
}

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

	c, err = NewCams(testPath + testConfig)
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

	if len(snap.Buf) == 0 {
		t.Fatalf("Empty snapshot data")
	}
}

func TestSnapShotSave(t *testing.T) {
	checkEnd(t, true)

	snap, err := cam.TakeSnapshot()
	if err != nil {
		t.Fatal(err)
	}

	fname, err := snap.SaveFile(testPathSnaps, "")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Snapshot saved to %s", fname)
}

func TestTakeAndSaveSnapshots(t *testing.T) {
	checkEnd(t, true)

	t.Log("Taking snapshots every 0.5 seconds for 5 seconds")

	err := cam.TakeAndSaveSnapshots(500*time.Millisecond, 5*time.Second, testPathSnaps)
	if err != nil {
		t.Fatal(err)
	}
}
