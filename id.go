package main

import (
	"github.com/sony/sonyflake"
	"strconv"
)

var flake = sonyflake.NewSonyflake(sonyflake.Settings{})

func NextID() string {
	id, _ := flake.NextID()
	return strconv.FormatUint(id, 16)
}
