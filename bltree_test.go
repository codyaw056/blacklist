package blacklist

import (
    "testing"
    "strings"
    "io/ioutil"
    "fmt"
)

func openList(filename string) []string {
    var items []string

    bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Could not read file: %s\n", filename)
		return items
	}

    data := string(bytes)
	for _, line := range strings.Split(data, "\n") {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		items = append(items, line)
	}

	return items

}

func BenchmarkNewTree(b *testing.B) {
    urls := openList("blacklist.txt")

    for n := 0; n < b.N; n++ {
        NewTree(urls)
    }
}

func BenchmarkLookupFail(b *testing.B) {
    urls := openList("blacklist.txt")
    tree := NewTree(urls)

    for n:=0; n < b.N; n++ {
        Lookup(tree, "yahoo.com")
    }
}

func BenchmarkLookupSuccess(b *testing.B) {
    urls := openList("blacklist.txt")
    tree := NewTree(urls)

    for n:=0; n < b.N; n++ {
        Lookup(tree, "gooogle.com")
    }
}
