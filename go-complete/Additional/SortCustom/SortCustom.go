package main

import "sort"
import "fmt"

type byLength []string

func (s byLength) Len() int {
	return len(s)
}

func (s byLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func main() {
	fruits := []string{"peach", "bananabigger", "k"}
	sort.Sort(byLength(fruits))
	fmt.Println(fruits)
}
