package fastCheck

import (
	"fmt"
	"testing"
)

func TestInSlice(t *testing.T) {
	var (
		a1 = []int{1, 2, 3}
	)

	fmt.Println(InSlice(a1, 0))
}
