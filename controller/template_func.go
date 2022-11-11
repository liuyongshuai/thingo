// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-29 14:37

package controller

import (
	"encoding/json"
	"github.com/liuyongshuai/negoutils/convertutils"
	"html/template"
	"reflect"
	"strings"
	"time"
)

//常用的模板函数
var CommonTplFuncs = map[string]interface{}{
	"substr":                  TplFuncSubStr,
	"htmlspecialchars":        TplFuncHtmlSpecialChars,
	"htmlspecialchars_decode": TplFuncHtmlSpecialcharsDecode,
	"json_encode":             TplFuncJsonEncode,
	"html_quote":              TplFuncHtmlQuote,
	"html_unquote":            TplFuncHtmlUnQuote,
	"str2html":                TplFuncStr2Html,
	"date_format":             TplFuncDateFormat,
	"date_parse":              TplFuncDateParse,
	"date":                    TplFuncDate,
	"eq":                      TplFuncEQ,
	"lt":                      TplFuncLT,
}

var (
	//类似于PHP的日期格式化选项
	datePatterns = []string{
		// year
		"Y", "2006",
		"y", "06",
		// month
		"m", "01",
		"n", "1",
		"M", "Jan",
		"F", "January",
		// day
		"d", "02",
		"j", "2",
		// week
		"D", "Mon",
		"l", "Monday",
		// time
		"g", "3",
		"G", "15",
		"h", "03",
		"H", "15",
		"a", "pm",
		"A", "PM",
		"i", "04",
		"s", "05",
		// time zone
		"T", "MST",
		"P", "-07:00",
		"O", "-0700",
		// RFC 2822
		"r", time.RFC1123Z,
	}
)

//高仿PHP的函数htmlspecialchars
func TplFuncHtmlSpecialChars(html string) template.HTML {
	replace := strings.NewReplacer(
		"&", "&amp;",
		"'", "&apos;",
		"\"", "&quot;",
		">", "&gt;",
		"<", "&lt;",
	)
	html = replace.Replace(html)
	return template.HTML(html)
}

//高仿PHP的htmlspecialchars_decode
func TplFuncHtmlSpecialcharsDecode(str string) template.HTML {
	replace := strings.NewReplacer(
		"&amp;", "&",
		"&apos;", "'",
		"&quot;", "\"",
		"&gt;", ">",
		"&lt;", "<",
	)
	str = replace.Replace(str)
	return template.HTML(str)
}

//json_encode
func TplFuncJsonEncode(t interface{}) template.HTML {
	ret, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return template.HTML(ret)
}

//转义html字符
func TplFuncHtmlQuote(text string) template.HTML {
	text = strings.Replace(text, "&", "&amp;", -1)
	text = strings.Replace(text, "<", "&lt;", -1)
	text = strings.Replace(text, ">", "&gt;", -1)
	text = strings.Replace(text, "'", "&#39;", -1)
	text = strings.Replace(text, "\"", "&quot;", -1)
	text = strings.Replace(text, "“", "&ldquo;", -1)
	text = strings.Replace(text, "”", "&rdquo;", -1)
	text = strings.Replace(text, " ", "&nbsp;", -1)
	return template.HTML(strings.TrimSpace(text))
}

//反转义html字符串
func TplFuncHtmlUnQuote(text string) template.HTML {
	text = strings.Replace(text, "&nbsp;", " ", -1)
	text = strings.Replace(text, "&rdquo;", "”", -1)
	text = strings.Replace(text, "&ldquo;", "“", -1)
	text = strings.Replace(text, "&quot;", "\"", -1)
	text = strings.Replace(text, "&#39;", "'", -1)
	text = strings.Replace(text, "&gt;", ">", -1)
	text = strings.Replace(text, "&lt;", "<", -1)
	text = strings.Replace(text, "&amp;", "&", -1)
	return template.HTML(strings.TrimSpace(text))
}

//字符串转为html类型
func TplFuncStr2Html(raw string) template.HTML {
	return template.HTML(raw)
}

//截取字符串
func TplFuncSubStr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)
	if start < 0 {
		start += length
	} else if start > length {
		start = start % length
	}
	if end < 0 {
		end += length
	} else if end > length {
		end = end % length
	}
	if start > end || end < 0 || start < 0 {
		return ""
	}
	return string(rs[start : end+1])
}

//日期格式化
func TplFuncDateFormat(t time.Time, layout string) template.HTML {
	return template.HTML(t.Format(layout))
}

//高仿PHP的日期解析
func TplFuncDateParse(dateString, format string) (time.Time, error) {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return time.ParseInLocation(format, dateString, time.Local)
}

//高仿PHP的date
func TplFuncDate(timestamp int64, format string) template.HTML {
	t := time.Unix(timestamp, 0)
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return template.HTML(t.Format(format))
}

//相等
func TplFuncEQ(arg1 interface{}, arg2 ...interface{}) (bool, error) {
	v1 := reflect.ValueOf(arg1)
	k1, err := convertutils.GetBasicKind(v1)
	if err != nil {
		return false, err
	}
	if len(arg2) == 0 {
		return false, convertutils.ErrorNoComparison
	}
	for _, arg := range arg2 {
		v2 := reflect.ValueOf(arg)
		k2, err := convertutils.GetBasicKind(v2)
		if err != nil {
			return false, err
		}
		if k1 != k2 {
			return false, convertutils.ErrorBadComparison
		}
		truth := false
		switch k1 {
		case convertutils.BoolKind:
			truth = v1.Bool() == v2.Bool()
		case convertutils.ComplexKind:
			truth = v1.Complex() == v2.Complex()
		case convertutils.FloatKind:
			truth = v1.Float() == v2.Float()
		case convertutils.IntKind:
			truth = v1.Int() == v2.Int()
		case convertutils.StringKind:
			truth = v1.String() == v2.String()
		case convertutils.UintKind:
			truth = v1.Uint() == v2.Uint()
		default:
			return false, convertutils.ErrorInvalidInputType
		}
		if truth {
			return true, nil
		}
	}
	return false, nil
}

//小于
func TplFuncLT(arg1, arg2 interface{}) (bool, error) {
	v1 := reflect.ValueOf(arg1)
	k1, err := convertutils.GetBasicKind(v1)
	if err != nil {
		return false, err
	}
	v2 := reflect.ValueOf(arg2)
	k2, err := convertutils.GetBasicKind(v2)
	if err != nil {
		return false, err
	}
	if k1 != k2 {
		return false, convertutils.ErrorBadComparison
	}
	truth := false
	switch k1 {
	case convertutils.BoolKind, convertutils.ComplexKind:
		return false, convertutils.ErrorBadComparisonType
	case convertutils.FloatKind:
		truth = v1.Float() < v2.Float()
	case convertutils.IntKind:
		truth = v1.Int() < v2.Int()
	case convertutils.StringKind:
		truth = v1.String() < v2.String()
	case convertutils.UintKind:
		truth = v1.Uint() < v2.Uint()
	default:
		return false, convertutils.ErrorInvalidInputType
	}
	return truth, nil
}
