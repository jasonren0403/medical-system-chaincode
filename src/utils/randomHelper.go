package utils

import (
	"crypto/rand"
	"math/big"
	"time"
)

type RandomGen struct {
	Len  int
	Type string
	Seed int64
}

func NewGenerator(len int, _type string) *RandomGen {
	return &RandomGen{
		Len:  len,
		Type: _type,
		Seed: time.Now().Unix(),
	}
}

func (r *RandomGen) RandIntRange(maxExclusive int64, minInclusive int64) *big.Int {
	if maxExclusive < 0 {
		return big.NewInt(0)
	}
	res, _ := rand.Int(rand.Reader, big.NewInt(maxExclusive-minInclusive))
	return res.Add(res, big.NewInt(minInclusive))

}

func (r *RandomGen) RandPrimeInt() *big.Int {
	if r.Type == "num" {
		res, _ := rand.Prime(rand.Reader, r.Len)
		return res
	}
	return big.NewInt(2)
}

func (r *RandomGen) RandStr() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < r.Len; i++ {
		index := r.RandIntRange(int64(len(bytes)), 0).Int64()
		result = append(result, bytes[index])
	}
	return string(result)
}

func (r *RandomGen) RandByte() []byte {
	b := make([]byte, r.Len)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
