package core

import (
	"bytes"
	"crypto/rand"
	"errors"
	"github.com/cloudflare/bn256"
	"math/big"
	"testing"
	"testing/iotest"
)

var (
	randReader = RandomSource
	errReader  = iotest.ErrReader(ErrRandom)
)

// ## Unit Tests ##

func TestKeyGen(t *testing.T) {
	t.Run("correctness", func(t *testing.T) {
		sk, pk, err := KeyGen()
		handleFatal(err, t)
		got := pk.Key.Marshal()
		want := new(bn256.G2).ScalarBaseMult(sk.Key).Marshal()
		if !bytes.Equal(want, got) {
			t.Logf("expected: %x, got: %x", want, got)
			t.Fatal("secret key is expected to correctly map to the public key")
		}
	})
	// setup erroneous random source
	RandomSource = errReader
	t.Run("error", func(t *testing.T) {
		_, _, err := KeyGen()
		got := err
		want := ErrRandom
		if !errors.Is(got, want) {
			t.Logf("expected: %v, got: %v", want, got)
			t.Fatal("KeyGen is expected to throw randomness error")
		}
	})
	// cleanup erroneous random source
	RandomSource = randReader
}

func TestKeyGenServer(t *testing.T) {
	t.Run("correctness", func(t *testing.T) {
		sk, pk, err := KeyGenServer()
		handleFatal(err, t)
		got := pk.Key.Marshal()
		want := new(bn256.GT).ScalarBaseMult(sk.Key).Marshal()
		if !bytes.Equal(want, got) {
			t.Logf("expected: %x, got: %x", want, got)
			t.Fatal("secret key is expected to correctly map to the public key")
		}
	})
	// setup erroneous random source
	RandomSource = errReader
	t.Run("error", func(t *testing.T) {
		_, _, err := KeyGenServer()
		got := err
		want := ErrRandom
		if !errors.Is(got, want) {
			t.Logf("expected: %v, got: %v", want, got)
			t.Fatal("KeyGenServer is expected to throw randomness error")
		}
	})
	// cleanup erroneous random source
	RandomSource = randReader
}

func TestSharedKey(t *testing.T) {
	// setup two public-secret key pairs
	sk1, pk1, err := KeyGen()
	handleFatal(err, t)
	sk2, pk2, err := KeyGen()
	handleFatal(err, t)

	// run tests
	t.Run("commutative", func(t *testing.T) {
		shared12 := SharedKey(pk1, sk2)
		shared21 := SharedKey(pk2, sk1)
		if !bytes.Equal(shared12, shared21) {
			t.Logf("%x != %x", shared12, shared21)
			t.Fatalf("shared key is expected to be same either way")
		}
	})
	t.Run("length", func(t *testing.T) {
		shared := SharedKey(pk1, sk2)
		got := len(shared)
		want := SizeG2
		if got != want {
			t.Logf("expected: %v, got: %v", want, got)
			t.Fatal("invalid length of shared key")
		}
	})
}

func TestPEKS(t *testing.T) {
	// setup random word
	word, err := getRandomBytes()
	handleFatal(err, t)

	// setup server public key
	_, server, err := KeyGenServer()
	handleFatal(err, t)
	// setup sender secret key
	sender, _, err := KeyGen()
	handleFatal(err, t)
	// setup receiver public key
	_, receiver, err := KeyGen()
	handleFatal(err, t)

	// run tests
	t.Run("correctness", func(t *testing.T) {
		ct, err := PEKS(word, server, receiver, sender)
		handleFatal(err, t)
		got := len(ct)
		want := 2 * SizeGT
		if got != want {
			t.Logf("expected: %v, got: %v", want, got)
			t.Fatal("invalid ciphertext length")
		}
	})

	// setup erroneous random source
	RandomSource = errReader
	t.Run("error", func(t *testing.T) {
		_, err := PEKS(word, server, receiver, sender)
		got := err
		want := ErrRandom
		if !errors.Is(got, want) {
			t.Logf("expected: %v, got %v", want, got)
			t.Fatal("PEKS is expected to throw randomness error")
		}
	})
	RandomSource = randReader
}

func TestTrapdoor(t *testing.T) {
	// setup random word
	word, err := getRandomBytes()
	handleFatal(err, t)

	// setup server public key
	_, server, err := KeyGenServer()
	handleFatal(err, t)
	// setup sender public key
	_, sender, err := KeyGen()
	handleFatal(err, t)
	// setup receiver secret key
	receiver, _, err := KeyGen()
	handleFatal(err, t)

	// run tests
	t.Run("correctness", func(t *testing.T) {
		ct, err := Trapdoor(word, server, sender, receiver)
		handleFatal(err, t)
		got := len(ct)
		want := 2 * SizeGT
		if got != want {
			t.Logf("expected: %v, got: %v", want, got)
			t.Fatal("invalid trapdoor length")
		}
	})

	// setup erroneous random source
	RandomSource = errReader
	t.Run("error", func(t *testing.T) {
		_, err := Trapdoor(word, server, sender, receiver)
		got := err
		want := ErrRandom
		if !errors.Is(got, want) {
			t.Logf("expected: %v, got %v", want, got)
			t.Fatal("Trapdoor is expected to throw randomness error")
		}
	})
	RandomSource = randReader
}

