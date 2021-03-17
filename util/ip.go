package util

import (
	"strconv"
	"strings"
)

func Ip2Num(ip string) (int, error) {
	values := strings.Split(ip, ".")
	n1, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, err
	}
	n2, err := strconv.Atoi(values[1])
	if err != nil {
		return 0, err
	}
	n3, err := strconv.Atoi(values[2])
	if err != nil {
		return 0, err
	}
	n4, err := strconv.Atoi(values[3])
	if err != nil {
		return 0, err
	}
	return n1*16777216 + n2*65536 + n3*256 + n4, nil
}
