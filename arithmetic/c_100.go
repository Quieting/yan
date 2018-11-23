package arithmetic

// 题目：有1、2、3、4个数字，能组成多少个互不相同且无重复数字的三位数？都是多少？
// 分析: 题目满足三个条件：1.需要的是三位数 2.三位数三个数字互补重复 3.可使用数字1、2、3、4

// CDiffNumber 根据所给数字集合 data，组合成互不相同的 n 位数切片 res 返回
// len(data) 有效范围 1-10, n 有效范围 1-10，超出范围会 panic 异常
func CDiffNumber(data []int64, n int) (res []int64) {
	if n == 0 || n > 10 {
		panic("无效的多位数选择")
	}
	length := len(data)
	if length == 0 || length > 10 {
		panic("无效的数据集合")
	}
	if n > length {
		panic("无效的参数")
	}

	return nil
}

// combination 组合  数学概念组合实现并返回不同的组合结果
// C(m, n) = m*m-1*...(m-n+1)=m!/n!*!(m-n)
func combination(data []int64, n int) (res []int64) {
	length := len(data)
	if n == length {
		return data
	}
	
	flag := make([]int64, 0, len(data))
	for key, _ := range data {
		src[key] 1
	}

	return
}
