package main

import (
	"fmt"
	"sort"
	"strings"
)

var dictionary []string = []string{
	"coal",
	"cool",
	"cull",
	"battery",
}

func isWord(dict []string, str string) bool {
	for _, v := range dict {
		if strings.HasPrefix(str, v) {
			return true
		}
	}
	return false
}

func seen(path [][2]int, val [2]int) bool {
	for _, v := range path {
		if v == val {
			return true
		}
	}
	return false
}

func appendIfNotSeen(path [][2]int, neighbors [][2]int, val [2]int) [][2]int {
	if !seen(path, val) {
		return append(neighbors, val)
	}
	return neighbors
}

// fmt.Println(getNeighbors([2]int{3, 3}, [][2]int{{0, 0}, {1, 0}, {2, 0}, {1, 1}}))
func getNeighbors(dim [2]int, path [][2]int) [][2]int {
	curr := path[len(path)-1]
	x, y := curr[0], curr[1]
	ns := [][2]int{}
	// go through the different cases
	if y > 0 {
		if x > 0 {
			ns = appendIfNotSeen(path, ns, [2]int{x - 1, y - 1})
		}
		ns = appendIfNotSeen(path, ns, [2]int{x, y - 1})
		if x < dim[0]-1 {
			ns = appendIfNotSeen(path, ns, [2]int{x + 1, y - 1})
		}
	}

	if x > 0 {
		ns = appendIfNotSeen(path, ns, [2]int{x - 1, y})
	}

	if x < dim[0]-1 {
		ns = appendIfNotSeen(path, ns, [2]int{x + 1, y})
	}

	if y < dim[1]-1 {
		if x > 0 {
			ns = appendIfNotSeen(path, ns, [2]int{x - 1, y + 1})
		}
		ns = appendIfNotSeen(path, ns, [2]int{x, y + 1})
		if x < dim[0]-1 {
			ns = appendIfNotSeen(path, ns, [2]int{x + 1, y + 1})
		}
	}

	return ns
}

func copyPath(path [][2]int) [][2]int {
	newPath := make([][2]int, len(path))
	for i := range path {
		newPath[i] = path[i]
	}
	return newPath
}

func nextPaths(dims [2]int, path [][2]int) [][][2]int {
	var neighbors [][2]int = getNeighbors(dims, path)
	if len(neighbors) == 0 {
		return [][][2]int{path}
	}

	var newPaths [][][2]int = [][][2]int{}
	for _, n := range neighbors {
		newPath := append(copyPath(path), n)
		nps := nextPaths(dims, newPath)
		for _, np := range nps {
			newPaths = append(newPaths, np)
		}
	}
	return newPaths
}

func pathToString(board [][]rune, path [][2]int) string {
	runes := []rune{}
	for _, xy := range path {
		runes = append(runes, board[xy[0]][xy[1]])
	}
	return string(runes)
}

func hasDupes(path [][2]int) bool {
	var set = map[string]bool{}
	for _, xy := range path {
		key := fmt.Sprint(xy)
		if !set[key] {
			set[key] = true
		} else {
			return true
		}
	}
	return false
}

// [c o] // [l, o]
//
//	path(x, y) ->  (x, y) + path(neighbors(x, y))
func findWords(board [][]rune) []string {
	output := []string{}
	for i, row := range board {
		for j := range row {
			nextPaths := nextPaths([2]int{len(row), len(board)}, [][2]int{{i, j}})
			for _, p := range nextPaths {
				str := pathToString(board, p)
				if isWord(dictionary, str) {
					output = append(output, pathToString(board, p))
				}
			}
		}
	}
	return output
}

func crushWords(board [][]rune) []string {
	output := []string{}
	nps := nextPaths([2]int{len(board[0]), len(board)}, [][2]int{{0, 0}})
	fmt.Printf("Next paths %v\n", nps)
	for _, p := range nps {
		str := pathToString(board, p)
		fmt.Println(str)
		//if isWord(dictionary, str) {
		output = append(output, pathToString(board, p))
		//}
	}
	return output
}

func eq(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	a := make([]string, len(x))
	copy(a, x)
	b := make([]string, len(y))
	copy(b, y)
	sort.Strings(a)
	sort.Strings(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func BoggleRun() {

	a := []string{"a", "c", "b", "d"}
	b := []string{"d", "c", "b", "a"}
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(eq(a, b))
	fmt.Println(eq([]string{"a", "c", "b", "d"}, []string{"d", "c", "b", "a"}))
	fmt.Println(a)
	fmt.Println(b)

	var simpleBoard [][]rune = [][]rune{
		{'o', 'x', 'x'},
		{'o', 'c', 'x'},
		{'x', 'l', 'x'},
	}

	//ns := [][2]int{}
	//fmt.Println(appendIfNotSeen([][2]int{{1, 2}, {1, 3}, {2, 3}}, ns, [2]int{3, 3}))
	//fmt.Println(appendIfNotSeen([][2]int{{1, 2}, {1, 3}, {2, 3}}, ns, [2]int{1, 3}))
	//fmt.Println(getNeighbors([2]int{3, 3}, [][2]int{{1, 2}, {0, 2}})) // [][2]int {
	fmt.Println("hi")
	fmt.Println(getNeighbors([2]int{3, 3}, [][2]int{{0, 0}, {1, 0}, {2, 0}, {1, 1}}))
	fmt.Println("hi")
	//fmt.Println(getNeighbors([2]int{3, 3}, [][2]int{{0, 0}}))
	//[[0 0] [1 0] [2 0] [1 1] [1 0] [0 1] [0 0]]

	fmt.Println(findWords(simpleBoard))

	fmt.Println("Hello World!")
}
