package smartMedicineSystem

import (
	"ccode/src/utils"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestRandom(t *testing.T) {
	numr := utils.NewGenerator(16, "num")
	assert.True(t, numr.RandPrimeInt().IsInt64(), "Returning value is prime int")
	num := numr.RandIntRange(114514, 64)
	assert.GreaterOrEqual(t, num.Int64(), int64(64), "Generated random num is greater than 64")
	assert.Less(t, num.Int64(), int64(114514), "Generated random num is less than 114514")
	strr := utils.NewGenerator(32, "string")
	rstr := strr.RandStr()
	assert.Len(t, rstr, 32, "Generated random string should be", 32)
}

func BenchmarkNumGen(b *testing.B) {
	b.ResetTimer()
	b.N = 50
	numr := utils.NewGenerator(16, "num")
	for i := 0; i < b.N; i++ {
		num := numr.RandIntRange(100, 16)
		log.Println("Generated num is", num.Int64())
	}
}

func BenchmarkStrGen(b *testing.B) {
	b.ResetTimer()
	b.N = 50
	strr := utils.NewGenerator(32, "string")
	for i := 0; i < b.N; i++ {
		str := strr.RandStr()
		log.Println("Generated string is", str)
	}
}
