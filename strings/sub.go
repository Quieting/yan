package strings

// SubStr 按字符截取字符串
// str: 操作的字符串
// start: 开始位置,+、- 表示计数方向(非负数数: 从左向右计数; 负整数: 从右向左计数)
// length: 需要的截取的字符长度,+、- 表示截取方向, 没有或者0表示从 start 到结束; 正数: 从左到右截取; 负数: 从右到左截取(字符顺序依旧是原有顺序)
// 1.按字符截取, 不按 byte 截取
// 2.假设 l := len([]rune(str)); |start| >= l: 超出限制返回空字符串; |start| < l && |start| + length[0] >= l: 返回 start 位置向右或向左所有字符
// 3.start 从 1 开始计数, 包含 start 处字符:
//   SubStr("abcdefg", 1, 3) --> "abc"
//   SubStr("abcdefg", -1, 3) --> "efg"
//   SubStr("abcdefg", 3, -2) --> "bc"
//   SubStr("abcdefg", -4, 2) --> "de"
func SubStr(str string, start int, length ...int) string {
	s := []rune(str)
	startAbs := absInt(start)
	if len(s) <= startAbs {
		return ""
	}

	end := 0
	if len(length) == 0 || length[0] == 0 {
		if start > 0 {
			start--
			end = len(s)
		} else {
			end = len(s) + start + 1
			start = 0
		}
		return string(s[start:end])
	}

	// 从左到右计数, 从左到右截取
	if start > 0 && length[0] > 0 {
		if length[0]+start >= len(s) {
			return string(s[start-1:])
		}
		end = start + length[0] - 1
		return string(s[start-1 : end])
	}

	// 从左到右计数, 从右到左截取
	lengAbs := absInt(length[0])
	if start > 0 && length[0] < 0 {
		end = start
		if lengAbs > start {
			return string(s[:end])
		}
		start = start - lengAbs
		return string(s[start:end])
	}

	// 从右到左计数, 从左到右截取
	if start < 0 && length[0] > 0 {
		start = len(s) + start
		if length[0] > startAbs {
			return string(s[startAbs:])
		}

		end = length[0] + start
		return string(s[start:end])
	}

	// 从右到左计数, 从右到左截取
	if start < 0 && length[0] < 0 {
		end = len(s) + start + 1
		if lengAbs > len(s)-startAbs {
			return string(s[:end])
		}
		start = len(s) + start + length[0] + 1
		return string(s[start:end])
	}
	return ""
}

// left 从第 start 个字符处向左截取 num 个字符
func left(str []rune, start, num int) string {
	return string(str[start+num : start])
}

// right 从第 start 个字符处向右截取 num 个字符
func right(str []rune, start, num int) string {
	return string(str[start : start+num])
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0
	}
	return x
}
