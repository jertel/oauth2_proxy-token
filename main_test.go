package main

import (
	"testing"
)

func TestConvertUser(t *testing.T) {
	actual := convertUser("myuser")
	if "myuser" != actual {
		t.Errorf("Convert user mismatch: %s", actual)
	}
}

func TestConvertUserEmail(t *testing.T) {
	actual := convertUser("myuser@somewhere.invalid")
	if "myuser" != actual {
		t.Errorf("Convert user mismatch: %s", actual)
	}
}

func TestConvertUserMissing(t *testing.T) {
	actual := convertUser("@somewhere.invalid")
	if "" != actual {
		t.Errorf("Convert user mismatch: %s", actual)
	}
}

func TestConvertUserBlank(t *testing.T) {
	actual := convertUser("")
	if "" != actual {
		t.Errorf("Convert user mismatch: %s", actual)
	}
}

func TestConvertUserCorrupt(t *testing.T) {
	actual := convertUser("myuser@somewhere@invalid")
	if "myuser" != actual {
		t.Errorf("Convert user mismatch: %s", actual)
	}
}
