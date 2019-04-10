// Copyright 2015 Sermo Digital, LLC. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package helpers

import (
	"errors"
	"net"
	"strconv"
)

const (
	digits   = "0123456789abcdefghijklmnopqrstuvwxyz"
	digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
	digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"
)

// Length finds the number of digits in a uint64. For example, 12 returns 2,
// 100 returns 3, and 1776 returns 4. The minimum width is 1.
func Length(x uint64) int {
	// TODO: use math/bits when it's merged in 1.9

	// Loop: for loop
	// Log10: math.Log10
	// Asm: https://graphics.stanford.edu/~seander/bithacks.html#IntegerLog10
	// 	with "IntegerLogBase2(v)" in assembly implemented similar to GCC's
	// 	__builtin_clzll.
	// Cond: the switch below
	//
	// 	Seq: sequential numbers
	// 	Rand: random numbers
	//
	// BenchmarkIntLenLoopSeq-4     	200000000         7.87 ns/op
	// BenchmarkIntLenLog10Seq-4    	30000000        43.8 ns/op
	// BenchmarkIntLenAsmSeq-4      	200000000         7.38 ns/op
	// BenchmarkIntLenCondSeq-4     	500000000         3.67 ns/op
	//
	// BenchmarkIntLenLoopRand-4    	10000000000000000000        21.8 ns/op
	// BenchmarkIntLenLog10Rand-4   	30000000        44.5 ns/op
	// BenchmarkIntLenAsmRand-4     	200000000         8.62 ns/op
	// BenchmarkIntLenCondRand-4    	200000000         8.09 ns/op

	switch {
	case x < 10:
		return 1
	case x < 100:
		return 2
	case x < 1000:
		return 3
	case x < 10000:
		return 4
	case x < 100000:
		return 5
	case x < 1000000:
		return 6
	case x < 10000000:
		return 7
	case x < 100000000:
		return 8
	case x < 1000000000:
		return 9
	case x < 10000000000:
		return 10
	case x < 100000000000:
		return 11
	case x < 1000000000000:
		return 12
	case x < 10000000000000:
		return 13
	case x < 100000000000000:
		return 14
	case x < 1000000000000000:
		return 15
	case x < 10000000000000000:
		return 16
	case x < 100000000000000000:
		return 17
	case x < 1000000000000000000:
		return 18
	case x < 10000000000000000000:
		return 19
	default:
		panic("unreachable")
	}
}

// ParseIP returns a valid IP address for the given *http.Request.RemoteAddr
// if available.
func ParseIP(remoteaddr string) (string, error) {
	host, _, err := net.SplitHostPort(remoteaddr)
	if err != nil {
		return "", err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return "", errors.New("Invalid IP Address")
	}
	return ip.String(), nil
}

// FormatUint serializes a uint64. It's borrowed from the standard library's
// strconv package, but with the signed cases removed and only formats in base
// 10. If the return value is being converted to a string using strconv
// directly might be faster.
func FormatUint(u uint64) []byte {
	// Special case.
	if u <= 9 {
		return []byte{byte(u) + '0'}
	}

	var a [64]byte
	i := 64

	if ^uintptr(0)>>32 == 0 {
		for u > uint64(^uintptr(0)) {
			q := u / 1e9
			us := uintptr(u - q*1e9) // us % 1e9 fits into a uintptr
			for j := 9; j > 0; j-- {
				i--
				qs := us / 10
				a[i] = byte(us - qs*10 + '0')
				us = qs
			}
			u = q
		}
	}

	// u guaranteed to fit into a uintptr
	us := uintptr(u)
	for us >= 10 {
		i--
		q := us / 10
		a[i] = byte(us - q*10 + '0')
		us = q
	}
	// u < 10
	i--
	a[i] = byte(us + '0')

	return a[i:]
}

// ParseUint is like ParseInt but for unsigned numbers. It's stolen from the
// strconv package and streamlined for uint64s.
func ParseUint(s []byte) (n uint64, err error) {
	const maxUint64 = (1<<64 - 1)
	var cutoff, maxVal uint64

	cutoff = maxUint64/10 + 1
	maxVal = 1<<uint(64) - 1

	for i := 0; i < len(s); i++ {
		var v byte
		d := s[i]
		switch {
		case '0' <= d && d <= '9':
			v = d - '0'
		case 'a' <= d && d <= 'z':
			v = d - 'a' + 10
		case 'A' <= d && d <= 'Z':
			v = d - 'A' + 10
		default:
			n = 0
			err = strconv.ErrSyntax
			goto Error
		}
		if v >= 10 {
			n = 0
			err = strconv.ErrSyntax
			goto Error
		}

		if n >= cutoff {
			// n*base overflows
			n = maxUint64
			err = strconv.ErrRange
			goto Error
		}
		n *= 10

		n1 := n + uint64(v)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			n = maxUint64
			err = strconv.ErrRange
			goto Error
		}
		n = n1
	}

	return n, nil

Error:
	return n, &strconv.NumError{Func: "ParseUint", Num: string(s), Err: err}
}
