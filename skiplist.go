package consisten

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

var ERREXISTED = errors.New("elem has existed")

func compareFunc(a, b uint64) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	}
	return -1
}

type Skiplist struct {
	maxLevel int
	locker   *sync.Mutex
	head     *SkiplistNode
	step     int
}

type SkiplistNode struct {
	ItemValue string
	Value     uint64
	Level     int
	Right     []*SkiplistNode // 0层-Level层 每层的右侧node
}

func newNode(lev int, elem uint64) *SkiplistNode {
	right := make([]*SkiplistNode, lev)
	for i := 0; i < lev; i++ {
		right[i] = nil
	}
	return &SkiplistNode{
		Value: elem,
		Level: lev,
		Right: right,
	}
}

func NewSkiplist(levels int, step uint) *Skiplist {
	head := newNode(levels, 0)
	return &Skiplist{
		maxLevel: levels,
		head:     head,
		locker:   new(sync.Mutex),
		step:     1<<uint(step) - 1, // 15
	}
}

func (sk *Skiplist) randomLevel() int {
	level := 1
	for rand.Int()&sk.step == 0 {
		level++
		if level == sk.maxLevel {
			break
		}
	}
	return level
}

func (sk *Skiplist) Insert(elem uint64, key string) error {
	sk.locker.Lock()
	defer sk.locker.Unlock()
	if sk.search(elem) != nil {
		return ERREXISTED
	}

	level := sk.randomLevel()
	node := newNode(level, elem)
	for i := level - 1; i >= 0; i-- {
		head := sk.head
		for {
			if head.Right[i] == nil {
				head.Right[i] = node
				break
			}
			_com := compareFunc(elem, head.Right[i].Value)
			if _com > 0 {
				head = head.Right[i]
				continue
			} else {
				node.Right[i] = head.Right[i]
				head.Right[i] = node
				break
			}
		}
	}
	node.ItemValue = key
	return nil
}

func (sk *Skiplist) Search(elem uint64) bool {
	sk.locker.Lock()
	defer sk.locker.Unlock()
	node := sk.search(elem)
	if node == nil {
		return false
	}
	return true
}

func (sk *Skiplist) search(elem uint64) *SkiplistNode {
	tempNode := sk.head
	for i := sk.maxLevel - 1; i >= 0; i-- {
		for {
			if tempNode.Right[i] == nil {
				break
			}
			_com := compareFunc(elem, tempNode.Right[i].Value)
			if _com == 0 {
				return tempNode.Right[i]
			}
			if _com > 0 {
				tempNode = tempNode.Right[i]
				continue
			}
			if _com < 0 {
				break
			}
		}
	}
	return nil
}

func (sk *Skiplist) Delete(elem uint64) bool {
	sk.locker.Lock()
	defer sk.locker.Unlock()
	node := sk.search(elem)
	if node == nil {
		return false
	}
	sk.delete(node)
	return true
}

func (sk *Skiplist) delete(node *SkiplistNode) {
	tempPreNode := sk.head
	for i := node.Level - 1; i >= 0; i-- {
		for {
			if node != tempPreNode.Right[i] {
				tempPreNode = tempPreNode.Right[i]
				continue
			}
			tempPreNode.Right[i] = node.Right[i]
			break
		}
	}
}

func (sk *Skiplist) SearchNext(elem uint64) *SkiplistNode {
	tempNode := sk.head
	for i := sk.maxLevel - 1; i >= 0; i-- {
		for {
			if tempNode.Right[i] == nil {
				if i == 0 {
					return sk.head.Right[0]
				}
				break
			}
			_com := compareFunc(elem, tempNode.Right[i].Value)
			if _com == 0 {
				return tempNode.Right[i]
			}
			if _com > 0 {
				tempNode = tempNode.Right[i]
				continue
			}
			if _com < 0 {
				if i == 0 {
					return tempNode.Right[0]
				}
				break
			}
		}
	}
	return nil
}

// 测试辅助打印
func (sk *Skiplist) Print() {
	for i := 0; i < sk.maxLevel; i++ {
		fmt.Print("line:", i, "----")
		head := sk.head
		for {
			if head.Right[i] == nil {
				break
			}
			fmt.Print(head.Right[i].Value, "-")
			head = head.Right[i]
		}
		fmt.Print("\n")
	}
}
