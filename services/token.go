package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"ledger/common"
	"ledger/model"
	"strconv"
	"time"
)

// 创建token
func TokenCreate(stub shim.ChaincodeStubInterface)pb.Response{

	_,args := stub.GetFunctionAndParameters()

	if len(args) != 2 {
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}

	///// check admin
	if !common.CheckAdminBySign(args[0],args[1]) {
		return common.SendError(common.Param_ERR,"only admin can call this function")
	}

	tokenParam := model.TokenParam{}
	err := json.Unmarshal([]byte(args[0]),&tokenParam)
	if err != nil {
		return common.SendError(common.MARSH_ERR,err.Error())
	}

	tokenname := common.Trim(tokenParam.Token)
	desc := tokenParam.Desc

	tokenByte,err := stub.GetState(common.TOKEN_PRE+tokenname)
	if err != nil {
		return common.SendError(common.GETSTAT_ERR,err.Error())
	}

	token := model.Token{}

	if tokenByte == nil{
		token.Type = common.TOKEN
		token.Status = true
		token.Action = common.TOKEN_INIT
		token.Desc = desc
		token.Name = tokenname
		token.Amount = float64(0)
		token.Issuer = common.ADMIN_Name
		tokenNewByte,err := json.Marshal(token)
		if err != nil {
			return common.SendError(common.MARSH_ERR,err.Error())
		}
		err = stub.PutState(common.TOKEN_PRE+tokenname,tokenNewByte)
		if err != nil {
			return common.SendError(common.PUTSTAT_ERR,err.Error())
		}
		return common.SendScuess(fmt.Sprintf("%s token had create",tokenname))
	}

	return common.SendError(common.TKNERR_EXIST,fmt.Sprintf("%s is exist",tokenname))
}

// 查询token

func TokenGet(stub shim.ChaincodeStubInterface,tokenname string)(model.Token, error){
	uptokename := common.Trim(tokenname)

	tokenByte,err := stub.GetState(common.TOKEN_PRE+uptokename)
	if err != nil {
		return model.Token{},err
	}
	if tokenByte == nil {
		return model.Token{},errors.New("the token is not exist")
	}
	token := model.Token{}

	err = json.Unmarshal(tokenByte,&token)
	if err != nil {
		return model.Token{},err
	}
	return token,nil
}

// 查询token

func TokenGetByName(stub shim.ChaincodeStubInterface)pb.Response{

	_,args := stub.GetFunctionAndParameters()

	if len(args) != 1 {
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}

	tokenname := common.Trim(args[0])

	tokenByte,err := stub.GetState(common.TOKEN_PRE+tokenname)
	if err != nil {
		return common.SendError(common.GETSTAT_ERR,err.Error())
	}
	if tokenByte == nil {
		return common.SendError(common.GETSTAT_ERR,err.Error())
	}
	return common.SendScuess(string(tokenByte))
}

// Token 状态

func TokenUpdateDisable(stub shim.ChaincodeStubInterface)pb.Response{

	_,args := stub.GetFunctionAndParameters()

	if len(args) != 2 {
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}

	///// check admin
	if !common.CheckAdminBySign(args[0],args[1]) {
		return common.SendError(common.Param_ERR,"only admin can call this function")
	}

	tokenname := common.Trim(args[0])

	tokenByte,err := stub.GetState(common.TOKEN_PRE+tokenname)
	if err != nil {
		return common.SendError(common.GETSTAT_ERR,err.Error())
	}
	token := model.Token{}

	err = json.Unmarshal(tokenByte,&token)

	if err != nil {
		return common.SendError(common.MARSH_ERR,err.Error())
	}

	token.Status = false

	//	 保存
	tokenByte , err = json.Marshal(token)
	if err != nil {
		return common.SendError(common.MARSH_ERR,err.Error())
	}
	err = stub.PutState(common.TOKEN_PRE+tokenname,tokenByte)
	if err != nil {
		return common.SendError(common.PUTSTAT_ERR,err.Error())
	}
	return common.SendScuess(fmt.Sprintf("%s had update disable",tokenname))
}

// Token 修改状态

func TokenUpdateEnable(stub shim.ChaincodeStubInterface)pb.Response{

	_,args := stub.GetFunctionAndParameters()

	if len(args) != 2 {
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}

	///// check admin
	if !common.CheckAdminBySign(args[0],args[1]) {
		return common.SendError(common.Param_ERR,"only admin can call this function")
	}

	tokenname := common.Trim(args[0])
	tokenByte,err := stub.GetState(common.TOKEN_PRE+tokenname)
	if err != nil {
		return common.SendError(common.GETSTAT_ERR,err.Error())
	}
	token := model.Token{}

	err = json.Unmarshal(tokenByte,&token)

	if err != nil {
		return common.SendError(common.MARSH_ERR,err.Error())
	}

	token.Status = true
	//	 保存
	tokenByte , err = json.Marshal(token)
	if err != nil {
		return common.SendError(common.MARSH_ERR,err.Error())
	}
	err = stub.PutState(common.TOKEN_PRE+tokenname,tokenByte)
	if err != nil {
		return common.SendError(common.PUTSTAT_ERR,err.Error())
	}
	return common.SendScuess(fmt.Sprintf("%s had update enable",tokenname))
}

// 查询token记录

func TokenGetHistory(stub shim.ChaincodeStubInterface)pb.Response{
	_,args := stub.GetFunctionAndParameters()

	if len(args) != 1{
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}

	tokenName := common.Trim(args[0])

	history, err := stub.GetHistoryForKey(common.TOKEN_PRE + tokenName)

	if err != nil {
		return common.SendError(common.GETSTAT_ERR,err.Error())
	}

	defer  history.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	for history.HasNext(){
		response ,err := history.Next()
		if err != nil {
			continue
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}
		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return common.SendScuess(buffer.String())
}

func TokenList(stub shim.ChaincodeStubInterface) pb.Response{
	queryString := "{\"selector\":{\"type\":\"token\"}}"
	resultsIterator, err := stub.GetQueryResult(queryString)
	defer resultsIterator.Close()
	if err != nil {
		return common.SendError(common.GETSTAT_ERR,err.Error())
	}
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse,err := resultsIterator.Next()
		if err != nil {
			return common.SendError(common.MARSH_ERR,err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return common.SendScuess(buffer.String())
}