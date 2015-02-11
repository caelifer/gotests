package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

func main() {
	tests := []string{
		"0",
		"0,1,2",
		"0-10,4-5",
		"0-2,4,6-8",
		"A",
		"-1,-2",
		"1,2,",
		"1 - 2",
		" ",
		"",
	}
	for _, t := range tests {
		fmt.Println(ParseComplexRangeExpr(t).eval())
	}
}

type ComplexRange struct {
	ranges []*Range
}

func (c *ComplexRange) eval() (res []int) {
	if c == nil {
		return nil
	}

	for _, v := range c.ranges {
		res = append(res, v.eval()...)
	}

	sort.Ints(res)

	distinct(res)

	return
}

func distinct(r []int) {
	if len(r) > 0 {
		filter := r[:1] // Reuse space
		filter[0] = r[0]

		old := r[0]
		for _, v := range r {
			// inplace
			if v != old {
				old = v
				filter = append(filter, v)
			}
		}
		r = filter
	}
}

func ParseComplexRangeExpr(expr string) *ComplexRange {
	var rgs []*Range
	for _, v := range strings.Split(expr, ",") {
		if r := ParseRangeExpr(v); r != nil {
			rgs = append(rgs, r)
		} else {
			log.Println("bad range", v, expr)
		}
	}

	if len(rgs) > 0 {
		return &ComplexRange{rgs}
	}
	return nil
}

type Range struct {
	ranges []int
}

func (r Range) eval() []int {
	return r.ranges
}

func ParseRangeExpr(expr string) *Range {
	var vals []int
	parts := strings.Split(expr, "-")

	switch l := len(parts); l {
	case 1:
		if ival1, err := strconv.Atoi(parts[0]); err == nil {
			// Deal with one
			// number range
			vals = append(vals, ival1)
		} else {
			log.Printf("bad number %q in sub-range %q - %v", parts[0], expr, err)
		}
	case 2:
		if ival1, err := strconv.Atoi(parts[0]); err == nil {
			if ival2, err := strconv.Atoi(parts[1]); err == nil {
				for i := ival1; i <= ival2; i++ {
					vals = append(vals, i)
				}
			} else {
				log.Printf("bad number[2] %q in sub-range %q - %v", parts[1], expr, err)
			}
		} else {
			log.Printf("bad number[1] %q in sub-range %q - %v", parts[0], expr, err)
		}
	default:
		return nil
	}

	return &Range{vals}
}
