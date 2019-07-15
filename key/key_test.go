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
//2e911c236ee73f4b26584cf5e48ff6f9fb9f6645ce9609f3192c6bbb6ba50566

func TestGetAddFromPri(t *testing.T) {
	add,err := GetAddFromPri("6cd5703701af0ed006ef4fddb612d5c0b373aca7c6920564b00b7e1df43d8eea")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(add)
}


//Desc       string `json:"desc"`
//Token      string `json:"token"`
//Sender     string `json:"sender"`
//Receiver   string `json:"receiver"`
//Amount     float64 `json:"amount"`

//{\"token\":\"hlc.t\",\"amount\":1.3,\"sender\":\"btcadmin\",\"receiver\":\"mmadmin\",\"desc\":\"desc1.3\"}
//{\"token\":\"hlc.t\",\"amount\":1.3,\"from\":\"btcadmin\",\"to\":\"mmadmin\"}
// {\"token\":\"hlc.t\",\"amount\":10000.0,\"from\":\"btcadmin\",\"to\":\"mmadmin\"}
// btcadmin   a372f90c470ec9d5fb2ace3f16a509a9119444b92234cdae47df4923ec9c2d49
func TestSignInfo(t *testing.T) {
	sig,err := SignInfo("ttt","2e911c236ee73f4b26584cf5e48ff6f9fb9f6645ce9609f3192c6bbb6ba50566")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hex.EncodeToString(sig))
}

// {\"token\":\"hlc.t\",\"amount\":1.3,\"sender\":\"btcadmin\",\"receiver\":\"mmadmin\",\"desc\":\"desc1.3\"}
//{
//"success": true,
//"txid": "8b7720d501aeb3023c753bfc0dd02176845cd45a66666577086899c228dc28bf",
//"info": "{\"status\":200,\"message\":\"the hlc.t token from btcadmin to mmadmin had request\"}"
//}

func TestSigToPub(t *testing.T) {

	add , err := SignToAddress("{\"name\":\"fieldlee\",\"age\":34}","050b1c64e097cf836ceb3006b4c3876503ccc03428550915a6ad1436260031405883a2b33919e6bd32d97836ee1b5c7f65fd01c339bd1b55494f6e57cee7135c01")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(add)
}