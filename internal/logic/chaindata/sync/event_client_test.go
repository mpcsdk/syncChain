package syncBlock

import (
	"os"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	defer func() {
		os.RemoveAll("logs")
	}()

	// Init()

	time.Sleep(10 * time.Hour)
}
