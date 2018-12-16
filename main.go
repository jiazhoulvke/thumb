package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jiazhoulvke/goutil"
)

//Size 尺寸
type Size struct {
	Width  int
	Height int
	Suffix string
}

var (
	sourcePath string
	// sizes      []Size
	sizestr string
	//AllowExts 允许的后缀
	allowExtStr string
	override    bool
)

func init() {
	flag.StringVar(&sourcePath, "path", "", "要生成缩略图的路径")
	flag.StringVar(&sizestr, "sizes", "100x100:_t", "缩略图的尺寸,如100x100:_t,200x200_s,前面的100x100是图片的长宽，后面的_t表示缩略图的文件名后缀")
	flag.StringVar(&allowExtStr, "allow_ext", "jpg,jpeg,png,gif", "允许的图片后缀")
	flag.BoolVar(&override, "override", false, "当缩略图存在时是否覆盖")
}

func main() {
	flag.Parse()

	var err error
	_, err = os.Stat(sourcePath)
	if os.IsNotExist(err) {
		fmt.Printf("目录 %s 不存在\n", sourcePath)
		os.Exit(1)
	}
	allowExts := make([]string, 0)
	if allowExtStr == "" {
		fmt.Println("allow_ext不能为空")
		os.Exit(1)
	}
	for _, s := range strings.Split(allowExtStr, ",") {
		if s == "" {
			continue
		}
		allowExts = append(allowExts, "."+strings.ToLower(s))
	}

	sizes := make([]Size, 0)
	for _, s := range strings.Split(sizestr, ",") {
		if s == "" {
			continue
		}
		a := strings.Split(s, ":")
		if len(a) != 2 {
			fmt.Printf("缩略图尺寸 %s 有误\n", s)
			os.Exit(1)
		}
		suffix := a[1]
		l := strings.Split(a[0], "x")
		if len(l) != 2 {
			fmt.Printf("缩略图尺寸 %s 有误\n", s)
			os.Exit(1)
		}
		w, err := strconv.Atoi(l[0])
		if err != nil {
			fmt.Printf("缩略图尺寸 %s 有误\n", s)
			os.Exit(1)
		}
		h, err := strconv.Atoi(l[1])
		if err != nil {
			log.Fatalf("缩略图尺寸 %s 有误\n", s)
		}
		sizes = append(sizes, Size{w, h, suffix})
	}
	fmt.Println("sizes:", sizes)
	suffixList := make([]string, 0)
	for _, size := range sizes {
		for _, ext := range allowExts {
			suffixList = append(suffixList, size.Suffix+ext)
		}
	}

	filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if info.Size() == 0 {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if !inStringSlice(ext, allowExts) {
			return nil
		}
		//如果是缩略图则跳过，避免生成缩略图的缩略图
		for _, s := range suffixList {
			if strings.HasSuffix(path, s) {
				return nil
			}
		}
		for _, size := range sizes {
			if err := thumb(path, size); err != nil {
				fmt.Printf("生成%s的缩略图失败:%s\n", path, err.Error())
				return nil
			}
		}
		return nil
	})
}

//thumb 生成缩略图
func thumb(filename string, size Size) error {
	var err error
	ext := strings.ToLower(filepath.Ext(filename))
	objPath := filepath.Join(filepath.Dir(filename), strings.TrimSuffix(filepath.Base(filename), ext)+size.Suffix+ext)
	//如果缩略图已存在，并且不覆盖，则直接退出
	if goutil.IsExist(objPath) && !override {
		return nil
	}
	img, err := imaging.Open(filename)
	if err != nil {
		return err
	}
	thumb := imaging.Fit(img, size.Width, size.Height, imaging.NearestNeighbor)
	fmt.Println("生成缩略图:", objPath)
	return imaging.Save(thumb, objPath)
}

func inStringSlice(s string, l []string) bool {
	for _, a := range l {
		if s == a {
			return true
		}
	}
	return false
}
