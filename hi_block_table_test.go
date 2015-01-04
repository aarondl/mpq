package mpq

import "testing"

func TestHiBlockTable(t *testing.T) {
	setup()

	if m.HiBlockTable != nil {
		t.Error("There should be no HiBlockTable.")
	}
}
