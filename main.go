package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"strings"
)

type ProcessInfo struct {
	PID      int            `json:"pid"`
	PPID     int            `json:"ppid"`
	CMD      string         `json:"cmd"`
	Children []*ProcessInfo `json:"children,omitempty"`
}

func atoi(s string) int {
	out, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		panic(err)
	}
	return out
}

func readCSV(name string) ([][]string, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return csv.NewReader(file).ReadAll()
}

func main() {

	filename := "test.csv"
	flag.Parse()
	if len(flag.Args()) >= 1 {
		filename = flag.Arg(0)
	}
	matrix, err := readCSV(filename)
	if err != nil {
		panic(err)
	}

	infos := make([]ProcessInfo, len(matrix))
	pid2index := map[int]int{}

	//Generate ProcessInfo struct from csv content
	//PID, PPID, CMD
	for i := 0; i < len(matrix); i++ {
		if len(matrix[i]) != 3 {
			panic("csv contain rows that not 3 columns")
		}
		infos[i] = ProcessInfo{
			PID:  atoi(matrix[i][0]),
			PPID: atoi(matrix[i][1]),
			CMD:  strings.TrimSpace(matrix[i][2]),
		}
		pid2index[infos[i].PID] = i
	}

	//Build process tree
	roots := []*ProcessInfo{}
	for i := 0; i < len(infos); i++ {
		parentIndex, ok := pid2index[infos[i].PPID]
		if ok {
			infos[parentIndex].Children = append(infos[parentIndex].Children, &infos[i])
		} else {
			roots = append(roots, &infos[i])
		}
	}

	out, err := json.MarshalIndent(roots, "", "  ")
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(append(out, '\n'))
}
