package key

import (
	"encoding/hex"
	"fmt"
	"sync"
	"testing"
)

func TestCreateKey(t *testing.T) {
	add,pri,err := CreateKey()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(add)
	fmt.Println(pri)
}

func TestAddress_Hex(t *testing.T) {
	//我们创建一个Pool，并实现New()函数
	sp := sync.Pool{
		//New()函数的作用是当我们从Pool中Get()对象时，如果Pool为空，则先通过New创建一个对象，插入Pool中，然后返回对象。
		New: func() interface{} {
			return make([]int, 16);
		},
	};
	item := sp.Get();
	//打印可以看到，我们通过New返回的大小为16的[]int
	fmt.Println("item : ", item);

	//然后我们对item进行操作
	//New()返回的是interface{}，我们需要通过类型断言来转换
	for i := 0; i < len(item.([]int)); i++ {
		item.([]int)[i] = i;
	}
	fmt.Println("item : ", item);

	//使用完后，我们把item放回池中，让对象可以重用
	sp.Put(item);

	//再次从池中获取对象
	item2 := sp.Get();
	//注意这里获取的对象就是上面我们放回池中的对象
	fmt.Println("item2 : ", item2);
	//我们再次获取对象
	item3 := sp.Get();
	//因为池中的对象已经没有了，所以又重新通过New()创建一个新对象，放入池中，然后返回
	//所以item3是大小为16的空[]int
	fmt.Println("item3 : ", item3);
}

//0xdce75625e4267CE1d15ab8EBa69811cDA23101F4
//0xdce75625e4267ce1d15ab8eba69811cda23101f4
//2e911c236ee73f4b26584cf5e48ff6f9fb9f6645ce9609f3192c6bbb6ba50566

func TestGetAddFromPri(t *testing.T) {
	add,err := GetAddFromPri("2e911c236ee73f4b26584cf5e48ff6f9fb9f6645ce9609f3192c6bbb6ba50566")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(add)
}
///{\"token\":\"hlc.t\",\"amount\":10000.0,\"from\":\"btcadmin\",\"to\":\"mmadmin\"}
func TestSignInfo(t *testing.T) {
	sig,err := SignInfo("EHB","2e911c236ee73f4b26584cf5e48ff6f9fb9f6645ce9609f3192c6bbb6ba50566")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hex.EncodeToString(sig))
}

func TestSigToPub(t *testing.T) {



	add,err := SignToAddress("{\"token\":\"hlc.t\",\"amount\":10000.0,\"from\":\"btcadmin\",\"to\":\"mmadmin\"}","02fe3d8c027e48ab57d5ee36f27659559dc11de804c13df720ddf89ea5cf7303124170ceb392cfb6f08d3c0777d4d7d8852b55ce7c9a59a1de936dc88e403aa300")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(add)
}

