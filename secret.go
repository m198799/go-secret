package secret

import (
	"bytes"
	"errors"
	"math"
	"math/rand"
	"time"
)

type byteSlice []byte

// Mapper 二维[]byte 保存加密的密钥 第一行用于索引 normal 为 0
type Mapper struct {
	mapper []byteSlice
	normal int
}

const (
	LENGTH = 36
	MaxID  = 10000000000000 - 1
	MixID  = 1000000000000
)

var mapperByte = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

// NewMapper ...
func NewMapper() Mapper {
	mapper := make([]byteSlice, LENGTH)
	rand.Seed(time.Now().UnixNano())
	for k, _ := range mapper {
		mapper[k] = make([]byte, LENGTH)
		copy(mapper[k], mapperByte)
		rand.Shuffle(len(mapper[k]), func(i, j int) {
			mapper[k][i], mapper[k][j] = mapper[k][j], mapper[k][i]
		})
	}
	return Mapper{
		mapper: mapper,
		normal: 0,
	}
}

func checkID(id int64) bool {
	if id < MixID || id > MaxID {
		return false
	}
	return true
}

// EncodeID 将 id int64 转换成一个 string
func (m Mapper) EncodeID(id int64) (str string, err error) {
	if !checkID(id) {
		err = errors.New("id is negative")
		return
	}
	numSlice := remainderSlice(math.MaxInt64 - id)
	end := len(numSlice) - 1
	index := numSlice[end]
	for _, v := range numSlice[:end] {
		str += string(m.mapper[index][v])
	}
	str += string(m.mapper[m.normal][index])
	return
}

// DecodeID 将 加密的 id string 转换成一个 id int64
func (m Mapper) DecodeID(str string) (int64, error) {
	strLen := len(str)
	if strLen == 0 {
		return -1, errors.New("the string cannot be empty")
	}
	indexByte := str[len(str)-1]
	index, err := m.getIndexPosition(indexByte)
	if err != nil {
		return -1, err
	}
	var num int64
	for _, v := range []byte(str)[:strLen-1] {
		n := bytes.IndexByte(m.mapper[index], v)
		if n == -1 {
			return -1, errors.New(string(v) + " is not in mapper!")
		}
		num = num*int64(LENGTH) + int64(n)
	}
	num = num*int64(LENGTH) + int64(index)
	num = math.MaxInt64 - num
	if !checkID(num) {
		return -1, errors.New("string: " + str + " is fake")
	}
	return num, nil
}

func (m Mapper) getIndexPosition(b byte) (index int, err error) {
	index = bytes.IndexByte(m.mapper[m.normal], b)
	if index == -1 {
		return index, errors.New(string(b) + " is not in mapper!")
	}
	return index, nil
}

// 将一个数字 除 LENGTH 取余数, 将余数放入 []int，将得到商 继续除 LENGTH 取余数, 将余数放入 []int，直到商为 0
func remainderSlice(rest int64) (stack []int) {
	if 0 == rest {
		stack = append(stack, 0)
	}
	for rest != 0 {
		stack = append(stack, int(rest-(rest/int64(LENGTH))*int64(LENGTH)))
		rest = rest / int64(LENGTH)
	}
	return reverse(stack)
}

// 反转 []int
func reverse(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// String Mapper 实现 Stringer 接口
func (m Mapper) String() string {
	buf := bytes.NewBufferString("{\n")
	for _, v := range m.mapper {
		buf.WriteString(v.String())
	}
	buf.WriteString("}\n")
	return buf.String()
}

// String byteSlice 实现 Stringer 接口
func (m byteSlice) String() string {
	buf := bytes.NewBufferString("{")
	for _, v := range m {
		buf.WriteString("'" + string(v) + "',")
	}
	buf.WriteString("},\n")
	return buf.String()
}
