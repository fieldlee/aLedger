package common

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"ledger/key"
	"ledger/log"
	"ledger/model"
	"math/big"
	"strconv"
	"strings"
	"unicode"
)


/** 获取交易发起方的MSPID **/
func GetMspid(stub shim.ChaincodeStubInterface) (string) {
	createrbyte, err := stub.GetCreator() //获得创建者
	if err != nil {
		log.Logger.Error("shim GetCreater error", err.Error())
		return ""
	}
	//解析MSPID
	newbytes := []byte{}
	headFlg := true
	for i := 0; i < len(createrbyte); i++ {
		if createrbyte[i] >= 33 && createrbyte[i] <= 126 {
			headFlg = false
			newbytes = append(newbytes, createrbyte[i])
		}
		if createrbyte[i] < 33 || createrbyte[i] > 126 {
			if !headFlg {
				break
			}
		}
	}
	return string(newbytes)
}

func ComMD(mole uint,deno uint) float64{
	fMole := new(big.Float)
	fMole.SetUint64(uint64(mole))
	fDeno := new(big.Float)
	fDeno.SetUint64(uint64(deno))

	dValue := new(big.Float).Quo(fMole,fDeno)
	v,_ := dValue.Float64()
	return Decimal(v)
}

func ComputeForMD(value float64,mole uint,deno uint) float64{
	fMole := new(big.Float)
	fMole.SetUint64(uint64(mole))
	fDeno := new(big.Float)
	fDeno.SetUint64(uint64(deno))

	fValue := new(big.Float)
	fValue.SetFloat64(value)

	mValue := new(big.Float).Mul(fValue,fMole)
	dValue := new(big.Float).Quo(mValue,fDeno)

	v,_ := dValue.Float64()
	return Decimal(v)
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func FloatSub(first float64,second float64) float64{
	bigF := big.NewFloat(first)
	bigT := big.NewFloat(second)
	bigSUb := new(big.Float).Sub(bigF,bigT)
	value , _ := bigSUb.Float64()
	return Decimal(value)
}

func FloatAdd(first float64,second float64) float64{
	bigF := big.NewFloat(first)
	bigT := big.NewFloat(second)
	bigAdd := new(big.Float).Add(bigF,bigT)
	value ,_:= bigAdd.Float64()
	return Decimal(value)
}

func GetMsp(stub shim.ChaincodeStubInterface)(string){
	id, err := cid.New(stub)
	log.Logger.Info(id)
	if err != nil {
		log.Logger.Error("shim getMsp error", err.Error())
	}
	mspid, err := id.GetMSPID()
	if err != nil {
		log.Logger.Error("shim getMsp error", err.Error())
	}
	log.Logger.Info(mspid)
	return mspid
}

func GetRight(stub shim.ChaincodeStubInterface)(string){
	id, err := cid.New(stub)
	if err != nil {
		log.Logger.Error("shim GetRight error", err.Error())
	}

	cert, err := id.GetX509Certificate()
	if err != nil {
		log.Logger.Error("shim GetRight error", err.Error())
	}

	//id.GetAttributeValue()

	log.Logger.Info(id)
	if cert.IsCA {
		return "Admin"
	}else{
		return "Member"
	}
}

func SendError(errno int, msg string) pb.Response {
	log.Logger.Error(msg)
	returnJson := model.ReturnJson{
		Status:errno,
		Message:msg,
	}
	returnByte,_ := json.Marshal(returnJson)
	return shim.Error(string(returnByte))
}

func SendScuess(msg string) pb.Response{
	returnJson := model.ReturnJson{
		Status:200,
		Message:msg,
	}
	returnByte,_ := json.Marshal(returnJson)
	return shim.Success(returnByte)
}

func GetCommonName(stub shim.ChaincodeStubInterface)( string, error){
	cert,err := cid.New(stub)
	if err != nil {
		log.Logger.Error(err)
		return "",err
	}
	certfiaction,err := cert.GetX509Certificate()
	if err != nil {
		log.Logger.Error(err)
		return "",err
	}
	return certfiaction.Subject.CommonName,nil
}

func Trim(trimS string)(string){
	return strings.ToUpper(strings.TrimSpace(trimS))
}

func GetIsAdmin(name string)( bool, error){

	log.Logger.Error("current username :", name)

	if strings.ToUpper(strings.TrimSpace(name)) == strings.ToUpper(strings.TrimSpace(ADMIN_Name)){
		return true , nil
	}

	return false,nil
}

func IsSuperAdmin(name string)(bool){

	//orgid := GetMsp(stub)
	//
	//isAdmin, err:= GetIsAdmin(stub)
	//
	//if err != nil {
	//	log.Logger.Error("IsSuperAdmin error", err.Error())
	//}
	//comName , err := GetCommonName(stub)
	//if err != nil {
	//	log.Logger.Error("GetCommonName error", err.Error())
	//}
	//if strings.ToLower(orgid) == strings.ToLower(ADMIN_ORG) {
	//	if isAdmin == true {
	//		return  true
	//	}
	//	if strings.ToLower(comName) == strings.ToLower(ADMIN_Name) {
	//		return true
	//	}
	//}
	//return false

	if strings.ToUpper(strings.TrimSpace(name)) == strings.ToUpper(strings.TrimSpace(ADMIN_Name)){
		return true
	}

	return false
}

func CheckAdminBySign(sign , signed string) (bool) {
	address,err := key.SignToAddress(sign,signed)
	if err != nil {
		return false
	}

	if Trim(address) == Trim(ADMIN_ADDRESS) {
		return true
	}
	return false
}

func CheckDigitLetter(codestr string) bool{
	hasDigit := false
	hasLetter := false
	for _,v := range codestr {
		if unicode.IsDigit(v){
			hasDigit = true
		}
		if unicode.IsLetter(v){
			hasLetter = true
		}
	}

	if hasLetter && hasDigit {
		return true
	}
	return  false
}

func CheckUserEnable(stub shim.ChaincodeStubInterface)(bool){

	return true
}
func CheckTokenEnable(stub shim.ChaincodeStubInterface,tokenName string)(bool){
	return true
}