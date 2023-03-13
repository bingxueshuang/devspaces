package core

import (
	"math/big"

	"github.com/cloudflare/bn256"
)

type EllipticKey interface {
	Bytes() []byte
	FromBytes(m []byte) error
	FromSKey(sk *SKey) error
}

type SKey struct {
	Key *big.Int
}

func (sk *SKey) Bytes() []byte {
	return sk.Key.Bytes()
}

func (sk *SKey) FromBytes(m []byte) error {
	sk.Key = new(big.Int).SetBytes(m)
	return nil
}

func (sk *SKey) FromSKey(skey *SKey) error {
	sk.Key = new(big.Int).Set(skey.Key)
	return nil
}

type PKey struct {
	Key *bn256.G2
}

func (pk *PKey) Bytes() []byte {
	return pk.Key.Marshal()
}

func (pk *PKey) FromBytes(m []byte) error {
	pk.Key = new(bn256.G2)
	_, err := pk.Key.Unmarshal(m)
	return err
}

func (pk *PKey) FromSKey(sk *SKey) error {
	pk.Key = new(bn256.G2).ScalarBaseMult(sk.Key)
	return nil
}

type PKeyServer struct {
	Key *bn256.GT
}

func (pk *PKeyServer) Bytes() []byte {
	return pk.Key.Marshal()
}

func (pk *PKeyServer) FromBytes(m []byte) error {
	pk.Key = new(bn256.GT)
	_, err := pk.Key.Unmarshal(m)
	return err
}

func (pk *PKeyServer) FromSKey(sk *SKey) error {
	pk.Key = new(bn256.GT).ScalarBaseMult(sk.Key)
	return nil
}

var _ EllipticKey = new(SKey)
var _ EllipticKey = new(PKey)
var _ EllipticKey = new(PKeyServer)
