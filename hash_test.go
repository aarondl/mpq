package mpq

import "testing"

func TestBlizz(t *testing.T) {
	if hash := blizz(`arr\units.dat`, 0); hash != 0xF4E6C69D {
		t.Error("Wrong value:", hash)
	}

	if hash := blizz(`unit\neutral\acritter.grp`, 0); hash != 0xA26067F3 {
		t.Error("Wrong value:", hash)
	}
}

func TestJenkins(t *testing.T) {
	tests := []struct {
		InitialPrimary   uint32
		InitialSecondary uint32
		ToHash           string
		Primary          uint32
		Secondary        uint32
	}{
		{0, 0, "", 0xdeadbeef, 0xdeadbeef},
		{0, 0xdeadbeef, "", 0xbd5b7dde, 0xdeadbeef},
		{0xdeadbeef, 0xdeadbeef, "", 0x9c093ccd, 0xbd5b7dde},
		{0, 0, "Four score and seven years ago", 0x17770551, 0xce7226e6},
		{0, 1, "Four score and seven years ago", 0xe3607cae, 0xbd371de4},
		{1, 0, "Four score and seven years ago", 0xcd628161, 0x6cbea4b3},
	}

	var pc, pb = new(uint32), new(uint32)
	for i, test := range tests {
		*pc = test.InitialPrimary
		*pb = test.InitialSecondary
		hashLittle2Str(test.ToHash, pc, pb)
		if *pc != test.Primary {
			t.Errorf("%d> Wrong Primary Value:   %08X", i, *pc)
		}
		if *pb != test.Secondary {
			t.Errorf("%d> Wrong Secondary Value: %08X", i, *pb)
		}
	}
}

func BenchmarkJenkins(b *testing.B) {
	strs := [][]byte{
		[]byte("Hello"),
		[]byte("Hello World"),
		[]byte("Hello long World With Friends"),
	}
	pc, pb := new(uint32), new(uint32)
	for i := 0; i < b.N; i++ {
		*pc, *pb = 0, 0
		hashLittle2(strs[i%3], pc, pb)
	}
}
