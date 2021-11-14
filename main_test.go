package main

import "testing"

func TestServer(t *testing.T) {

	want := "Ok"

	got := test()

	if want != got {
		t.Fatalf("want %s, got %s\n", want, got)
	}
}
