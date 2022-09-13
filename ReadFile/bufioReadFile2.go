package ReadFile

import (
	"bufio"
	"io"
	"os"
	"fmt"
)
// https://mp.weixin.qq.com/s/ww27OPuD_Pse_KDNQWyjzA

// 3. 每次只读取固定字节数

func bufioTry() {
	// 先创建一个文件句柄，可以使用 os.Open 或者  os.OpenFile
	fi, err := os.Open("a.txt")
	if err != nil {
		panic(err)
	}

	// 创建 Reader
	r := bufio.NewReader(fi)

	// 每次读取 1024 个字节
	buf := make([]byte, 1024)
	// 然后在 for 循环里调用  Reader 的 Read 函数，每次仅读取固定字节数量的数据。
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}
		fmt.Println(string(buf[:n]))
	}
}
