package consisten

import (
	"errors"
	"fmt"
	"hash/crc64"
	"sync"
	"unsafe"
)

type Consisten struct {
	maxItem    int
	circleSum  int
	repeatNum  int
	circleSlot map[string][]uint64
	circle     *Skiplist
	locker     *sync.Mutex
}

// repeatNum 单个元素 在环上重复出现次数 =  1 << repeatNum -1
// circleSum  环上元素总数 = 1<< circleSum
func NewConsistenObject(circleSum, repeatNum uint) *Consisten {
	size := uint(unsafe.Sizeof(uint(1)))
	if repeatNum >= circleSum {
		circleSum = size*6 - 1
		repeatNum = 4
	} else {
		maxSize := size*8 - 2
		if circleSum > maxSize {
			circleSum = maxSize
		}
		if repeatNum > 6 {
			repeatNum = 6
		}
	}

	sk := NewSkiplist(6, 2)
	return &Consisten{
		maxItem:    1 << (circleSum - repeatNum - 3),
		circleSum:  1<<circleSum - 1,
		repeatNum:  1<<repeatNum - 1,
		circleSlot: make(map[string][]uint64),
		circle:     sk,
		locker:     new(sync.Mutex),
	}
}

var ECMATable = crc64.MakeTable(crc64.ECMA)

func (c *Consisten) Insert(item string) (err error) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if _, ok := c.circleSlot[item]; ok {
		err = errors.New(item + " has existed!")
		return
	}
	if len(c.circleSlot)+1 == c.maxItem {
		err = errors.New("excess max item limit!")
		return
	}

	c.circleSlot[item] = make([]uint64, c.repeatNum)
	cnt := 0
	for i := 0; cnt < c.repeatNum; i++ {
		key := fmt.Sprintf("%d%s", i, item)
		csum := crc64.Checksum([]byte(key), ECMATable)
		haskKey := csum & uint64(c.circleSum)
		err = c.circle.Insert(haskKey, item)
		if err == nil {
			c.circleSlot[item][cnt] = haskKey
			cnt++
			continue
		}
	}
	return
}

func (c *Consisten) Delete(item string) (err error) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if _, ok := c.circleSlot[item]; !ok {
		err = errors.New(item + " has not existed!")
	}
	for _, key := range c.circleSlot[item] {
		c.circle.Delete(key)
	}
	return
}

func (c *Consisten) GetItems() []string {
	s := make([]string, len(c.circleSlot))
	i := 0
	for key, _ := range c.circleSlot {
		s[i] = key
		i++
	}
	return s
}

func (c *Consisten) GetItemInfo() map[string][]uint64 {
	return c.circleSlot
}

func (c *Consisten) GetHashValue(key string) string {
	if len(c.circleSlot) == 0 {
		return ""
	}
	csum := crc64.Checksum([]byte(key), ECMATable)
	haskKey := csum & uint64(c.circleSum)
	node := c.circle.SearchNext(haskKey)
	return node.ItemValue
}
