package gtoken

import (
	"testing"
)

func TestMsg(t *testing.T) {
	if msgLog("123") != "[GToken]123" {
		t.Error("msg err")
	}
	if msgLog("123-%s-%d", "123", 44) != "[GToken]123-123-44" {
		t.Error("msg sprintf err")
	}
}
