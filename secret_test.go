package secret

import (
	"math/rand"
	"testing"
)

var mapper Mapper

func init() {
	mapper = NewMapper()
}

func TestLong2String(t *testing.T) {
	mapper.PrintMapper()
	for i := 1; i < 100000; i++ {
		j := rand.Int63n(9000000000000) + 1000000000000
		str, err := mapper.Long2String(j)
		if err != nil {
			t.Fatal(err)
		}
		id, err := mapper.String2Long(str)
		if err != nil {
			t.Fatal(j, err)
		}
		if id != j {
			t.Fatal("inconsistent results;")
		}
	}
}

func TestString2Long(t *testing.T) {
	for i := 1; i < 10000000; i++ {
		j := rand.Int63n(9000000000000) + 1000000000000
		str, err := mapper.Long2String(j)
		if err != nil {
			t.Fatal(err)
		}
		id, err := mapper.String2Long(str)
		if err != nil {
			t.Fatal(j, err)
		}
		if id != j {
			t.Fatal("inconsistent results")
		}
	}

}

func TestPrintMapper(t *testing.T) {
	mapper.PrintMapper()
}


func BenchmarkLong2String(b *testing.B) {
	j := rand.Int63n(9000000000000) + 1000000000000
	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := mapper.Long2String(j)
		if err != nil {
			b.Fatal(err)
		}
	}
}
