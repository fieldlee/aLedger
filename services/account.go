package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"ledger/common"
	"ledger/key"
	"ledger/log"
	"ledger/model"
	"strconv"
	"time"
)

/// check 用户名是否可用
func AccountCheck(stub shim.ChaincodeStubInterface)pb.Response{
	_,args := stub.GetFunctionAndParameters()
	if len(args) != 1{
		log.Logger.Error("Parameters error ,please check Parameters")
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}

	accountName := common.Trim(args[0])

	if accountName == ""{
		log.Logger.Error("The Name is Blank")
		return common.SendError(common.Param_ERR,"The Name is Blank")
	}

	if len(accountName) < 3 && len(accountName) > 64 {
		log.Logger.Error("The Name must low 64 strings and higher 3 strings")
		return common.SendError(common.FORMAT_ERR,"The Name must low 64 strings and higher 3 strings")
	}

	accByte,err := stub.GetState(common.ACCOUNT_PRE + accountName)
	if err != nil {
		log.Logger.Error(err.Error())
		return common.SendError(common.GETSTAT_ERR,err.Error())
	}
	// 查询account 是否存在
	if accByte == nil {
		// 返回check 状态
		return common.SendScuess(fmt.Sprintf("%s check success",accountName))
	}

	return common.SendError(common.ACCOUNT_EXIST,fmt.Sprintf("%s is exist",args[0]))
}

////AccountChangePri
func AccountChangePri(stub shim.ChaincodeStubInterface)pb.Response{
	_,args := stub.GetFunctionAndParameters()
	if len(args) != 2{
		log.Logger.Error("Parameters error ,please check Parameters")
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}

	commonName := args[0]
	accountName := common.Trim(commonName)
	accountByte,err := stub.GetState(common.ACCOUNT_PRE + accountName)
	if err != nil {
		log.Logger.Error(err.Error())
		return common.SendError(common.GETSTAT_ERR,err.Error())
	}

	///// check admin
	if !common.CheckAdminBySign(args[0],args[1]) {
		return common.SendError(common.Param_ERR,"only admin can call this function")
	}


	if accountByte == nil {
		log.Logger.Error("the common name not check, please first call check api")
		return common.SendError(common.ACCOUNT_NOT_EXIST,"the account not exist")
	}else{
		// 创建address 和 privatekey
		address,prikey , err := key.CreateKey()
		if err != nil {
			log.Logger.Error(err.Error())
			return common.SendError(common.MARSH_ERR,err.Error())
		}
		account := model.Account{}
		err = json.Unmarshal(accountByte,&account)
		if err != nil {
			log.Logger.Error(err.Error())
			return common.SendError(common.MARSH_ERR,err.Error())
		}

		///// admin set address
		if accountName == common.Trim(common.ADMIN_Name) {
			return common.SendError(common.Param_ERR," admin can not change private key")
		}

		account.Address = common.Trim(address)

		newAccByte,err := json.Marshal(account)
		if err != nil {
			return common.SendError(common.MARSH_ERR,err.Error())
		}

		err = stub.PutState(common.ACCOUNT_PRE + account.DidName,newAccByte)
		if err != nil {
			log.Logger.Error(err.Error())
			return common.SendError(common.PUTSTAT_ERR,err.Error())
		}

		jStr := fmt.Sprintf("{\"name\":\"%s\",\"address\":\"%s\",\"privatekey\":\"%s\"}",commonName,address,prikey)

		return common.SendScuess(jStr)
	}

}
/// admin or 用户 confirm 用户并制为有效
func AccountConfirm(stub shim.ChaincodeStubInterface)pb.Response{

	_,args := stub.GetFunctionAndParameters()
	if len(args) != 2{
		log.Logger.Error("Parameters error ,please check Parameters")
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}

	accountStr    := args[0]
	accountParam  := model.AccountParam{}
	err := json.Unmarshal([]byte(accountStr),&accountParam)
	if err != nil {
		log.Logger.Error("Unmarshal:",err)
		return common.SendError(common.MARSH_ERR,err.Error())
	}

	accountName := common.Trim(accountParam.Name)
	accountByte,err := stub.GetState(common.ACCOUNT_PRE + accountName)
	if err != nil {
		log.Logger.Error(err.Error())
		return common.SendError(common.GETSTAT_ERR,err.Error())
	}

	if !common.CheckDigitLetter(accountParam.Code){
		return common.SendError(common.Param_ERR,"the code must include letter and digit")
	}

	///// check admin
	if !common.CheckAdminBySign(args[0],args[1]) {
		return common.SendError(common.Param_ERR,"only admin can call this function")
	}

	account := model.Account{}
	if accountByte == nil {
		account.Type = common.ACCOUNT
		account.DidName = accountName
		account.CommonName = accountParam.Name
		account.MspID = common.GetMsp(stub)
		account.Status = true
		account.Code = accountParam.Code
		newAccByte,err := json.Marshal(account)
		if err != nil {
			return common.SendError(common.MARSH_ERR,err.Error())
		}
		err = stub.PutState(common.ACCOUNT_PRE + account.DidName,newAccByte)
		if err != nil {
			log.Logger.Error(err.Error())
			return common.SendError(common.PUTSTAT_ERR,err.Error())
		}
		jStr := fmt.Sprintf("{\"name\":\"%s\",\"code\":\"%s\"}",accountParam.Name,accountParam.Code)
		return common.SendScuess(jStr)
	}else{
		log.Logger.Error("the common name not check, please first call check api")
		return common.SendError(common.ACCOUNT_NOT_EXIST,"the account had created")
	}
}

