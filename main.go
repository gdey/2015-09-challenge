/*
This package implements the challenge for SD Gophers.

The general algorithm for this is going to be the following:

We are going to only do this for 2x2 or bigger matrix. For things less then
that we will treat it as a special case.


we will start with the following matrix:

  3 2
  4 1

If it's a 2x2 we are done.
For a 3 x 3 we will do the following:
First add a new row

   3 2
   4 1
   5 6  // This is the new row added.

Next we will flip row and columns.

   5 4 3
   6 1 2

Next we add another row.

   5 4 3
   6 1 2
   7 8 9

Now we rotate the matrix, so that the largest number
is in the last row of the first column. This way we can
repeate the steps; if we need a larger matrix.

   7 6 5
   8 1 4
   9 2 3

Now we print the Matrix.

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
)

var dflag = flag.Uint("d", 2, "The dimentions of the square.")
var sflag = flag.Uint("s", 1, "The dimentions of the square.")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

const (
	StartSize = 2
)

type Sprial struct {
	TargetSize  int
	StartOffset uint
	s           [][]uint
}

func (s *Sprial) Init() {
	if s.TargetSize < StartSize {
		s.TargetSize = StartSize
	}
	m := make([][]uint, StartSize, s.TargetSize)
	for i := 0; i < len(m); i++ {
		m[i] = make([]uint, StartSize, s.TargetSize)
	}
	m[0][0] = 3 + s.StartOffset
	m[0][1] = 2 + s.StartOffset
	m[1][0] = 4 + s.StartOffset
	m[1][1] = 1 + s.StartOffset
	s.s = m
}

func (s Sprial) largestValue() uint {
	// The largest value in the sprial will always be kept in first column of
	// the last row.
	return s.s[len(s.s)-1][0]
}

// addNextRow add next next row of numbers to the bottom of the sprial. This function assumes that the
// largest number is in the first column of the last row.
func (s *Sprial) addNextRow() {
	if s.s == nil || len(s.s) == 0 {
		s.Init()
	}

	start := s.largestValue() + 1 // Starting value
	row := make([]uint, len(s.s[0]), s.TargetSize)
	for i := 0; i < len(s.s[0]); i++ {
		row[i] = uint(i) + start
	}
	s.s = append(s.s, row)
}

// rotate roates the sprial so that the largest number is in the first column of the last row.
func (s *Sprial) rotate() {
	// Make a new sprial we are going to copy the number from the original into the correct spots
	// into this new sprial. this new sprial will be N x M where the original sprial was M x N.
	ns := make([][]uint, len(s.s[0]), s.TargetSize)
	for i := 0; i < len(ns); i++ {
		ns[i] = make([]uint, len(s.s), s.TargetSize)
	}
	slen := len(s.s) - 1
	for i := 0; i < len(s.s[0]); i++ {
		for j := slen; j >= 0; j-- {
			ns[i][slen-j] = s.s[j][i]
		}
	}
	s.s = ns
}

func (s *Sprial) AddNextRow() {
	s.addNextRow()
	// We want to leave it so that the m[0][len(m[0])-1] is the largest number in the spiral. This way we can call addNextRow without worrying about the orintation of the square.
	s.rotate()
}

func (s Sprial) String() string {
	if s.TargetSize < 2 {
		if s.TargetSize == 1 {
			return "1"
		}
		return ""
	}
	if s.s == nil {
		s.Init()
		for i := 0; i < s.TargetSize-StartSize; i++ {
			s.AddNextRow()
			s.AddNextRow()
		}
	}
	str := strconv.Itoa(int(s.largestValue()))
	format := fmt.Sprintf("%%0.%dd ", len(str))
	str = ""
	for _, r := range s.s {
		for _, c := range r {
			str += fmt.Sprintf(format, c)
		}
		str += fmt.Sprintf("\n")
	}
	return str
}

func main() {
	flag.Parse()
	var offset uint
	if *sflag > 0 {
		offset = *sflag - 1
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	s := Sprial{TargetSize: int(*dflag), StartOffset: offset}
	fmt.Printf("%v\n", s)
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}
