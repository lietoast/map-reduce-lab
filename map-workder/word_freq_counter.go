package map_workder

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/lietoast/map-reduce-lab/input"
	. "github.com/lietoast/map-reduce-lab/pub-sub"
)

// WordFreqCounter 词频统计
type WordFreqCounter struct{}

func (wfc *WordFreqCounter) Map(in input.InputReader) error {
	for value := in.Next(); value != ""; {
		value = removePunctuation(value)
		words := strings.Split(value, " ")
		for _, w := range words {
			EmitIntermediate(w, "1")
		}
	}
	return nil
}

func removePunctuation(input string) string {
	// 正则表达式匹配标点符号和特殊符号
	re := regexp.MustCompile(`[[:punct:]]`)
	// 使用正则表达式替换为空格
	output := re.ReplaceAllString(input, " ")

	// 替换连续的空格为单个空格
	re = regexp.MustCompile(`\s+`)
	output = re.ReplaceAllString(output, " ")

	// 过滤非字母和空格的字符，同时处理中文标点
	filtered := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			return r
		} else if unicode.IsPunct(r) {
			return ' '
		}
		return -1
	}, output)

	return filtered
}
