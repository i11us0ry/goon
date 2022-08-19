package public

import (
	"github.com/axgle/mahonia"
	"golang.org/x/net/html/charset"
	"regexp"
	"strings"
)

/**
 * 对外公开的编码转换接口，传入的字符串会自动检测编码，并转换成utf-8
 */
func ToUtf8(content string) string {
	return toUtf8(content, "")
}


func TitletoUtf8(content string, title string) string {
	var contentType string
	var htmlEncode string
	var htmlEncode2 string
	var htmlEncode3 string
	/* 通过charset获取网站编码 */
	reg := regexp.MustCompile(`(?is)<meta[^>]*charset\s*=["']?\s*([A-Za-z0-9\-]+)`)
	match := reg.FindStringSubmatch(content)
	if len(match) > 1 {
		contentType = strings.ToLower(match[1])
		if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
			htmlEncode2 = "gb18030"
		} else if strings.Contains(contentType, "big5") {
			htmlEncode2 = "big5"
		} else if strings.Contains(contentType, "utf-8") {
			htmlEncode2 = "utf-8"
		}
	}

	_, contentType, _ = charset.DetermineEncoding([]byte(title), "")
	contentType = strings.ToLower(contentType)
	if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
		htmlEncode3 = "gb18030"
	} else if strings.Contains(contentType, "big5") {
		htmlEncode3 = "big5"
	} else if strings.Contains(contentType, "utf-8") {
		htmlEncode3 = "utf-8"
	}

	if htmlEncode3 != "" && htmlEncode2 != htmlEncode3 {
		htmlEncode2 = htmlEncode3
	}
	if htmlEncode2 != "" && htmlEncode != htmlEncode2 {
		htmlEncode = htmlEncode2
	}
	if htmlEncode != "" && htmlEncode != "utf-8" {
		/* 编码为utf-8 */
		title = Convert(title, htmlEncode, "utf-8")
	}
	return title
}
/**
 * 内部编码判断和转换，会自动判断传入的字符串编码，并将它转换成utf-8
 * windows-1252 并不是一个具体的编码，直接拿它来转码会失败
 */
func toUtf8(content string, contentType string) string {

	var htmlEncode string
	var htmlEncode2 string
	var htmlEncode3 string

	/* 指定初始编码 */
	if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
		htmlEncode = "gb18030"
	} else if strings.Contains(contentType, "big5") {
		htmlEncode = "big5"
	} else if strings.Contains(contentType, "utf-8") {
		htmlEncode = "utf-8"
	}

	/* 通过charset获取网站编码 */
	reg := regexp.MustCompile(`(?is)<meta[^>]*charset\s*=["']?\s*([A-Za-z0-9\-]+)`)
	match := reg.FindStringSubmatch(content)
	if len(match) > 1 {
		contentType = strings.ToLower(match[1])
		if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
			htmlEncode2 = "gb18030"
		} else if strings.Contains(contentType, "big5") {
			htmlEncode2 = "big5"
		} else if strings.Contains(contentType, "utf-8") {
			htmlEncode2 = "utf-8"
		}
	}

	/* 通过DetermineEncoding()获取网站title编码 */
	reg = regexp.MustCompile(`(?is)<title[^>]*>(.*?)<\/title>`)
	match = reg.FindStringSubmatch(content)
	reg1 := regexp.MustCompile(`document\.title[\s]*=[\s]*['\"](.*?)['\"]`)
	match1 := reg1.FindStringSubmatch(content)
	if len(match)<1 && len(match)>1{
		match = match1
	}
	if len(match) > 1 {
		aa := match[1]
		/* 判断网页标题编码 */
		_, contentType, _ = charset.DetermineEncoding([]byte(aa), "")
		contentType = strings.ToLower(contentType)
		if strings.Contains(contentType, "gbk") || strings.Contains(contentType, "gb2312") || strings.Contains(contentType, "gb18030") || strings.Contains(contentType, "windows-1252") {
			htmlEncode3 = "gb18030"
		} else if strings.Contains(contentType, "big5") {
			htmlEncode3 = "big5"
		} else if strings.Contains(contentType, "utf-8") {
			htmlEncode3 = "utf-8"
		}
	}

	/*
		进行编码
		如果htmlEncode3（title）编码和htmlEncode2（charset）不一致，则使用（htmlEncode3）title编码
		如果htmlEncode2（上一步筛选结果） 编码和 htmlEncode 原始编码不一致，则把htmlEncode编码转换为htmlEncode2编码
		如果最终 htmlEncode 编码不为空或不为utf-8，则转换为utf-8编码
	*/
	if htmlEncode3 != "" && htmlEncode2 != htmlEncode3 {
		htmlEncode2 = htmlEncode3
	}
	if htmlEncode2 != "" && htmlEncode != htmlEncode2 {
		htmlEncode = htmlEncode2
	}
	if htmlEncode != "" && htmlEncode != "utf-8" {
		/*  */
		content = Convert(content, htmlEncode, "utf-8")
	}

	return content
}

/**
 * 编码转换
 * 需要传入原始编码和输出编码，如果原始编码传入出错，则转换出来的文本会乱码
 */
func Convert(src string, srcCode string, tagCode string) string {
	if srcCode == tagCode {
		return src
	}
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}