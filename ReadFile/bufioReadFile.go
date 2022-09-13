package ReadFile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// 一次性读取所有的数据，太耗费内存，因此可以指定每次只读取一行数据。
// 1. bufio.ReadBytes('\n')
// 2. bufio.ReadString('\n')

// 每次只读取一行

// 使用bufio.ReadBytes('\n')
func bufio1(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(file)

	for {
		lineByte, err := r.ReadBytes("\n")
		line := strings.TrimSpace(string(lineByte))
		if err != nil && err != io.EOF {
			panic(err)
		}

		if err == io.EOF {
			break
		}

		fmt.Println(line)
	}
}


// 使用bufio.ReadString('\n')
func bufio2(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(file)

	for {
		lineByte, err := r.ReadString("\n")
		line := strings.TrimSpace(lineByte)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if err == io.EOF {
			break
		}

		fmt.Println(line)
	}
}
