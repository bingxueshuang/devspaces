package core

import (
	"bytes"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/cloudflare/bn256"
)

var (
	SizeG1 = 64
	SizeG2 = 129
	SizeGT = 384
	SizeSK = 32
)

var RandomSource = rand.Reader

var (
	ErrCiphertext = errors.New("invalid ciphertext")
	ErrTrapdoor   = errors.New("invalid trapdoor")
	ErrRandom     = errors.New("invalid source of randomness")
)

func SharedKey(pk *PKey, sk *SKey) []byte {
	a := sk.Key
	bP := pk.Key
	abP := new(bn256.G2).ScalarMult(bP, a)
	return abP.Marshal()
}

func KeyGen() (sk *SKey, pk *PKey, err error) {
	a, aP, err := bn256.RandomG2(RandomSource)
	if err != nil {
		err = errors.Join(ErrRandom, err)
		return
	}
	return &SKey{a}, &PKey{aP}, nil
}

func KeyGenServer() (sk *SKey, pk *PKeyServer, err error) {
	b, bQ, err := bn256.RandomGT(RandomSource)
	if err != nil {
		err = errors.Join(ErrRandom, err)
		return
	}
	return &SKey{b}, &PKeyServer{bQ}, nil
}

func PEKS(word []byte, server *PKeyServer, receiver *PKey, sender *SKey) ([]byte, error) {
	ct1, pr, e, err := encryptHelper(word, server.Key, receiver.Key, sender.Key)
	if err != nil {
		return nil, err
	}
	c1 := ct1.Marshal()
	c2 := new(bn256.GT).Add(pr, e).Marshal()
	return bytes.Join([][]byte{c1, c2}, nil), nil
}

func Trapdoor(word []byte, server *PKeyServer, sender *PKey, receiver *SKey) ([]byte, error) {
	ct1, pr, e, err := encryptHelper(word, server.Key, sender.Key, receiver.Key)
	if err != nil {
		return nil, err
	}
	e.Neg(e)
	t1 := ct1.Marshal()
	t2 := new(bn256.GT).Add(pr, e).Marshal()
	return bytes.Join([][]byte{t1, t2}, nil), nil
}

func Test(ciphertext, trapdoor []byte, server *SKey) (ok bool, err error) {
	s1, s2, err := testHelper(ciphertext, trapdoor)
	if err != nil {
		return
	}
	A := s1.ScalarMult(s1, server.Key).Marshal()
	B := s2.Marshal()
	return bytes.Equal(A, B), nil
}

func testHelper(c, t []byte) (s1, s2 *bn256.GT, err error) {
	n := SizeGT
	// prevent index out of bounds
	if len(c) != 2*n {
		err = ErrCiphertext
		return
	}
	if len(t) != 2*n {
		err = ErrTrapdoor
		return
	}
	c1 := new(bn256.GT)
	c2 := new(bn256.GT)
	t1 := new(bn256.GT)
	t2 := new(bn256.GT)
	// only possible error is bad number of bytes
	// which is already handled above
	_, _ = c1.Unmarshal(c[:n])
	_, _ = c2.Unmarshal(c[n:])
	_, _ = t1.Unmarshal(t[:n])
	_, _ = t2.Unmarshal(t[n:])
	s1 = new(bn256.GT).Add(c1, t1)
	s2 = new(bn256.GT).Add(c2, t2)
	return
}

func encryptHelper(w []byte, pubkey *bn256.GT, pk *bn256.G2, sk *big.Int) (ct1, pr, e *bn256.GT, err error) {
	r, err := rand.Int(RandomSource, bn256.Order)
	if err != nil {
		err = errors.Join(ErrRandom, err)
		return
	}
	ct1 = new(bn256.GT).ScalarBaseMult(r)
	h := bn256.HashG1(w, nil)
	k := new(bn256.G2).ScalarMult(pk, sk)
	e = bn256.Pair(h, k)
	pr = new(bn256.GT).ScalarMult(pubkey, r)
	return
}
