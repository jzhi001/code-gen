package idl_conv

import "testing"

func TestCharIterator_Next(t *testing.T) {

	s := "abcde"

	cit := NewCharIterator(s)

	for i := 0; i < len(s); i++ {

		expected := s[i : i+1]
		c, err := cit.Next()
		if err != nil {
			t.Fatal("early end " + err.Error())
		}

		if c != expected {
			t.Fatalf("expect %s but got %s", expected, c)
		}
	}
}

func TestCharIterator_JumpTo(t *testing.T) {

	s := "abcde"

	cit := NewCharIterator(s)

	err := cit.JumpTo([]rune(s)[2])

	if err != nil {
		t.Fatalf("jump should success, but failed: %s", err.Error())
	}

	if cit.i != 2 {
		t.Fatalf("i should equals 2, but got %d", cit.i)
	}
}

func TestCharIterator_SkipUntil(t *testing.T) {

	s := "abcde"

	cit := NewCharIterator(s)

	err := cit.SkipUntil([]rune(s)[2])

	if err != nil {
		t.Fatalf("sip should success, but failed: %s", err.Error())
	}

	if cit.i != 1 {
		t.Fatalf("i should equals 1, but got %d", cit.i)
	}
}
