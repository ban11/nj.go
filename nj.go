package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type DM struct {
	mat  [][]float64
	name []string
}

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func (dm *DM) Generate_matrix() {
	file, err := os.Open("out.csv")
	failOnError(err)
	defer file.Close()

	reader := csv.NewReader(file)

	for i := 0; ; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}
		if i == 0 {
			dm.mat = make([][]float64, len(record))
			for a := 0; a < len(record); a++ {
				dm.mat[a] = make([]float64, len(record))
			}
		}

		for j, s := range record {
			dm.mat[i][j], _ = strconv.ParseFloat(s, 64)
		}

	}
}

func (dm *DM) Generate_namelist() {
	file, err := os.Open("out_name.csv")

	failOnError(err)
	defer file.Close()

	reader := csv.NewReader(file)

	for i := 0; ; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}
		if i == 0 {
			dm.name = []string{}
		}
		dm.name = append(dm.name, record[0])
	}
}

func (dm *DM) Row_sum(x int) float64 {
	nseq := len(dm.mat)
	sum := 0.0

	for i := 0; i < nseq; i++ {
		sum += dm.mat[x][i]
	}

	return sum
}

func (dm *DM) Delete(y int) {
	nseq := len(dm.mat)
	result := make([][]float64, 0, nseq-1)

	for i := 0; i < nseq; i++ {
		if i != y {
			result = append(result, dm.mat[i])
		}
	}

	result2 := make([][]float64, nseq-1)
	for i := 0; i < nseq-1; i++ {
		result2[i] = make([]float64, 0, nseq-1)
	}

	for i := 0; i < nseq-1; i++ {
		for j := 0; j < nseq; j++ {
			if j != y {
				result2[i] = append(result2[i], result[i][j])
			}
		}
	}

	dm.mat = result2
}

func (dm *DM) Minimum_element() (int, int) {
	nseq := len(dm.mat)
	var (
		d, min_d     float64
		min_i, min_j int
	)

	for i := 0; i < nseq; i++ {
		for j := i + 1; j < nseq; j++ {
			d = float64(nseq-2)*dm.mat[i][j] - dm.Row_sum(i) - dm.Row_sum(j)

			if d < min_d {
				min_d = d
				min_i = i
				min_j = j
			}
		}
	}
	return min_i, min_j
}

func (dm *DM) Update(x, y int) {
	nseq := len(dm.mat)

	for i := 0; i < nseq; i++ {
		if i == y {
			continue
		}

		dm.mat[x][i] = (dm.mat[x][i] + dm.mat[y][i] - dm.mat[x][y]) / 2
		dm.mat[i][x] = dm.mat[x][i]
	}
	dm.Delete(y)
}

func (dm *DM) Rename(x, y int) {
	dm.name[x] = "(" + dm.name[x] + "," + dm.name[y] + ")"
	new_name := []string{}
	for i, s := range dm.name {
		if i != y {
			new_name = append(new_name, s)
		}
	}

	dm.name = new_name
	return
}

func (dm *DM) Nj() string {
	nseq := len(dm.mat)
	var min_i, min_j int

	for i := 0; i < nseq-2; i++ {
		min_i, min_j = dm.Minimum_element()
		dm.Update(min_i, min_j)
		dm.Rename(min_i, min_j)
		fmt.Println(i)
	}

	newick := "(" + dm.name[0] + "," + dm.name[1] + ");"
	return newick
}

func main() {
	var dm DM
	dm.Generate_matrix()
	dm.Generate_namelist()
	fmt.Println(dm.mat)
	fmt.Println(dm.name)

	newick := dm.Nj()
	fmt.Println(newick)
	file, err := os.Create("go_newick")
	failOnError(err)
	defer file.Close()

	file.Write(([]byte)(newick))
}
