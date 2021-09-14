package secret

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

type Mapper struct {
	mapper    [][]byte
	normalMap []byte
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
	for i := 0; i < LENGTH; i++ {
		tmp := []byte(MapperStr)
		sliceRange(tmp)
		mapper = append(mapper, tmp)
	}
	return Mapper{
		mapper:    mapper,
		normalMap: mapper[0],
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
	end := len(numSlice)-1
	for _, v := range numSlice[:end]{
		buffer += string(m.mapper[index][v])
	}
	buffer += string(m.normalMap[index])
	return
}

func charAt(str string, n int) byte {
	if len(str) > n {
		return []byte(str)[n]
	}
	return *new(byte)
}

// String2Long ...
func (m Mapper) String2Long(appID string) (int64, error) {
	if len(appID) == 0 {
		return -1, errors.New("must has an app id")
	}
	index, err := m.getIndexKey(appID)
	if err != nil {
		return -1, err
	}
	var num int64
	end := len(appID)-1
	for _, v := range []byte(appID)[:end]{
		n := m.getMapperIndex(v, index)
		if n == -1{
			return -1, errors.New(string(v) + " is not a valid key char!")
		}
		num = num*int64(LENGTH) + int64(n)
	}
	num = num*int64(LENGTH) + int64(index)
	num = math.MaxInt64 - num
	if num < MixID || num > MaxID {
		return -1, errors.New("appID " + appID + " is fake")
	}
	return num, nil
}


func (m Mapper) getIndexKey(appID string) (int, error) {
	indexByte := charAt(appID, len(appID)-1)
	if indexByte == *new(byte) {
		return -1, errors.New("未找到")
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
	return bytes.IndexByte(m.normalMap, c)
}

func GetRemainderSlice(rest int64) (stack []int) {
	if 0 == rest {
		stack = append(stack,0)
	}
	for rest != 0 {
		stack = append(stack, int(rest - (rest/int64(LENGTH))*int64(LENGTH)))
		rest = rest / int64(LENGTH)
	}
	return reverse(stack)
}

func reverse(s []int) []int{
	for i, j := 0, len(s)-1; i < j; i,j = i+1, j-1{
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// PrintMapper ...
func (m Mapper) PrintMapper() {

	fmt.Println("{")
	for i := 0; i < LENGTH -1; i++ {
		m.printSlice(i)
	}
	fmt.Println("}")
}

func (m Mapper) printSlice(n int) {
	fmt.Printf("{")
	for i := 0; i < LENGTH; i++ {
		fmt.Printf("'" + string(m.mapper[n][i]) + "'")
		if i != LENGTH-1 {
			fmt.Printf(",")
		}
	}
	fmt.Println("},")
}

func sliceRange(byteSlice []byte) {
	for i := len(byteSlice); i > 1; i-- {
		last := i-1
		idx := rand.Intn(i)
		byteSlice[last], byteSlice[idx] = byteSlice[idx], byteSlice[last]
	}
}