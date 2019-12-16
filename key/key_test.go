package key

import (
	"encoding/hex"
	"fmt"
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
	sig,err := SignInfo("{\"token\":\"usaapl.t\",\"amount\":100000.0}","2e911c236ee73f4b26584cf5e48ff6f9fb9f6645ce9609f3192c6bbb6ba50566")
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