func GeneratePriKey(stub shim.ChaincodeStubInterface)pb.Response{
	_,args := stub.GetFunctionAndParameters()
	if len(args) != 1{
		log.Logger.Error("Parameters error ,please check Parameters")
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}

	accountStr := args[0]
	accountParam :=	model.AccountParam{}
	err := json.Unmarshal([]byte(accountStr),&accountParam)
	if err != nil {
		log.Logger.Error("Unmarshal:",err)
		return common.SendError(common.MARSH_ERR,err.Error())
	}

	accountName := common.Trim(accountParam.Name)
	accountByte,err := stub.GetState(common.ACCOUNT_PRE + accountName)
	if err != nil {
		log.Logger.Error(err.Error())
		return common.SendError(common.GETSTAT_ERR,err.Error())
	}
	account := model.Account{}
	if accountByte == nil {
		log.Logger.Error("the common name not check, please first call check api")
		return common.SendError(common.ACCOUNT_NOT_EXIST,"the account had not created")
	}else{
		err = json.Unmarshal(accountByte,&account)
		if err != nil {
			log.Logger.Error(err.Error())
			return common.SendError(common.MARSH_ERR,err.Error())
		}
		if account.Address != ""{
			jStr := fmt.Sprintf("{\"name\":\"%s\",\"code\":\"%s\",\"address\":\"%s\",\"status\":false,\"desc\":\"the account had generate private key\"}",account.CommonName,account.Code,account.Address)
			return common.SendScuess(jStr)
		}
		// 创建address 和 privatekey
		address,prikey , err := key.CreateKey()
		if err != nil {
			log.Logger.Error(err.Error())
			return common.SendError(common.MARSH_ERR,err.Error())
		}
		///// admin set address
		if accountName == common.Trim(common.ADMIN_Name) {
			address = common.ADMIN_ADDRESS
			prikey = ""
		}

		account.Address = common.Trim(address)

		newAccByte,err := json.Marshal(account)
		if err != nil {
			return common.SendError(common.MARSH_ERR,err.Error())
		}

		err = stub.PutState(common.ACCOUNT_PRE + account.DidName,newAccByte)
		if err != nil {
			log.Logger.Error(err.Error())
			return common.SendError(common.PUTSTAT_ERR,err.Error())
		}

		jStr := fmt.Sprintf("{\"name\":\"%s\",\"code\":\"%s\",\"address\":\"%s\",\"status\":true,\"private\":\"%s\"}",account.CommonName,account.Code,address,prikey)

		return common.SendScuess(jStr)
	}
}

