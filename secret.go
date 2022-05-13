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
	MapperStr = "0123456789abcdefghijklmnopqrstuvwxyz"
	LENGTH    = 36
	MaxID     = 10000000000000 - 1
	MixID     = 1000000000000
)

// NewMapper ...
func NewMapper() Mapper {
	var mapper [][]byte
	rand.Seed(time.Now().UnixNano())
	tmp := []byte(MapperStr)
	for i := 0; i < LENGTH; i++ {
		tmp = []byte(MapperStr)
		rand.Shuffle(len(tmp), func(i, j int) {
			tmp[i], tmp[j] = tmp[j], tmp[i]
		})
		mapper = append(mapper, tmp)
	}
	return Mapper{
		mapper: mapper,
		normal: 0,
	}
}

// Long2String ...
func (m Mapper) Long2String(appID int64) (buffer string, err error) {
	if appID < 0 {
		err = errors.New("appID is negative")
		return
	}
	numSlice := GetRemainderSlice(math.MaxInt64 - appID)
	index := numSlice[len(numSlice)-1]
	end := len(numSlice) - 1
	for _, v := range numSlice[:end] {
		buffer += string(m.mapper[index][v])
	}
	buffer += string(m.mapper[m.normal][index])
	return
}

func charAt(str string, n int) byte {
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
			return -1, errors.New(string(v) + " is not a valid key char!")
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

func (m Mapper) getIndexKey(appID string) (int, error) {
	indexByte := charAt(appID, len(appID)-1)
	if indexByte == *new(byte) {
		return -1, errors.New("not found")
	}
	index := m.getNormalIndex(indexByte)
	if index == -1 {
		return -1, errors.New(strconv.Itoa(int(indexByte)) + " is not a valid key char!")
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
