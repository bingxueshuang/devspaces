package core

import (
	"bytes"
	"crypto/rand"
	"math/big"

	"github.com/cloudflare/bn256"
)

var (
	SizeG1 = 64
	SizeG2 = 129
	SizeGT = 384
)

func KeyGen() (sk *big.Int, pk *PKey, err error) {
	a, aP, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		return
	}
	return a, &PKey{aP}, nil
}

func KeyGenServer() (sk *big.Int, pk *PKeyServer, err error) {
	b, bQ, err := bn256.RandomGT(rand.Reader)
	if err != nil {
		return
	}
	return b, &PKeyServer{bQ}, nil
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
	c1 := new(bn256.GT)
	c2 := new(bn256.GT)
	t1 := new(bn256.GT)
	t2 := new(bn256.GT)
	_, err = c1.Unmarshal(c[:n])
	if err != nil {
		return
	}
	_, err = c2.Unmarshal(c[n:])
	if err != nil {
		return
	}
	_, err = t1.Unmarshal(t[:n])
	if err != nil {
		return
	}
	_, err = t2.Unmarshal(t[n:])
	if err != nil {
		return
	}
	s1 = new(bn256.GT).Add(c1, t1)
	s2 = new(bn256.GT).Add(c2, t2)
	return
}

func encryptHelper(w []byte, pubkey *bn256.GT, pk *bn256.G2, sk *big.Int) (ct1, pr, e *bn256.GT, err error) {
	r, err := rand.Int(rand.Reader, bn256.Order)
	if err != nil {
		return
	}
	ct1 = new(bn256.GT).ScalarBaseMult(r)
	h := bn256.HashG1(w, nil)
	k := new(bn256.G2).ScalarMult(pk, sk)
	e = bn256.Pair(h, k)
	pr = new(bn256.GT).ScalarMult(pubkey, r)
	return
}
