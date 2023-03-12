package core

import (
	"math/big"

	"github.com/cloudflare/bn256"
)

type SKey struct {
	Key *big.Int
}

func (sk *SKey) Bytes() []byte {
	return sk.Key.Bytes()
}

func (sk *SKey) FromBytes(m []byte) {
	sk.Key = new(big.Int).SetBytes(m)
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
