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
		fmt.Printf("%-10q - %+v\n\n", t, ParseComplexRangeExpr(t).Eval())
	}
}

type ComplexRange struct {
	ranges []*Range
}

func (c *ComplexRange) Eval() (res []int) {
	if c == nil {
		return nil
	}

	for _, v := range c.ranges {
		res = append(res, v.Eval()...)
	}

	sort.Ints(res)

	res = distinct(res)

	return
}

func distinct(r []int) []int {
	if len(r) > 0 {
		filter := r[:1]  // Reuse space
		filter[0] = r[0] // Always accept first number

		old := r[0]
		for _, v := range r {
			// in place
			if v != old {
				old = v
				filter = append(filter, v)
			}
		}
		return filter
	}
	return r
}

func ParseComplexRangeExpr(expr string) *ComplexRange {
	var rgs []*Range
	for _, v := range strings.Split(expr, ",") {
		if r := ParseRangeExpr(v); r != nil {
			rgs = append(rgs, r)
		} else {
			log.Printf("bad range %q %q\n", v, expr)
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

func (r Range) Eval() []int {
	return r.ranges
}

func ParseRangeExpr(expr string) *Range {
	var vals []int
	parts := strings.Split(expr, "-")
	nparts := len(parts)

	if nparts == 0 {
		return nil
	}

	var (
		ival1 int
		err   error
	)

	if ival1, err = strconv.Atoi(parts[0]); err == nil {
		// Deal with one number range
		vals = append(vals, ival1)
	} else {
		log.Printf("bad number %q in sub-range %q - %v", parts[0], expr, err)
	}

	if nparts == 2 {
		if ival2, err := strconv.Atoi(parts[1]); err == nil {
			for i := ival1; i <= ival2; i++ {
				vals = append(vals, i)
			}
		} else {
			log.Printf("bad number[2] %q in sub-range %q - %v", parts[1], expr, err)
		}
	}

	return &Range{vals}
}
