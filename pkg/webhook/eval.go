package pkg_webhook

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/antonmedv/expr"
)

func isMatch(rule, value string) bool {
	// simple
	if rule == value {
		return true
	}

	//eval
	input := prepack(rule)
	prep_input := prepare(input, value)

	result, err := expr.Eval(strings.Join(prep_input, ""), map[string]interface{}{})
	if err != nil {
		panic(err)
	}
	return result.(bool)
}

func prepare(input []string, value string) []string {
	var result = make([]string, len(input))
	for cnt, s := range input {
		if s == "|" || s == "&" {
			result[cnt] = s + s
			continue
		}
		if s == "(" || s == ")" {
			result[cnt] = s
			continue
		}
		if s == "^" {
			result[cnt] = "!"
			continue
		}
		result[cnt] = strconv.FormatBool(s == value)
	}

	return result
}

// prepack rule string into array
func prepack(s string) []string {
	zp := regexp.MustCompile("[(|&^)]").FindAllStringIndex(s, -1)

	// flat
	var flat []int
	for _, arr := range zp {
		if len(flat) == 0 || arr[0] != flat[len(flat)-1] {
			flat = append(flat, arr...)
		} else {
			flat = append(flat, arr[1])
		}
	}

	//split
	res := make([]string, len(flat)+1)
	p := 0
	for cnt, v := range flat {
		res[cnt] = s[p:v]
		p = v
	}
	res[len(flat)] = s[p:]

	//trim
	if res[0] == "" {
		res = res[1:]
	}
	if res[len(res)-1] == "" {
		res = res[:len(res)-1]
	}
	return res
}
