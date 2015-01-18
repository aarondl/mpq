package mpq

import "testing"

func TestBlockTable(t *testing.T) {
	setup()

	testEntries := []BlockTableEntry{
		{0x000004D0, 0x0000035B, 0x000004C1, 0x81000200},
		{0x0000083B, 0x00000682, 0x00000D01, 0x81000200},
		{0x00000ECD, 0x0000466A, 0x00009980, 0x81000200},
		{0x00005557, 0x00056DD2, 0x000DDCA6, 0x81000200},
		{0x0005C489, 0x0000054E, 0x0000073B, 0x81000200},
		{0x0005C9E7, 0x00000028, 0x00000038, 0x81000200},
		{0x0005CA1F, 0x000004A2, 0x000006B3, 0x81000200},
		{0x00000000, 0x00000000, 0x00000000, 0x80000200},
		{0x0005CED1, 0x0000E0DA, 0x00086BF7, 0x81000200},
		{0x0006AFEB, 0x00000384, 0x00000596, 0x81000200},
		{0x0006B37F, 0x00000462, 0x000013B0, 0x81000200},
		{0x0006B7F1, 0x00000190, 0x00000335, 0x81000200},
		{0x0006B991, 0x00000097, 0x00000104, 0x81000200},
		{0x0006BA38, 0x00000120, 0x00000120, 0x81000200},
	}

	entries := m.BlockTable.Entries()
	if len(entries) != len(testEntries) {
		t.Error("Size mismatch:", len(entries))
	}

	for i, test := range testEntries {
		if entries[i].FilePosition != test.FilePosition {
			t.Errorf("%d> FilePos wrong: %02X", i, entries[i].FilePosition)
		}
		if entries[i].CompressedSize != test.CompressedSize {
			t.Errorf("%d> CompressedSize wrong: %02X", i, entries[i].CompressedSize)
		}
		if entries[i].FileSize != test.FileSize {
			t.Errorf("%d> FileSize wrong: %02X", i, entries[i].FileSize)
		}
		if entries[i].Flags != test.Flags {
			t.Errorf("%d> Flags wrong: %02X", i, entries[i].Flags)
		}
	}

}
