package ReadFile

import (
	"fmt"
	"io/ioutil"
	"os"
)

// 整个文件读取入内存
// 直接将数据直接读取入内存，是效率最高的一种方式，但此种方式，仅适用于小文件，对于大文件，则不适合，因为比较浪费内存。


// 1.1 直接指定文件名读取
// ioutil.ReadFile 就等价于 os.ReadFile，二者是完全一致的

func OSReadFile1(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(content))
}


func IOReadFile(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(content))
}


// 1.2 先创建句柄再读取

func OSFile(filename string) {
	file, err := os.Open(filename)
	// file, err := os.OpenFile("a.txt", os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	fmt.Println(string(content))
}
