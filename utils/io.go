package utils

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func Read(conn net.Conn) error {
	r := bufio.NewReader(conn)
	_, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	return nil
}

func Write(conn net.Conn, msg string) error {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func LogResult(result []float64) {
	strArr := make([]string, len(result))
	for i, v := range result {
		strArr[i] = fmt.Sprintf("%.1f", v) // You can adjust the formatting as needed
	}
	res := strings.Join(strArr, ", ")
	fmt.Printf("[%s]\n", res)
}
