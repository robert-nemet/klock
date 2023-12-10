package pkg_webhook

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	alg "github.com/robert-nemet/blobs/algorithms"
	"github.com/robert-nemet/blobs/datas"
)

type Evaluator interface {
	IsMatch(rule, value string) (bool, error)
}

type evaluator struct {
	sy   alg.ShuntingYard
	eval alg.Evaluator[bool]
}

// IsMatch implements Evaluator
func (e evaluator) IsMatch(rule string, value string) (bool, error) {
	// transform to array
	infix := prepack(rule)
	infixs := strings.Join(prepare(infix, value), " ")
	po := e.sy.Transform(infixs)
	return e.eval.Evaluate(po)
}

func NewEvaluator() Evaluator {
	ops := map[string]alg.Operation[bool]{
		"&": func(stack datas.Stack[bool]) (bool, error) {
			fop := stack.Pop()
			if fop == nil {
				return false, errors.New("no operand")
			}
			sop := stack.Pop()
			if sop == nil {
				return false, errors.New("no operand")
			}
			return *fop && *sop, nil
		},
		"|": func(stack datas.Stack[bool]) (bool, error) {
			fop := stack.Pop()
			if fop == nil {
				return false, errors.New("no operand")
			}
			sop := stack.Pop()
			if sop == nil {
				return false, errors.New("no operand")
			}
			return *fop || *sop, nil
		},
		"^": func(stack datas.Stack[bool]) (bool, error) {
			op := stack.Pop()
			if op == nil {
				return false, errors.New("no operand")
			}
			return !*op, nil
		},
	}
	sy := alg.NewShuntingYard(map[string]int{"&": 1, "|": 1, "^": 2})
	ev := alg.NewEvaluator(ops, func(input string) (bool, error) {
		return strconv.ParseBool(input)
	})
	return evaluator{
		sy:   sy,
		eval: ev,
	}
}

func prepare(input []string, value string) []string {
	var result = make([]string, len(input))
	for cnt, s := range input {
		if isReserved(s) {
			result[cnt] = s
			continue
		}
		result[cnt] = strconv.FormatBool(s == value)
	}
	return result
}

func isReserved(s string) bool {
	return s == "&" || s == "^" || s == "|" || s == "(" || s == ")"
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
