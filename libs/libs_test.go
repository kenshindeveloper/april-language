package libs

import "testing"

func TestStack(t *testing.T) {
	s := NewStack()
	s.Push(5)
	s.Push(1)
	s.Push(2)
	s.Push(3)

	if s.Len() != 4 {
		t.Fatalf("Error length must be '4'. got='%d'", s.Len())
	}

	for i := 0; i < 150; i++ {
		s.Pop()
	}

	if s.Len() != 0 {
		t.Fatalf("Error length must be '0'. got='%d'", s.Len())
	}
}
