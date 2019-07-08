package gtoken_test

import (
	"gtoken/gtoken"
	"testing"
)

func TestEncryptDecryptToken(t *testing.T) {
	t.Log("encrypt and decrypt token test ")
	gtoken := gtoken.GfToken{}
	gtoken.Init()

	userKey := "123123"
	token := gtoken.EncryptToken(userKey)
	if !token.Success() {
		t.Error(token.Json())
	}
	t.Log(token.DataString())

	token2 := gtoken.DecryptToken(token.GetString("token"))
	if !token2.Success() {
		t.Error(token2.Json())
	}
	t.Log(token2.DataString())
	if userKey != token2.GetString("userKey") {
		t.Error("token decrypt userKey error")
	}
	if token.GetString("uuid") != token2.GetString("uuid") {
		t.Error("token decrypt uuid error")
	}

}

func BenchmarkEncryptDecryptToken(b *testing.B) {
	b.Log("encrypt and decrypt token test ")
	gtoken := gtoken.GfToken{}
	gtoken.Init()

	userKey := "123123"
	token := gtoken.EncryptToken(userKey)
	if !token.Success() {
		b.Error(token.Json())
	}
	b.Log(token.DataString())

	for i := 0; i < b.N; i++ {
		token2 := gtoken.DecryptToken(token.GetString("token"))
		if !token2.Success() {
			b.Error(token2.Json())
		}
		b.Log(token2.DataString())
		if userKey != token2.GetString("userKey") {
			b.Error("token decrypt userKey error")
		}
		if token.GetString("uuid") != token2.GetString("uuid") {
			b.Error("token decrypt uuid error")
		}
	}

}
