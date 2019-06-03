# Consisten 一致性哈希算法

## 介绍

在分布式中，一致性HASH算法有着很多的应用场景，在避免机器缩容、扩容时都有着优良的rehash结果。 减少移动数据的绝对值数量，也减少热机器的压力，将待移动key均匀的分布到多台机器。

[一致性HASH算法原理](https://mp.weixin.qq.com/s/cV75gXnZhycWneD6TthBjg)

使用跳跃链表的结构来存储一致性HASH中的数据环，在hash值查找过程时有着优秀的性能。 

## Example

``` go
  import "zhaozy93/consisten" 
  // NewConsistenObject(circleSum, repeatNum uint){}
  // repeatNum 表示单个元素在环上重复出现次数:  1 << repeatNum -1
  // circleSum 表示环上元素总数:  1<< circleSum 
  c := NewConsistenObject(10, 4)
	c.Insert("baidu")
	c.Insert("alibaba")
	c.Insert("tencent")

	fmt.Println(c.GetItemInfo())

	fmt.Println("didi will hash to location", c.GetHashValue("didi"))
```


 
## 微信
 
![wechat-group-chat](./assets/wechat_1.jpeg)
![wechat-group-chat](./assets/wechat_2.jpeg)
