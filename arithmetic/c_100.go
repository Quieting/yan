package arithmetic

import (
	"fmt"
)

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

// combination 组合，数学概念组合实现并返回不同的组合结果
// C(m, n) = m*m-1*...(m-n+1)=m!/n!*!(m-n)
// 算法:
// 1.新建长度 len(data) 的数值切片 falg, 前 n(n < len(data))切片值为 1, 其余为 0
// 2.遍历新切片 flag, 遇见第一个连续两个切片的值为1、0时, 交换其值
// 3.将 2步骤中找到的1、0组合前的1全部移动到切片最前面
// 4.重复 2、3 步骤, 直至没有1、0组合
// 5.每次 2、3 步骤得到完成时, 取出当时 flag[i] = 1 对应的 data[i] 的值即为一个组合结果
func combination(data []int64, n int) (res []int64) {

	return
}

func combination2(data []int64, n int) (res [][]int64) {
	if n == len(data) {
		return [][]int64{data}
	}
	res = make([][]int64, 0, 10)
	for i := 0; i < len(data); i++ {
		c := []int64{data[i]}
		for j := i + 1; j < len(data); j++ {
			addone(res, append(c, data[j]))
		}
	}
	fmt.Printf("第一次计算结果:%#v\n", res)

	return
}

func com(data, curr []int64, resu [][]int64, index, dept int) {
	if dept == 1 {
		for j := index + 1; j < len(data); j++ {
			resu = append(resu, append(curr, data[j]))
		}
	} else {
		for i := index + 1; i < len(data); i++ {
			curr = append(curr, data[index])
			com(data, curr, resu, i, dept-1)
		}
	}
}

func addone(res [][]int64, d []int64) {
	res = append(res, d)
}

// 请教一个问题，在 golang 中 []int64 类型作为函数参数传入时是传递的地址值，而 [][]int64 却不是什么
//
func test() {
	sliceOne := make([]int64, 0, 10)
	sliceOne = append(sliceOne, 2)
	fmt.Printf("sliceOne: %#v, %p\n", sliceOne, sliceOne)
	addOne(sliceOne, 1)
	fmt.Printf("sliceOne: %#v, %p\n", sliceOne, sliceOne)
	fmt.Printf("len(sliceOne):%d\n", len(sliceOne))
}
func addOne(res []int64, d int64) {
	fmt.Printf("进入前 res: %#v, %p, cap(res):%d\n", res, res, cap(res))
	res = append(res, d)
	fmt.Printf("追加后: res: %#v, %p, cap(res):%d\n", res, res, cap(res))
}
