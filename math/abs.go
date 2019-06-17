package math

// Abs 获取绝对数
func Abs(num int64) uint64 {
	if num > 0 {
		return uint64(num)
	}
	return uint64(^num) + 1
}
