package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

var (
	table []map[int]bool
	nodes int
	links int
	num   int
)

func setup() {
	// Make network
	table = make([]map[int]bool, num, num)
	for i := range table {
		table[i] = make(map[int]bool)
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())

	// Size of output network
	flag.IntVar(&num, "num", 100, "size of resulting network")

	var (
		m0, m, k int
		beta     float64
		mode     string
	)

	flag.IntVar(&m0, "m0", 10, "initial size of full-mesh network (freescale)")
	flag.IntVar(&m, "m", 5, "number of links added for each new node (freescale)")

	flag.IntVar(&k, "k", 2, "number (even) of initial neighbours (smallworld)")
	flag.Float64Var(&beta, "beta", 0.5, "probability of re-linking (smallworld)")

	flag.StringVar(&mode, "mode", "", "mode (freescale/smallworld)")

	flag.Parse()

	switch mode {
	case "freescale": //freescale

		if m0 == 0 || m == 0 {
			fmt.Println("lol")
			flag.PrintDefaults()
		}

		if m > m0 || m0 > num {
			fmt.Fprintf(os.Stderr, "m0 (%v) should be less than m (%v)", m0, m)
			os.Exit(2)
		}

		ba_populate(m, m0)
	case "smallworld": // smallworld

		if k <= 0 || k != k/2+k/2 {
			fmt.Fprintln(os.Stderr, "K should be positive even")
			os.Exit(2)
		}

		if !(0 <= beta && beta <= 1) {
			fmt.Fprintln(os.Stderr, "Beta should be between 0 and 1")
			os.Exit(2)
		}

		if !(num >= k && float64(k) >= math.Log(float64(num)) && math.Log(float64(num)) >= 1.0) {
			fmt.Fprintf(os.Stderr, "Condition not matched: N (%v) >> K (%v) >> Log(N) (%v) >> 1 \n", num, k, math.Log(float64(num)))
			os.Exit(2)
		}

		fc_populate(k, beta)

	default:
		fmt.Println("Lol")
		flag.PrintDefaults()
	}

	dump()
}

func connect(i, j int) {
	table[i][j] = true
	table[j][i] = true
	links += 1
}

func count(id int) (out int) {
	for _, v := range table[id] {
		if v {
			out += 1
		}
	}

	return
}

func ba_populate(m, m0 int) {
	setup()
	ba_setup(m0)

	counts := make([]int, nodes, num)
	for i, _ := range counts {
		counts[i] = count(i)
	}

	for i := nodes; i < num; i++ {
		var added int
		for added < m { //add exactly M links
			j := int(rand.Int31n(int32(nodes)))

			if !table[i][j] && rand.Int31n(int32(links)) < int32(counts[j]) {
				connect(i, j)
				added += 1
			}
		}
		nodes += 1
		counts = append(counts, added)
	}

}

func disconnect(i, j int) {
	delete(table[i], j)
	delete(table[j], i)
	links -= 1
}

// Create initial nodes
func ba_setup(m0 int) {

	for i := 0; i < m0; i++ {
		for j := 0; j < m0; j++ {
			table[i][j] = true
		}
	}
	nodes = m0
	links = m0 * (m0 - 1) / 2
}

// Create initial links
func fc_setup(k int) {
	for i := 0; i < num; i++ {
		for di := 1; di <= k/2; di += 1 {
			connect(i, (i+di)%num)
			connect(i, (num+i-di)%num)
		}
	}
}

func fc_populate(k int, beta float64) {
	setup()
	fc_setup(k)

	for i, conns := range table {
		for j, _ := range conns {
			if j > i && rand.Float64() < beta {
				disconnect(i, j)
				added := false
				for added == false {
					k := int(rand.Int31n(int32(num)))
					if i != k && !table[i][k] {
						connect(i, k)
						added = true
					}
				}
			}
		}
	}
}

func dump() {
	fmt.Println(num)
	for _, row := range table {
		for i, linked := range row {
			if linked {
				fmt.Printf("%v ", i)
			}
		}
		fmt.Println()
	}

}
