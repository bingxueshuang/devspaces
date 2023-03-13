package core

import (
	"bytes"
	"crypto/rand"
	"github.com/cloudflare/bn256"
	"testing"
)

func TestSKey(t *testing.T) {
	// any curve is fine, since all secret keys have same size range
	validInt, _, err := bn256.RandomG1(rand.Reader)
	handleFatal(err, t)
	validBytes := validInt.Bytes()
	validSKey := &SKey{
		Key: validInt,
	}
	t.Run("ToBytes", func(t *testing.T) {
		got := validSKey.Bytes()
		want := validBytes
		if !bytes.Equal(got, want) {
			t.Logf("expected: %v, got: %v", want, got)
			t.Fatal("incorrect secret key to bytes")
		}
	})
	t.Run("FromBytes", func(t *testing.T) {
		sk := new(SKey)
		err := sk.FromBytes(validBytes)
		handleFatal(err, t)
		got := sk.Key
		want := validInt
		if want.Cmp(got) != 0 {
			t.Logf("expected: %x, got: %x", want, got)
			t.Fatal("incorrect secret key from bytes")
		}
	})
	t.Run("FromSKey", func(t *testing.T) {
		sk := new(SKey)
		err := sk.FromSKey(validSKey)
		handleFatal(err, t)
		got := sk.Key
		want := validInt
		if want.Cmp(got) != 0 {
			t.Logf("expected: %x, got: %x", want, got)
			t.Fatal("incorrect secret key from SKey")
		}
	})
}

func TestPkey(t *testing.T) {
	zero := []byte{0x0}
	sk, validKey, err := bn256.RandomG2(rand.Reader)
	validSKey := &SKey{Key: sk}
	handleFatal(err, t)
	validBytes := validKey.Marshal()
	t.Run("ToBytes", func(t *testing.T) {
		pk := &PKey{
			Key: validKey,
		}
		// check for panics
		got := len(pk.Bytes())
		want := SizeG2
		if got != want {
			t.Logf("expected: %v, got: %v", want, got)
			t.Fatal("incorrect public key to bytes")
		}
	})
	t.Run("FromBytes", func(t *testing.T) {
		pk := new(PKey)
		err := pk.FromBytes(validBytes)
		handleFatal(err, t)
		tmp := new(bn256.G2).Neg(pk.Key)
		tmp.Add(tmp, validKey)
		if !bytes.Equal(zero, tmp.Marshal()) {
			t.Fatal("incorrect public key from bytes")
		}
	})
	t.Run("FromSKey", func(t *testing.T) {
		pk := new(PKey)
		err := pk.FromSKey(validSKey)
		handleFatal(err, t)
		tmp := new(bn256.G2).Neg(pk.Key)
		tmp.Add(tmp, validKey)
		if !bytes.Equal(zero, tmp.Marshal()) {
			t.Fatal("incorrect public key from SKey")
		}
	})
}

func TestPkeyServer(t *testing.T) {
	zero := make([]byte, SizeGT)
	zero[len(zero)-1] = 0x1
	sk, validKey, err := bn256.RandomGT(rand.Reader)
	validSKey := &SKey{Key: sk}
	handleFatal(err, t)
	validBytes := validKey.Marshal()
	t.Run("ToBytes", func(t *testing.T) {
		pk := &PKeyServer{
			Key: validKey,
		}
		// check for panics
		got := len(pk.Bytes())
		want := SizeGT
		if got != want {
			t.Logf("expected: %v, got: %v", want, got)
			t.Fatal("incorrect public key to bytes")
		}
	})
	t.Run("FromBytes", func(t *testing.T) {
		pk := new(PKeyServer)
		err := pk.FromBytes(validBytes)
		handleFatal(err, t)
		tmp := new(bn256.GT).Neg(pk.Key)
		tmp.Add(tmp, validKey)
		if !bytes.Equal(zero, tmp.Marshal()) {
			t.Logf("expected: %v, got: %d", zero, len(tmp.Marshal()))
			t.Fatal("incorrect public key from bytes")
		}
	})
	t.Run("FromSKey", func(t *testing.T) {
		pk := new(PKeyServer)
		err := pk.FromSKey(validSKey)
		handleFatal(err, t)
		tmp := new(bn256.GT).Neg(pk.Key)
		tmp.Add(tmp, validKey)
		if !bytes.Equal(zero, tmp.Marshal()) {
			t.Fatal("incorrect public key from SKey")
		}
	})
}
