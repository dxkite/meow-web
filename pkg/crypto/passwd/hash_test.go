package passwd

import "testing"

func TestVerifyPasswordHash(t *testing.T) {
	password := "dxkite"
	hash, err := NewHash(password)
	if err != nil {
		t.Error(err)
	}
	// kbIhOmdPBV2GMX8vrBs95-JZSSExA34nheMXiiSoHE4.cMIoDUBDw9YXN8bjAg5ahg
	t.Log("password", password, hash)
	if rst, err := VerifyHash(password, hash); err != nil {
		t.Error(err)
	} else {
		if rst != true {
			t.Fatalf("verify password error")
		}
	}
}
