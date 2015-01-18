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

/*void hashlittle2(
  const void *key,       /* the key to hash
  size_t      length,    /* length of the key
  uint32_t   *pc,        /* IN: primary initval, OUT: primary hash
  uint32_t   *pb)        /* IN: secondary initval, OUT: secondary hash
{
  uint32_t a,b,c;                                          /* internal state
  union { const void *ptr; size_t i; } u;     /* needed for Mac Powerbook G4

  /* Set up the internal state
  a = b = c = 0xdeadbeef + ((uint32_t)length) + *pc;
  c += *pb;

  u.ptr = key;
  if (HASH_LITTLE_ENDIAN && ((u.i & 0x3) == 0)) {
    const uint32_t *k = (const uint32_t *)key;         /* read 32-bit chunks
    const uint8_t  *k8;

    /*------ all but last block: aligned reads and affect 32 bits of (a,b,c)
    while (length > 12)
    {
      a += k[0];
      b += k[1];
      c += k[2];
      mix(a,b,c);
      length -= 12;
      k += 3;
    }

    /*----------------------------- handle the last (probably partial) block
    /*
     * "k[2]&0xffffff" actually reads beyond the end of the string, but
     * then masks off the part it's not allowed to read.  Because the
     * string is aligned, the masked-off tail is in the same word as the
     * rest of the string.  Every machine with memory protection I've seen
     * does it on word boundaries, so is OK with this.  But VALGRIND will
     * still catch it and complain.  The masking trick does make the hash
     * noticably faster for short strings (like English words).

#else /* make valgrind happy

    k8 = (const uint8_t *)k;
    switch(length)
    {
    case 12: c+=k[2]; b+=k[1]; a+=k[0]; break;
    case 11: c+=((uint32_t)k8[10])<<16;  /* fall through
    case 10: c+=((uint32_t)k8[9])<<8;    /* fall through
    case 9 : c+=k8[8];                   /* fall through
    case 8 : b+=k[1]; a+=k[0]; break;
    case 7 : b+=((uint32_t)k8[6])<<16;   /* fall through
    case 6 : b+=((uint32_t)k8[5])<<8;    /* fall through
    case 5 : b+=k8[4];                   /* fall through
    case 4 : a+=k[0]; break;
    case 3 : a+=((uint32_t)k8[2])<<16;   /* fall through
    case 2 : a+=((uint32_t)k8[1])<<8;    /* fall through
    case 1 : a+=k8[0]; break;
    case 0 : *pc=c; *pb=b; return;  /* zero length strings require no mixing
    }

#endif /* !valgrind

  } else if (HASH_LITTLE_ENDIAN && ((u.i & 0x1) == 0)) {
    const uint16_t *k = (const uint16_t *)key;         /* read 16-bit chunks
    const uint8_t  *k8;

    /*--------------- all but last block: aligned reads and different mixing
    while (length > 12)
    {
      a += k[0] + (((uint32_t)k[1])<<16);
      b += k[2] + (((uint32_t)k[3])<<16);
      c += k[4] + (((uint32_t)k[5])<<16);
      mix(a,b,c);
      length -= 12;
      k += 6;
    }

    /*----------------------------- handle the last (probably partial) block
    k8 = (const uint8_t *)k;
    switch(length)
    {
    case 12: c+=k[4]+(((uint32_t)k[5])<<16);
             b+=k[2]+(((uint32_t)k[3])<<16);
             a+=k[0]+(((uint32_t)k[1])<<16);
             break;
    case 11: c+=((uint32_t)k8[10])<<16;     /* fall through
    case 10: c+=k[4];
             b+=k[2]+(((uint32_t)k[3])<<16);
             a+=k[0]+(((uint32_t)k[1])<<16);
             break;
    case 9 : c+=k8[8];                      /* fall through
    case 8 : b+=k[2]+(((uint32_t)k[3])<<16);
             a+=k[0]+(((uint32_t)k[1])<<16);
             break;
    case 7 : b+=((uint32_t)k8[6])<<16;      /* fall through
    case 6 : b+=k[2];
             a+=k[0]+(((uint32_t)k[1])<<16);
             break;
    case 5 : b+=k8[4];                      /* fall through
    case 4 : a+=k[0]+(((uint32_t)k[1])<<16);
             break;
    case 3 : a+=((uint32_t)k8[2])<<16;      /* fall through
    case 2 : a+=k[0];
             break;
    case 1 : a+=k8[0];
             break;
    case 0 : *pc=c; *pb=b; return;  /* zero length strings require no mixing
    }

  } else {                        /* need to read the key one byte at a time
    const uint8_t *k = (const uint8_t *)key;

    /*--------------- all but the last block: affect some 32 bits of (a,b,c)
    while (length > 12)
    {
      a += k[0];
      a += ((uint32_t)k[1])<<8;
      a += ((uint32_t)k[2])<<16;
      a += ((uint32_t)k[3])<<24;
      b += k[4];
      b += ((uint32_t)k[5])<<8;
      b += ((uint32_t)k[6])<<16;
      b += ((uint32_t)k[7])<<24;
      c += k[8];
      c += ((uint32_t)k[9])<<8;
      c += ((uint32_t)k[10])<<16;
      c += ((uint32_t)k[11])<<24;
      mix(a,b,c);
      length -= 12;
      k += 12;
    }

    /*-------------------------------- last block: affect all 32 bits of (c)
    switch(length)                   /* all the case statements fall through
    {
    case 12: c+=((uint32_t)k[11])<<24;
    case 11: c+=((uint32_t)k[10])<<16;
    case 10: c+=((uint32_t)k[9])<<8;
    case 9 : c+=k[8];
    case 8 : b+=((uint32_t)k[7])<<24;
    case 7 : b+=((uint32_t)k[6])<<16;
    case 6 : b+=((uint32_t)k[5])<<8;
    case 5 : b+=k[4];
    case 4 : a+=((uint32_t)k[3])<<24;
    case 3 : a+=((uint32_t)k[2])<<16;
    case 2 : a+=((uint32_t)k[1])<<8;
    case 1 : a+=k[0];
             break;
    case 0 : *pc=c; *pb=b; return;  /* zero length strings require no mixing
    }
  }

  final(a,b,c);
  *pc=c; *pb=b;
}*/
