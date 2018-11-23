package arithmetic

import (
	"fmt"
	"testing"
)

func TestAbandon(t *testing.T) {
	data := []int64{5, 6, 7, 8}
	n := 2
	dst := make([]int64, 0)

	abandon(dst, data, 0, n, 0)
	fmt.Printf("%+v", dst)
}
