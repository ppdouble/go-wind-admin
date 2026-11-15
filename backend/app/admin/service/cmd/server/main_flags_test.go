package main

import (
	"reflect"
	"testing"
)

func TestNormalizeLegacyFlags(t *testing.T) {
	input := []string{"-conf", "./configs", "-env=dev", "-chost", "127.0.0.1:8500", "-ctype=consul", "-daemon", "-c", "./other", "--", "-conf"}
	want := []string{"--conf", "./configs", "--env=dev", "--chost", "127.0.0.1:8500", "--ctype=consul", "--daemon", "-c", "./other", "--", "-conf"}
	got := normalizeLegacyFlags(input)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeLegacyFlags() = %q, want %q", got, want)
	}
	if !reflect.DeepEqual(input, []string{"-conf", "./configs", "-env=dev", "-chost", "127.0.0.1:8500", "-ctype=consul", "-daemon", "-c", "./other", "--", "-conf"}) {
		t.Fatal("normalizeLegacyFlags modified the caller's arguments")
	}
}
