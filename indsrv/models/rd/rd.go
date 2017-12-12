package rd

import (
	"errors"
	"math/rand"
	"time"
)

//rd  random
var rd *rand.Rand

func Init() error {
	rd = rand.New(rand.NewSource(time.Now().UnixNano()))
	if rd == nil {
		return errors.New("Init rand failed !")
	}
	return nil
}

func I64() int64 {
	return rd.Int63()
}

//example s = 5 , e = 10 ,生成结果是包含5 和 10的
func I64n(s int64, e int64) int64 {
	e++
	return rd.Int63n(e-s) + s
}

func I32() int32 {
	return rd.Int31()
}

func I32n(s int32, e int32) int32 {
	e++
	return rd.Int31n(e-s) + s
}

func Int() int {
	return rd.Int()
}

func Intn(s int, e int) int {
	e++
	return rd.Intn(e-s) + s
}
