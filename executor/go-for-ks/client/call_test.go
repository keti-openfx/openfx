package client

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkCall(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := Call("localhost:50051", []byte("HI, HELLO!"), time.Second)
		fmt.Printf("%d: %s\n", i, res)
	}
}
