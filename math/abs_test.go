package math

import "testing"

func TestAbs(t *testing.T) {
	var num int64
	num = 100
	if Abs(num) != uint64(num) {
		t.Errorf("Err: %s", "正数获取绝对数错误")
	}
	num = -100
	if Abs(num) != uint64(0-num) {
		t.Errorf("Err: %s, Abs(%d) = %d", "负数获取绝对数错误", num, Abs(num))
	}
}

