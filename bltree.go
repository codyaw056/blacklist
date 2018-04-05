package blacklist

import (
    "crypto/sha256"
    "bytes"
    "unicode/utf8"
)

const workFactor = 1
const chunkSize = 2
const increase = 4


// Type BLTree represents a black list hash tree.
type BLTree struct {
    Hash []byte
    WorkFactor int
    Children []*BLTree
}

// Get the index of the child whose hash matches the given hash.
func (b *BLTree) index(hash []byte) int {
    for i := range b.Children {
        if bytes.Equal(hash, b.Children[i].Hash) {
            return i
        }
    }

    return -1
}

// Add a new url to the tree.
func (b *BLTree) Add(url []string) {
    if len(url) == 0 {
        return
    }

    hash := getHash(url[0], b.WorkFactor)
    i := b.index(hash)

    switch {
    case i == -1:
        child := new(BLTree)
        child.Hash = hash
        child.WorkFactor = increase * b.WorkFactor
        child.Add(url[1:])

        b.Children = append(b.Children, child)
    default:
        b.Children[i].Add(url[1:])
    }
}

// Find the given url in the tree.
func (b *BLTree) Find(url []string) bool {
    if len(url) == 0 {
        return true
    }

    hash := getHash(url[0], b.WorkFactor)
    i := b.index(hash)

    switch {
    case i == -1:
        return false
    default:
        return b.Children[i].Find(url[1:])
    }
}

func getHash(data string, factor int) []byte {
    var digest = []byte(data)

    for i:=0; i<factor; i++ {
        sha := sha256.Sum256(digest)
        digest = sha[:]
    }

    return digest
}

func reverse(url string) string {
    runes := make([]rune, utf8.RuneCountInString(url))
    length := len(runes)

    for i, c := range url {
        runes[length-(i+1)] = c
    }

    return string(runes)
}

func pad(str string) string {
    for {
        if len(str) % chunkSize == 0 {
            break
        }

        str = str + "="
    }

    return str

}

func chunk(str string) []string {
    var chunks []string

    for i:=0; i<len(str); i += chunkSize {
        chunks = append(chunks, str[i:i+chunkSize])
    }

    return chunks
}

func NewTree(urls []string) *BLTree {
    tree := new(BLTree)
    tree.WorkFactor = workFactor

    for _, url := range urls {
        url = reverse(url)
        url = pad(url)
        chunks := chunk(url)

        tree.Add(chunks)
    }

    return tree
}

func Lookup(tree *BLTree, url string) bool {
    url = reverse(url)
    url = pad(url)
    chunks := chunk(url)

    return tree.Find(chunks)
}
