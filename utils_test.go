package main

import (
	"fmt"
	"testing"
	"time"
)

func assertTrue(t *testing.T, a bool, message string) {
	if a {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != true", a)
	}
	t.Fatal(message)
}

func TestIsHoursAgo(t *testing.T) {
	assertTrue(t, IsDurationAgo("2018-07-02 06:00:34.573", 1*time.Hour), "1小时之前")
}
