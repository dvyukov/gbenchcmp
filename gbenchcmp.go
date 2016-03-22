package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	flagMode = flag.String("mode", "last", "select last/best results")
	flagCpu  = flag.Bool("cpu", true, "show CPU time")
	flagLoad = flag.Bool("load", false, "show CPU load")
)

type Res struct {
	Name string
	Time []uint64
	Cpu  []uint64
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Fprintf(os.Stderr, "usage: benchcmpcc [-flags] old.txt new.txt\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *flagMode != "last" && *flagMode != "best" {
		fmt.Fprintf(os.Stderr, "flag -mode has bad value: %v, expect: last/best\n", *flagMode)
		flag.PrintDefaults()
		os.Exit(1)
	}
	res0 := parse(flag.Args()[0])
	res1 := parse(flag.Args()[1])
	map1 := make(map[string]*Res)
	for _, r := range res1 {
		map1[r.Name] = r
	}
	var data [][]string
	data = append(data, []string{"Benchmark", "Time(ns): old", "new", "diff"})
	if *flagCpu {
		data[0] = append(data[0], "CPU(ns): old", "new", "diff")
	}
	if *flagLoad {
		data[0] = append(data[0], "Load(%): old", "new", "diff")
	}
	for _, r0 := range res0 {
		r1, ok := map1[r0.Name]
		if !ok {
			continue
		}
		time0, cpu0 := choose(r0, *flagMode)
		time1, cpu1 := choose(r1, *flagMode)
		vals := []string{
			r0.Name,
			fmt.Sprintf("%v", time0),
			fmt.Sprintf("%v", time1),
			fmt.Sprintf("%+.2f%%", diff(time0, time1)),
		}
		if *flagCpu {
			vals = append(vals,
				fmt.Sprintf("%v", cpu0),
				fmt.Sprintf("%v", cpu1),
				fmt.Sprintf("%+.2f%%", diff(cpu0, cpu1)),
			)
		}
		if *flagLoad {
			load0 := float64(cpu0) / float64(time0)
			load1 := float64(cpu1) / float64(time1)
			vals = append(vals,
				fmt.Sprintf("%.1f", load0),
				fmt.Sprintf("%.1f", load1),
				fmt.Sprintf("%+.2f%%", difff(load0, load1)),
			)
		}
		data = append(data, vals)
	}
	const ws = 2
	lineWidth := -ws
	var width []int
	for i := range data[0] {
		w := 0
		for _, row := range data {
			if w < len(row[i]) {
				w = len(row[i])
			}
		}
		width = append(width, w)
		lineWidth += w + ws
	}
	for ri, row := range data {
		for i, s := range row {
			x := width[i] - len(s)
			if i != 0 {
				x += ws
			}
			pad := strings.Repeat(" ", x)
			if i != 0 {
				s, pad = pad, s
			}
			fmt.Printf("%v%v", s, pad)
		}
		fmt.Printf("\n")
		if ri == 0 {
			fmt.Printf("%v\n", strings.Repeat("=", lineWidth))
		}
	}
}

func choose(r *Res, mode string) (time, cpu uint64) {
	switch mode {
	case "last":
		return r.Time[len(r.Time)-1], r.Cpu[len(r.Cpu)-1]
	case "best":
		time = ^uint64(0)
		for i, t := range r.Time {
			if time > t {
				time = t
				cpu = r.Cpu[i]
			}
		}
		return
	default:
		panic("flags should have been verified before")
	}
}

func diff(v0, v1 uint64) float64 {
	return difff(float64(v0), float64(v1))
}

func difff(v0, v1 float64) float64 {
	return v1/v0*100 - 100
}

func parse(fn string) []*Res {
	f, err := os.Open(fn)
	if err != nil {
		failf("failed to open input file '%v': %v", fn, err)
	}
	defer f.Close()
	var res []*Res
	dups := make(map[string]*Res)
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if !strings.HasPrefix(line, "BM_") {
			continue
		}
		for strings.Contains(line, "  ") {
			line = strings.Replace(line, "  ", " ", -1)
		}
		parts := strings.Split(line, " ")
		if len(parts) < 3 {
			continue
		}
		time, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseUint(parts[2], 10, 64)
		if err != nil {
			continue
		}
		name := parts[0][3:]
		r := &Res{Name: name}
		if prev := dups[name]; prev != nil {
			r = prev
		} else {
			dups[name] = r
			res = append(res, r)
		}
		r.Time = append(r.Time, time)
		r.Cpu = append(r.Cpu, cpu)
	}
	if err := s.Err(); err != nil {
		failf("failed to read input file '%v': %v", fn, err)
	}
	return res
}

func failf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
