package database

import "testing"

func TestReplace(t *testing.T) {
	got := replace(`"id" int not null`)
	want := "`id` int not null"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
