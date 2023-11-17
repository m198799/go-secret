package secret

import (
	"fmt"
	"math/rand"
	"testing"
)

var mapper Mapper

func init() {
	mapper = NewMapper()
}

func TestMapper_String(t *testing.T) {
	fmt.Println(mapper)
}

func TestMapper_DecodeID(t *testing.T) {
	for i := 1; i < 10000000; i++ {
		j := rand.Int63n(MaxID-MixID) + MixID
		str, err := mapper.EncodeID(j)
		if err != nil {
			t.Fatal(err)
		}
		id, err := mapper.DecodeID(str)
		if err != nil {
			t.Fatal(j, err)
		}
		if id != j {
			t.Fatal("inconsistent results;")
		}
	}
}

func BenchmarkMapper_EncodeID(b *testing.B) {
	j := rand.Int63n(9000000000000) + 1000000000000
	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		str, err := mapper.EncodeID(j)
		if err != nil {
			b.Fatal(err)
		}
		_, err = mapper.DecodeID(str)
		if err != nil {
			b.Fatal(j, err)
		}
	}
}
