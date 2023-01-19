package tools

import (
	"bufio"
	"os"
)

func ReadFromConfig(filename string) (list []string) {
	r, _ := os.Open(filename)
	defer r.Close()
	s := bufio.NewScanner(r)
	for s.Scan() { // 循环直到文件结束
		line := s.Text() // 这个 line 就是每一行的文本了，string 类型
		list = append(list, line)
	}
	return list
}
