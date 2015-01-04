package mpq

import "testing"

func TestBETTable(t *testing.T) {
	setup()

	if m.BETTable.Version != 1 {
		t.Errorf("Incorrect Value for Version: %d", m.BETTable.Version)
	}
	if m.BETTable.DataSize != 0x11F {
		t.Errorf("Incorrect Value for DataSize: %X", m.BETTable.DataSize)
	}
	if m.BETTable.TableSize != 0x11F {
		t.Errorf("Incorrect Value for TableSize: %X", m.BETTable.TableSize)
	}
	if m.BETTable.EntryCount != 14 {
		t.Errorf("Incorrect Value for EntryCount: %d", m.BETTable.EntryCount)
	}
	if m.BETTable.Unknown08 != 0x10 {
		t.Errorf("Incorrect Value for Unknown08: %X", m.BETTable.Unknown08)
	}
	if m.BETTable.TableEntrySize != 60 {
		t.Errorf("Incorrect Value for TableEntrySize: %d", m.BETTable.TableEntrySize)
	}
	if m.BETTable.BitIndexFilePos != 0 {
		t.Errorf("Incorrect Value for BitIndexFilePos: %d", m.BETTable.BitIndexFilePos)
	}
	if m.BETTable.BitIndexFileSize != 19 {
		t.Errorf("Incorrect Value for BitIndexFileSize: %d", m.BETTable.BitIndexFileSize)
	}
	if m.BETTable.BitIndexCmpSize != 39 {
		t.Errorf("Incorrect Value for BitIndexCmpSize: %d", m.BETTable.BitIndexCmpSize)
	}
	if m.BETTable.BitIndexFlagIndex != 58 {
		t.Errorf("Incorrect Value for BitIndexFlagIndex: %d", m.BETTable.BitIndexFlagIndex)
	}
	if m.BETTable.BitIndexUnknown != 0x3C {
		t.Errorf("Incorrect Value for BitIndexUnknown: %X", m.BETTable.BitIndexUnknown)
	}
	if m.BETTable.BitCountFilePos != 19 {
		t.Errorf("Incorrect Value for BitCountFilePos: %d", m.BETTable.BitCountFilePos)
	}
	if m.BETTable.BitCountFileSize != 20 {
		t.Errorf("Incorrect Value for BitCountFileSize: %d", m.BETTable.BitCountFileSize)
	}
	if m.BETTable.BitCountCmpSize != 19 {
		t.Errorf("Incorrect Value for BitCountCmpSize: %d", m.BETTable.BitCountCmpSize)
	}
	if m.BETTable.BitCountFlagIndex != 2 {
		t.Errorf("Incorrect Value for BitCountFlagIndex: %d", m.BETTable.BitCountFlagIndex)
	}
	if m.BETTable.BitCountUnknown != 0 {
		t.Errorf("Incorrect Value for BitCountUnknown: %d", m.BETTable.BitCountUnknown)
	}
	if m.BETTable.HashSizeTotal != 56 {
		t.Errorf("Incorrect Value for HashSizeTotal: %d", m.BETTable.HashSizeTotal)
	}
	if m.BETTable.HashSizeExtra != 0 {
		t.Errorf("Incorrect Value for HashSizeExtra: %d", m.BETTable.HashSizeExtra)
	}
	if m.BETTable.HashSize != 56 {
		t.Errorf("Incorrect Value for HashSize: %d", m.BETTable.HashSize)
	}
	if m.BETTable.HashArraySize != 98 {
		t.Errorf("Incorrect Value for HashArraySize: %d", m.BETTable.HashArraySize)
	}
	if m.BETTable.FlagCount != 2 {
		t.Errorf("Incorrect Value for FlagCount: %d", m.BETTable.FlagCount)
	}

	flags := []uint32{0x81000200, 0x80000200}
	for i := 0; i < len(flags); i++ {
		if flags[i] != m.BETTable.Flags[i] {
			t.Errorf("Incorrect Value for Flags: % 02X", m.BETTable.Flags[i])
		}
	}

	entryTests := []BETTableEntry{
		{0x00EE9DCAA075CAEC, 0x004D0, 0x004C1, 0x0035B, 0, 0x81000200},
		{0x00330B43A504AA7D, 0x0083B, 0x00D01, 0x00682, 0, 0x81000200},
		{0x00596196BFAC3121, 0x00ECD, 0x09980, 0x0466A, 0, 0x81000200},
		{0x0072D81E4CC62E5A, 0x05557, 0xDDCA6, 0x56DD2, 0, 0x81000200},
		{0x0045290E3E8FF05C, 0x5C489, 0x0073B, 0x0054E, 0, 0x81000200},
		{0x00FDEF2A65E4FCF2, 0x5C9E7, 0x00038, 0x00028, 0, 0x81000200},
		{0x00EAC8D601B705B5, 0x5CA1F, 0x006B3, 0x004A2, 0, 0x81000200},
		{0x004BED954CB26449, 0x00000, 0x00000, 0x00000, 1, 0x80000200},
		{0x002EEDABFC07F3B2, 0x5CED1, 0x86BF7, 0x0E0DA, 0, 0x81000200},
		{0x003C593F942FDD7A, 0x6AFEB, 0x00596, 0x00384, 0, 0x81000200},
		{0x0073D410423F4F5D, 0x6B37F, 0x013B0, 0x00462, 0, 0x81000200},
		{0x00BB513A4C07D865, 0x6B7F1, 0x00335, 0x00190, 0, 0x81000200},
		{0x00C43BB0B2F3866A, 0x6B991, 0x00104, 0x00097, 0, 0x81000200},
		{0x0055E9E414C107A7, 0x6BA38, 0x00120, 0x00120, 0, 0x81000200},
	}

	entries, err := m.BETTable.Entries()
	if err != nil {
		t.Error("Error loading table entries:", err)
	}

	if len(entryTests) != len(entries) {
		t.Error("Number of entries is wrong:", len(entries))
	}

	for i, test := range entryTests {
		if entries[i].NameHash2 != test.NameHash2 {
			t.Errorf("%d> NameHash wrong: %02X", i, entries[i].FilePosition)
		}
		if entries[i].FilePosition != test.FilePosition {
			t.Errorf("%d> File Position wrong: %02X", i, entries[i].FilePosition)
		}
		if entries[i].FileSize != test.FileSize {
			t.Errorf("%d> File Size wrong: %02X", i, entries[i].FilePosition)
		}
		if entries[i].CompressedSize != test.CompressedSize {
			t.Errorf("%d> Compressed size wrong: %02X", i, entries[i].CompressedSize)
		}
		if entries[i].FlagIndex != test.FlagIndex {
			t.Errorf("%d> Flag Index wrong: %02X", i, entries[i].FlagIndex)
		}
		if entries[i].Flags != test.Flags {
			t.Errorf("%d> Flags wrong: %02X", i, entries[i].Flags)
		}
	}
}
