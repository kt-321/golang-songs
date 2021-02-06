package interfaces

import (
	"golang-songs/model"
	"testing"
)

var user = model.User{Email: "a@test.co.jp", Password: "aaaaaa"}

func BenchmarkCreateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createToken(user)
	}
}
