package secret

import (
	"bytes"
	"errors"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type Mapper struct {
	mapper [][]byte
	normal int
}

const (
	LENGTH = 36
	MaxID  = 10000000000000 - 1
	MixID  = 1000000000000
)

var MapperByte = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

// NewMapper ...
func NewMapper() Mapper {
	mapper := make([][]byte, LENGTH)
	rand.Seed(time.Now().UnixNano())
	for k, _ := range mapper {
		mapper[k] = make([]byte, LENGTH)
		copy(mapper[k], MapperByte)
		rand.Shuffle(len(mapper[k]), func(i, j int) {
			mapper[k][i], mapper[k][j] = mapper[k][j], mapper[k][i]
		})
	}
	return Mapper{
		mapper: mapper,
		normal: 0,
	}
}

// Long2String ...
func (m Mapper) Long2String(id int64) (buffer string, err error) {
	if id < 0 {
		err = errors.New("id is negative")
		return
	}
	numSlice := GetRemainderSlice(math.MaxInt64 - id)
	index := numSlice[len(numSlice)-1]
	end := len(numSlice) - 1
	for _, v := range numSlice[:end] {
		buffer += string(m.mapper[index][v])
	}
	buffer += string(m.mapper[m.normal][index])
	return
}

func byteAt(str string, n int) byte {
	if len(str) > n {
		return []byte(str)[n]
	}
	return *new(byte)
}

// String2Long ...
func (m Mapper) String2Long(id string) (int64, error) {
	if len(id) == 0 {
		return -1, errors.New("must has an id")
	}
	index, err := m.getIndexKey(id)
	if err != nil {
		return -1, err
	}
	var num int64
	end := len(id) - 1
	for _, v := range []byte(id)[:end] {
		n := m.getMapperIndex(v, index)
		if n == -1 {
			return -1, errors.New(string(v) + " is not a valid key byte!")
		}
		num = num*int64(LENGTH) + int64(n)
	}
	num = num*int64(LENGTH) + int64(index)
	num = math.MaxInt64 - num
	if num < MixID || num > MaxID {
		return -1, errors.New("id " + id + " is fake")
	}
	return num, nil
}

func (m Mapper) getIndexKey(id string) (int, error) {
	indexByte := byteAt(id, len(id)-1)
	if indexByte == *new(byte) {
		return -1, errors.New("not found")
	}
	index := m.getNormalIndex(indexByte)
	if index == -1 {
		return -1, errors.New(strconv.Itoa(int(indexByte)) + " is not a valid key byte!")
	}
	return index, nil
}

func (m Mapper) getMapperIndex(c byte, index int) int {
	return bytes.IndexByte(m.mapper[index], c)
}

func (m Mapper) getNormalIndex(c byte) int {
	return bytes.IndexByte(m.mapper[m.normal], c)
}

func GetRemainderSlice(rest int64) (stack []int) {
	if 0 == rest {
		stack = append(stack, 0)
	}
	for rest != 0 {
		stack = append(stack, int(rest-(rest/int64(LENGTH))*int64(LENGTH)))
		rest = rest / int64(LENGTH)
	}
	return reverse(stack)
}

func reverse(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// String ...
func (m Mapper) String() string {
	buf := bytes.NewBufferString("{\n")
	for k, _ := range m.mapper {
		buf.WriteString(m.string(k))
	}
	buf.WriteString("}\n")
	return buf.String()
}

func (m Mapper) string(n int) string {
	buf := bytes.NewBufferString("{")
	for _, v := range m.mapper[n] {
		buf.WriteString("'" + string(v) + "',")
	}
	buf.WriteString("},\n")
	return buf.String()
}
