package sequence

//
//import (
//	"fmt"
//	"math/rand/v2"
//	"strings"
//	"testing"
//
//	"github.com/HaChains/infra/blockchain/types"
//)
//
//const hashFormat = "%s-%d"
//const confirmations = 7
//
//var forks = []string{"A", "B"}
//
//type block struct {
//	height uint64
//	fork   string
//	parent string
//}
//
//func newBlock(height uint64, _fork ...string) *block {
//	parentHeight := height - 1
//
//	fork := forks[rand.IntN(len(forks))]
//	if len(_fork) != 0 {
//		fork = _fork[0]
//	}
//
//	var parent string
//	switch {
//	case height == 0:
//		parent = ""
//	case height == 1:
//		parent = "0"
//	case height%(confirmations+1) == 0:
//		// joint block, select parent from a random fork
//		parentFork := forks[rand.IntN(len(forks))]
//		parent = fmt.Sprintf(hashFormat, parentFork, parentHeight)
//		fork = ""
//	case height%(confirmations+1) == 1:
//		// block after joint, parent must be the joint block
//		parent = fmt.Sprint(height - 1)
//	default:
//		parent = fmt.Sprintf(hashFormat, fork, parentHeight)
//	}
//
//	return &block{
//		height: height,
//		fork:   fork,
//		parent: parent,
//	}
//}
//
//func (b *block) Height() uint64 {
//	return b.height
//}
//func (b *block) Hash() string {
//	if b.fork == "" {
//		return fmt.Sprint(b.height)
//	}
//	return fmt.Sprintf(hashFormat, b.fork, b.height)
//}
//func (b *block) ParentHash() string {
//	return b.parent
//}
//func (b *block) BlockTime() uint64 {
//	return 0
//}
//func (b *block) String() string {
//	return b.Hash()
//}
//
//type fetcher2 struct{}
//
//func (f *fetcher2) GetBlockByHash(hash string) (types.IBlock, error) {
//	return nil, nil
//}
//
//type fakeFetcher struct{}
//
//func (ff *fakeFetcher) GetBlockByHash(hash string) (types.IBlock, error) {
//	var (
//		height uint64
//		fork   string
//	)
//	msg := strings.ReplaceAll(hash, "-", " ")
//	_, err := fmt.Sscanf(msg, "%s %d", &fork, &height)
//	if err != nil {
//		return nil, fmt.Errorf("invalid hash format: %s", hash)
//	}
//	b := newBlock(height, fork)
//	blockInfo := fmt.Sprintf("%s, parent: %s", b, b.ParentHash())
//	fmt.Printf("🌏 fetched block %s by hash: %s\n", blockInfo, hash)
//	return b, nil
//}
//
//func TestSequence(t *testing.T) {
//	ff := &fakeFetcher{}
//	seq, err := New(confirmations, &block{height: 0}, ff)
//	if err != nil {
//		t.Fatal(err)
//	}
//	// insert blocks
//	for h := uint64(1); h <= 30; h++ {
//		b := newBlock(h)
//		fmt.Printf("Appending block: %s, parent: %s\n", b, b.ParentHash())
//		err = seq.Append(b)
//		if err != nil {
//			t.Fatalf("failed to append block %s: %v", b.Hash(), err)
//		}
//		sequence, err := seq.Get()
//		if err != nil {
//			t.Fatalf("failed to get block sequence: %v", err)
//		}
//		fmt.Println("sequence:", sequence)
//	}
//}
