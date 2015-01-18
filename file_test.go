package mpq

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestFile_Files(t *testing.T) {
	setup()

	fileTest := []string{
		"(attributes)",
		"(listfile)",
		"replay.attributes.events",
		"replay.details",
		"replay.game.events",
		"replay.initData",
		"replay.load.info",
		"replay.message.events",
		"replay.resumable.events",
		"replay.server.battlelobby",
		"replay.smartcam.events",
		"replay.sync.events",
		"replay.sync.history",
		"replay.tracker.events",
	}

	files, err := m.Files()
	if err != nil {
		t.Fatal("Could not retrieve file list:", err)
	}

	fmt.Println(files)

	if len(fileTest) != len(files) {
		t.Fatal("Length of the files are different:", len(fileTest), len(files))
	}

	sort.Strings(fileTest)
	sort.Strings(files)

	for i, test := range fileTest {
		if test != files[i] {
			t.Error("Expected file:", test, "got:", files[i])
		}
	}
}

func TestFile_FileInfo(t *testing.T) {
	m = nil

	file, err := m.fileInfo("(listfile)")
	if err == nil || strings.Contains(err.Error(), "HET, BET, Hash and Block") {
		t.Error(`Expected an error about "HET, BET, Hash and Block"`)
	}

	setup()

	file, err = m.fileInfo("(listfile)")
	if file == nil || err != nil {
		t.Error("Expected HET and BET lookup to succeed.")
	}

	m.BETTable = nil
	m.HETTable = nil

	file, err = m.fileInfo("(listfile)")
	if file == nil || err != nil {
		t.Error("Expected Hash and Block lookup to succeed.")
	}

	m.Close()
	m = nil
}

func TestFile_FromHETAndBET(t *testing.T) {
	setup()

	file, err := m.findFromHETAndBET("(listfile)")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if file.Name != "(listfile)" {
		t.Errorf("Wrong File Name: %s", file.Name)
	}
	if file.FileSize != 0x104 {
		t.Errorf("Wrong File Size: % 02X", file.FileSize)
	}
	if file.CompressedSize != 0x97 {
		t.Error("Wrong File Compressed Size: % 02X", file.CompressedSize)
	}
	if file.Position != 0x6B991 {
		t.Errorf("Wrong Position: % 02X", file.Position)
	}
	if file.Flags != 0x81000200 {
		t.Errorf("Wrong File Flags: % 02X", file.Flags)
	}
}

func TestFile_FromHashAndBlock(t *testing.T) {
	setup()

	file, err := m.findFromHashAndBlock("(listfile)")
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if file.Name != "(listfile)" {
		t.Errorf("Wrong File Name: %s", file.Name)
	}
	if file.Locale != LocaleNeutral {
		t.Error("Wrong Locale: % 02X", file.Locale)
	}
	if file.FileSize != 0x104 {
		t.Errorf("Wrong File Size: % 02X", file.FileSize)
	}
	if file.CompressedSize != 0x97 {
		t.Error("Wrong File Compressed Size: % 02X", file.CompressedSize)
	}
	if file.Position != 0x6B991 {
		t.Errorf("Wrong File Position: % 02X", file.Position)
	}
	if file.Flags != 0x81000200 {
		t.Errorf("Wrong File Flags: % 02X", file.Flags)
	}
}
