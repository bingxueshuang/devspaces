package core

import (
	"math/big"

	"github.com/cloudflare/bn256"
)

type SKey struct {
	Key *big.Int
}

func (sk *SKey) Decode() []byte {
	return sk.Key.Bytes()
}

func (sk *SKey) Encode(m []byte) {
	sk.Key = new(big.Int).SetBytes(m)
}

type PKey struct {
	Key *bn256.G2
}

func (pk *PKey) Decode() []byte {
	return pk.Key.Marshal()
}

func (pk *PKey) Encode(m []byte) error {
	pk.Key = new(bn256.G2)
	_, err := pk.Key.Unmarshal(m)
	return err
}

type PKeyServer struct {
	Key *bn256.GT
}

func (pk *PKeyServer) Decode() []byte {
	return pk.Key.Marshal()
}

func (pk *PKeyServer) Encode(m []byte) error {
	pk.Key = new(bn256.GT)
	_, err := pk.Key.Unmarshal(m)
	return err
}
