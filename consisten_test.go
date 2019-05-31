package consisten

import (
	"fmt"
	"log"
	"testing"
)

func TestLink(t *testing.T) {
	fmt.Println("test skiplist starte")
	sk := NewSkiplist(10, 1)
	sk.Insert(4, "4")
	sk.Insert(3, "4")
	sk.Insert(2, "4")
	sk.Insert(1, "4")

	for i := 0; i < 1000; i++ {
		sk.Insert(uint64(i), "4")
	}

	log.Println(sk.Search(101))
	log.Println(sk.Search(35))
	log.Println(sk.Search(60))
	sk.Delete(85)
	sk.Print()

	fmt.Println("test skiplist end")

	fmt.Println("test consist start")
	c := NewConsistenObject(10, 4)
	c.Insert("baidu")
	c.Insert("alibaba")
	c.Insert("tencent")

	c.circle.Print()
	fmt.Println(c.GetItemInfo())

	fmt.Println("didi will hash to location", c.GetHashValue("didi"))
	fmt.Println("test consist end")
}