func TestTest(t *testing.T) {
	// setup server key pair
	skServer, pkServer, err := KeyGenServer()
	handleFatal(err, t)
	// setup sender key pair
	skSender, pkSender, err := KeyGen()
	handleFatal(err, t)
	// setup receiver key pair
	skReceiver, pkReceiver, err := KeyGen()
	handleFatal(err, t)

	// setup random word
	word, err := getRandomBytes()
	handleFatal(err, t)
	// setup another random word
	word2, err := getRandomBytes()
	handleFatal(err, t)

	// correct ciphertext
	ciphertext, err := PEKS(word, pkServer, pkReceiver, skSender)
	handleFatal(err, t)
	// invalid ciphertext
	invalidCT, err := getRandomBytes()
	handleFatal(err, t)

	// correct trapdoor
	tdTruthy, err := Trapdoor(word, pkServer, pkSender, skReceiver)
	handleFatal(err, t)
	// incorrect trapdoor
	tdFalsey, err := Trapdoor(word2, pkServer, pkSender, skReceiver)
	handleFatal(err, t)
	// invalid trapdoor
	invalidTD, err := getRandomBytes()
	handleFatal(err, t)

	// run tests
	t.Run("truthy", func(t *testing.T) {
		ok, err := Test(ciphertext, tdTruthy, skServer)
		handleFatal(err, t)
		if !ok {
			t.Logf("expected: %v, got: %v", true, ok)
			t.Fatal("Test is expected to pass")
		}
	})
	t.Run("falsey", func(t *testing.T) {
		ok, err := Test(ciphertext, tdFalsey, skServer)
		handleFatal(err, t)
		if ok {
			t.Logf("expected: %v, got: %v", false, ok)
			t.Fatal("Test is expected to fail")
		}
	})
	t.Run("error", func(t *testing.T) {
		t.Run("ciphertext", func(t *testing.T) {
			_, err := Test(invalidCT, tdTruthy, skServer)
			got := err
			want := ErrCiphertext
			if !errors.Is(got, want) {
				t.Logf("expected: %v, got: %v", want, got)
				t.Fatal("invalid ciphertext error is expected")
			}
		})
		t.Run("trapdoor", func(t *testing.T) {
			_, err := Test(ciphertext, invalidTD, skServer)
			got := err
			want := ErrTrapdoor
			if !errors.Is(got, want) {
				t.Logf("expected: %v, got: %v", want, got)
				t.Fatal("invalid trapdoor error is expected")
			}
		})
	})
}

// ## Benchmarks ##

func BenchmarkKeyGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := KeyGen()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkKeyGenServer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := KeyGenServer()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSharedKey(b *testing.B) {
	// setup two public-secret key pairs
	sk, _, err := KeyGen()
	handleFatal(err, b)
	_, pk, err := KeyGen()
	handleFatal(err, b)

	// run benchmark
	for i := 0; i < b.N; i++ {
		_ = SharedKey(pk, sk)
	}
}

func BenchmarkTrapdoor(b *testing.B) {
	// setup random word
	word, err := getRandomBytes()
	handleFatal(err, b)

	// setup server public key
	_, server, err := KeyGenServer()
	handleFatal(err, b)
	// setup sender public key
	_, sender, err := KeyGen()
	handleFatal(err, b)
	// setup receiver private key
	receiver, _, err := KeyGen()
	handleFatal(err, b)

	// run benchmark
	for i := 0; i < b.N; i++ {
		_, err := Trapdoor(word, server, sender, receiver)
		handleFatal(err, b)
	}
}

func BenchmarkPEKS(b *testing.B) {
	// setup random word
	word, err := getRandomBytes()
	handleFatal(err, b)

	// setup server public key
	_, server, err := KeyGenServer()
	handleFatal(err, b)
	// setup sender private key
	sender, _, err := KeyGen()
	handleFatal(err, b)
	// setup receiver public key
	_, receiver, err := KeyGen()
	handleFatal(err, b)

	// run benchmark
	for i := 0; i < b.N; i++ {
		_, err := PEKS(word, server, receiver, sender)
		handleFatal(err, b)
	}
}

func BenchmarkTest(b *testing.B) {
	// setup server public key
	skServer, pkServer, err := KeyGenServer()
	handleFatal(err, b)
	// setup sender private key
	skSender, pkSender, err := KeyGen()
	handleFatal(err, b)
	// setup receiver public key
	skReceiver, pkReceiver, err := KeyGen()
	handleFatal(err, b)

	// setup random word
	word, err := getRandomBytes()
	handleFatal(err, b)
	// setup another random word
	word2, err := getRandomBytes()
	handleFatal(err, b)

	// ciphertext
	ciphertext, err := PEKS(word, pkServer, pkReceiver, skSender)
	handleFatal(err, b)
	// correct trapdoor
	tdTruthy, err := Trapdoor(word, pkServer, pkSender, skReceiver)
	handleFatal(err, b)
	// incorrect trapdoor
	tdFalsey, err := Trapdoor(word2, pkServer, pkSender, skReceiver)
	handleFatal(err, b)

	// run benchmarks
	b.Run("true", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := Test(ciphertext, tdTruthy, skServer)
			handleFatal(err, b)
		}
	})
	b.Run("false", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := Test(ciphertext, tdFalsey, skServer)
			handleFatal(err, b)
		}
	})
}

// ## Helpers ##

func getRandomBytes() ([]byte, error) {
	// random byte slice whose length is 2^7 to 2^15
	length, err := rand.Int(rand.Reader, big.NewInt(1<<7))
	if err != nil {
		return nil, err
	}
	bytebuffer := make([]byte, length.Int64()+1<<8)
	_, err = rand.Read(bytebuffer)
	if err != nil {
		return nil, err
	}
	return bytebuffer, nil
}

func handleFatal(e error, i interface{ Fatal(args ...any) }) {
	if e != nil {
		i.Fatal(e)
	}
}