// 查找用户
func AccountGetByName(stub shim.ChaincodeStubInterface,didName string)(model.Account,error){

	accountName := common.Trim(didName)

	didByte,err := stub.GetState(common.ACCOUNT_PRE + accountName)
	if err != nil {
		log.Logger.Error(err.Error())
		return model.Account{},err
	}
	if didByte == nil {
		log.Logger.Error(err.Error())
		return model.Account{},nil
	}

	didAccount := model.Account{}
	err = json.Unmarshal(didByte,&didAccount)
	if err != nil {
		return model.Account{},err
	}
	return didAccount,nil
}
/// 锁定账号
func AccountLock(stub shim.ChaincodeStubInterface)pb.Response{
	_,args := stub.GetFunctionAndParameters()

	if len(args) != 2{
		log.Logger.Error("Parameters error ,please check Parameters")
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}
	///// check admin
	if !common.CheckAdminBySign(args[0],args[1]) {
		return common.SendError(common.Right_ERR,"only admin can call this function")
	}


		accountName := common.Trim(args[0])
		byteAccount,err := stub.GetState(common.ACCOUNT_PRE + accountName)
		if err != nil {
			log.Logger.Error(err.Error())
			return common.SendError(common.GETSTAT_ERR,err.Error())
		}
		account := model.Account{}
		err = json.Unmarshal(byteAccount,&account)
		if err != nil {
			log.Logger.Error(err.Error())
			return common.SendError(common.MARSH_ERR,err.Error())
		}

		if account.Status == false {
			return common.SendScuess(fmt.Sprintf("%s had locked",accountName))
		}
		account.Type = common.ACCOUNT
		account.Status = false

		accountByte,err := json.Marshal(account)
		if err != nil{
			log.Logger.Error(err.Error())
			return common.SendError(common.MARSH_ERR,err.Error())
		}
		err = stub.PutState(common.ACCOUNT_PRE + accountName,accountByte)
		if err != nil{
			log.Logger.Error(err.Error())
			return common.SendError(common.PUTSTAT_ERR,err.Error())
		}
		return common.SendScuess(fmt.Sprintf("%s had locked",account.DidName))

}
/// 解锁账号
func AccountUNLock(stub shim.ChaincodeStubInterface)pb.Response{
	_,args := stub.GetFunctionAndParameters()

	if len(args) != 2{
		log.Logger.Error("Parameters error ,please check Parameters")
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}
	///// check admin
	if !common.CheckAdminBySign(args[0],args[1]) {
		return common.SendError(common.Right_ERR,"only admin can call this function")
	}


		accountName := common.Trim(args[0])
		byteAccount,err := stub.GetState(common.ACCOUNT_PRE + accountName)
		if err != nil {
			return common.SendError(common.GETSTAT_ERR,err.Error())
		}
		account := model.Account{}
		err = json.Unmarshal(byteAccount,&account)
		if err != nil {
			log.Logger.Error(err.Error())
			return common.SendError(common.MARSH_ERR,err.Error())
		}
		if account.Status == true {
			return common.SendScuess(fmt.Sprintf("%s had unlock",accountName))
		}
		account.Type = common.ACCOUNT
		account.Status = true
		accountByte,err := json.Marshal(account)
		if err != nil{
			log.Logger.Error(err.Error())
			return common.SendError(common.MARSH_ERR,err.Error())
		}
		err = stub.PutState(common.ACCOUNT_PRE + accountName,accountByte)
		if err != nil{
			log.Logger.Error(err.Error())
			return common.SendError(common.PUTSTAT_ERR,err.Error())
		}
		return common.SendScuess(fmt.Sprintf("%s had unlock",accountName))

}
/// get history
func AccountGetHistory(stub shim.ChaincodeStubInterface)pb.Response{
	_,args := stub.GetFunctionAndParameters()
	if len(args) != 1{
		log.Logger.Error("Parameters error ,please check Parameters")
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}
	accountName := common.Trim(args[0])
	history, err := stub.GetHistoryForKey(common.ACCOUNT_PRE + accountName)
	if err != nil {
		log.Logger.Error(err.Error())
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

////　ｇｅｔ account
func AccountGet(stub shim.ChaincodeStubInterface)pb.Response{
	_,args := stub.GetFunctionAndParameters()
	if len(args) != 1{
		log.Logger.Error("Parameters error ,please check Parameters")
		return common.SendError(common.Param_ERR,"Parameters error ,please check Parameters")
	}

	accountName := common.Trim(args[0])

	accByte,err := stub.GetState(common.ACCOUNT_PRE + accountName)
	if err != nil {
		log.Logger.Error(err.Error())
		return common.SendError(common.GETSTAT_ERR,err.Error())
	}

	return common.SendScuess(string(accByte))
}

////　get all account
func AccountGetAll(stub shim.ChaincodeStubInterface)pb.Response{
	queryString := "{\"selector\":{\"type\":\"account\"}}"
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