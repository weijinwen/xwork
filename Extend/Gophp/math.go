package Gophp

import (
	"encoding/hex"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// Rand rand()
// Range: [0, 2147483647]
func Rand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	// PHP: getrandmax()
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min
}

// Round round()
func Round(value float64) float64 {
	return math.Floor(value + 0.5)
}

// Max max()
func Max(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}
	max := nums[0]
	for i := 1; i < len(nums); i++ {
		max = math.Max(max, nums[i])
	}
	return max
}

// Min min()
func Min(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}
	min := nums[0]
	for i := 1; i < len(nums); i++ {
		min = math.Min(min, nums[i])
	}
	return min
}

// 16进制转字符转
// Hex2bin hex2bin()
func Hex2bin(data string) string {
	b, err := hex.DecodeString(data)
	if err != nil {
		return ""
	}
	return string(b)
}

// 字符转转16进制
// Bin2hex bin2hex()
func Bin2hex(str string) string {
	return hex.EncodeToString([]byte(str))
}

// IsNan is_nan()
func IsNan(val float64) bool {
	return math.IsNaN(val)
}
