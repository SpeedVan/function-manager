package controller

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(url2funcConfig("http://gitlab.puhuitech.cn/shizhiyin/temporary/raw/bf1c72c26a9e80a9ac41d6860df482cd52dfcefd/indicator/python/ZHJ.py"))

	fmt.Println(url2funcConfig("http://gitlab.puhuitech.cn/lidongchen/FF_PB_workflow/raw/master/"))
	fmt.Println(url2funcConfig("http://gitlab.puhuitech.cn/lidongchen/FF_PB_workflow/raw/master"))

	fmt.Println(url2funcConfig("http://gitlab.puhuitech.cn/lidongchen/FF_PB_workflow/"))
	fmt.Println(url2funcConfig("http://gitlab.puhuitech.cn/lidongchen/FF_PB_workflow"))
}
