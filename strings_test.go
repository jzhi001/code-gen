package code_gen

import "testing"

func TestCamel2Snake(t *testing.T) {

	s := "CreateAt"

	s = Camel2Snake(s)

	if s != "create_at" {
		t.Fatalf("expected 'create_at', got %s", s)
	}
}
