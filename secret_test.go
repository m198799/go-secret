/*
 * @Author: panxu
 * @Date: 2022-05-13 10:15:56
 * @LastEditors: panxu
 * @LastEditTime: 2022-05-13 15:03:09
 * @FilePath: /go-secret/secret_test.go
 */
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

func TestLong2String(t *testing.T) {
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

func BenchmarkLong2String(b *testing.B) {
	j := rand.Int63n(9000000000000) + 1000000000000
	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		str, err := mapper.Long2String(j)
		if err != nil {
			b.Fatal(err)
		}
		_, err = mapper.String2Long(str)
		if err != nil {
			b.Fatal(j, err)
		}
	}
}
