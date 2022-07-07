package main
// //package jsonparser
// // we have to always use pointers because the data area always shifts 
import (
	"bytes"
//	"errors"
	"fmt"
//	d "runtime/debug"
//	"strconv"
	"os"
	"log"
	"unsafe"
	jp "jpack"
	"strconv"
)

// Replace the value of a named variable
// TODO check old value, if one is given
func ReplaceValue (data []byte, aofs int, orig string, rep string , keys ...string) (value []byte, err error) {
	dsfix := uintptr(unsafe.Pointer(&data[0]))

	sdata,_,_,_,err := jp.Get(data[aofs:], keys...)
    if err != nil {
        fmt.Printf("Error is :[%v] \n",err)
    // } else {
    // //     fmt.Printf("Sync :[%v] xy %v xx %v sf %p da %p sfxx [%v] \n",
	// // 	     string(sdata),xy, xx, &sdata[0], &data[0], string(sdata[xx-3:xx]))
	//  	fmt.Printf("type [%v] value [%v] \n",xt, string(data[xy-1:xy+9]))
	}

	// for this to work xt must be jp.String or jp.Number
	// how far into the object is the keys resolution
	usfix := uintptr(unsafe.Pointer(&sdata[0]))
	asfix := usfix - dsfix

    origb :=[]byte(orig)
    repb :=[]byte(rep)
	lxx := len(sdata)
	lxy := len(repb)
	value,err = ReplaceIxIf(data, int(asfix), int(asfix+uintptr(lxx)), &origb, &repb)

	fmt.Printf("---> value after replace  \n%v \n",string(value[:asfix+uintptr(lxy)+20]))
	return value,err
}

//Find the start of the var by moving back until we detect two '"'
func VarFind(data []byte, idx int) int {
	i := idx
	for true {
		if i >= 0 && data[i] != '"' {
			i -=1
		} else {
			break
		}
	}
	if i > 0 {
		i -=1
	}
	for true {
		if i >= 0 && data[i] != '"' {
			i -=1
		} else {
			break
		}
	}
    return i
}

//idx2 := ValueFindStart(data, idx)
// keep moving back till we find a comma or a '{'
// TODO fix for arrays
func ValueFindStart(data []byte, idx int) int {
	i := idx

	for true {
		if i >= 0 && (data[i] != ',' && data[i] != '{')  {
			i -=1
		} else {
			break
		}
	}
	return i
}
//idx2 := ValueFindStart(data, idx)
// keep moving forward till we find a comma or a '{'
// TODO fix for arrays
func ValueFindNext(data []byte, idx int) int {
	imax := len(data)
    i := idx
	for true {
		if i < imax && (data[i] != ',' && data[i] != '{')  {
			i +=1
		} else {
			break
		}
	}
	return i
}
//idx2 := ValueFindStart(data, idx)
// keep moving back till we find a comma or a '{'
// TODO fix for arrays
func ValueFindNextVar(data []byte, idx int) int {
	imax := len(data)
    i := idx
	for true {
		if i < imax && (data[i] != ':' && data[i] != '}')  {
			i +=1
		} else {
			break
		}
	}
	return i
}

func ValueFindNextArrayVar(data []byte, idx int) int {
	imax := len(data)
    i := idx
	for true {
		//fmt.Printf(" find next arrayvar i = %v  data [%v]\n", i, string(data[i]) )
		if i < imax && (data[i] != ':' && data[i] != ']')  {
			i +=1
		} else {
			break
		}
	}
	return i
}
//idx2 := ValueFindStart(data, idx)
// keep moving back till we find a comma or a '{'
// TODO fix for arrays
func ValueFindNextArray(data []byte, idx int) int {
	imax := len(data)
    i := idx
	for true {
		fmt.Printf(" find next array i = %v  data [%v]\n", i, string(data[i]) )
		if i < imax && (data[i] != '[')  {
			i +=1
		} else {
			break
		}
	}
	return i
}

//idx2 := ValueFindEnd(data, idx)
// keep moving forward till we find a comma or a '}'
// escape if we find '{'
// TODO fix for arrays
func ValueFindEnd(data []byte, idx int) int {
	i := idx
    depth := 0 
	imax := len(data)
	for true {
		if i < imax && data[i] == '{' {
			depth += 1
		}
		if i < imax && data[i] == '}' {
			if depth > 0 {
				depth -= 1
			}
		}
		if i < imax && (data[i] != ',' && (data[i] != '}') && (depth == 0))  {
			i +=1
		} else {
			break
		}
	}
	return i
}

//idx2 := ValueFindEnd(data, idx)
// keep moving forward till we find a comma or a '}'
// escape if we find '{'
// fix for arrays
func ValueFindEndArray(data []byte, idx int) int {
	i := idx
    depth := 0 
	imax := len(data)
	// we have to step past the first array detection
	for true {
		if i < imax && data[i] == '[' {
			depth += 1
		}
		if i < imax && data[i] == ']' {
			if depth > 0 {
				depth -= 1
			}
		}
		//fmt.Printf("    end array checking depth [%v] i [%v]  data [%v]\n", depth, i, string(data[i]))
		if i < imax && (data[i] != ']' && depth >= 0)  {
			i +=1
		} else {
			break
		}
	}
	return i
}
//idx2 := ValueFindEnd(data, idx)
// keep moving forward till we find a comma or a '}'
// escape if we find '{'
// fix for arrays
func ValueFindEndObject(data []byte, idx int) int {
	i := idx
    depth := 0 
	imax := len(data)
	// we have to step past the first array detection
	for true {
		if i < imax && data[i] == '{' {
			depth += 1
		}
		if i < imax && data[i] == '}' {
			if depth > 0 {
				depth -= 1
			}
		}
		if i < imax && (data[i] != '}' && depth >= 0)  {
			i +=1
		} else {
			break
		}
	}
	return i
}


//val,_ := ReplaceVarName(val, "sync_feeder" ,"my_sync_feeder","assets","feeders")
// replace a variable name (oname) with a new one (nname) 
func ReplaceVarName (data []byte, aofs int, oname string, vval string, nname string, keys ...string) (value []byte, err error) {
	dsfix := uintptr(unsafe.Pointer(&data[0]))

	sdata1,_,_,_,err := jp.Get(data[aofs:], keys...)
    if err != nil {
        fmt.Printf("Error #1 is :[%v] \n",err)
    // } else {
	//  	fmt.Printf("Variable feeders type [%v] value [%v] \n",xt, string(data[xy-1:xy+10]))
	}

	sdata2,xt,_,_,err := jp.Get(sdata1, oname)
    if err != nil {
        fmt.Printf("Error #2 is :[%v] \n",err)
    // } else {
	//  	fmt.Printf("Variable [%v] type [%v] value [%v] \n",oname, xt, string(sdata2[xy-1:xy+10]))
	}
	isNumber := xt == jp.Number
	isString := xt == jp.String
	isBool := xt == jp.Boolean

	usfix := uintptr(unsafe.Pointer(&sdata2[0]))
	asfix := usfix - dsfix

	idx := VarFind(data, int(asfix))    
	if len(vval) > 0 && (isNumber || isString || isBool) {
		fmt.Printf("ReplaceVarName check vval [%v]  data [%v] idx [%v]\n", 
												vval, string(data[idx:idx+30]), idx)

	}

	//fmt.Printf(" test var after test  [%v]\n",string(data[idx:idx+20]))
	// data[idx] is the start of the variable name so we can run 
	origb :=[]byte(oname)

	// add quotes to the name 
	xname:= strconv.Quote(nname)
    repb :=[]byte(xname)
	
	fmt.Printf("---> test before value change  [%v]\n",string(data[0:uintptr(asfix)+20]))

	value,err = ReplaceIxIf(data, idx, idx+len(oname)+2, &origb, &repb)

	fmt.Printf("---> test after value change  [%v]\n",string(value[0:uintptr(asfix)+20]))
	return value,err
}

//val,_ := RemoveItem(val, 0, "sync_feeder" ,"my_sync_feeder","assets","feeders")
// TODO add if value == vval 
func RemoveItem(data []byte, aofs int, vname string, vval string, keys ...string) (value []byte, err error) {
	dsfix := uintptr(unsafe.Pointer(&data[0]))

	sdata1,_,_,_,err := jp.Get(data[aofs:], keys...)
    if err != nil {
        fmt.Printf("Error #1 is :[%v] \n",err)
	}

	sdata2,xt,_,_,err := jp.Get(sdata1, vname)
    if err != nil {
        fmt.Printf("Error #2 is :[%v] \n",err)
	}
	isArray := xt == jp.Array
	isObject := xt == jp.Object
	isNumber := xt == jp.Number
	isString := xt == jp.String
	isBool := xt == jp.Boolean

	value = sdata2
	usfix := uintptr(unsafe.Pointer(&sdata2[0]))
	asfix := usfix - dsfix

	// //reduce asfix until we pass two '"'  to get to the start of the var 
	// //lxy := 20
	idx := VarFind(data, int(asfix))
	fmt.Printf("RemoveItem Variable #2  [%v]\n",string(data[idx:idx+20]))

	// keep moving back till we find a comma or a '{'
	idx2 := ValueFindStart(data, idx)
	fmt.Printf("RemoveItem Variable #3  [%v] idx2 [%v] \n",string(data[idx2:idx+20]), idx2)
	// keep moving forward till we find a comma or a '}'
	idx3 := ValueFindEnd(data, int(idx2+1))
	fmt.Printf("RemoveItem Variable #4  [%v] idx2 [%v] idx3 [%v] \n", string(data[idx2:idx3+1]), idx2, idx3)
	if len(vval) > 0 && (isNumber || isString || isBool) {
		fmt.Printf("RemoveItem Variable #5 check vval [%v]  data [%v] idx2 [%v] idx3 [%v] \n", 
												vval, string(data[idx2:idx3+1]), idx2, idx3)

	}
	if isArray {
		idx3 = ValueFindEndArray(data, int(idx2+1))
		idx3 += 1
		fmt.Printf("RemoveItem array    idx2 [%v] idx3 [%v] data [%v]  \n", idx2 , idx3, string(data[idx3:idx3+20]))

	}
	if isObject {
		idx3 = ValueFindEndObject(data, int(idx2+1))
		idx3 += 1
		fmt.Printf("RemoveItem object    idx2 [%v] idx3 [%v] data [%v]  \n", idx2 , idx3, string(data[idx3:idx3+20]))
	}

	// if data[idx3] == ',' then remove idx2 to idx3
	// else remove idx2-1 to idx3
	var newdata []byte
	//var newdata []byte
	//newdata = append(data[:idx3], append (nvb, data[idx3:]...)...)
    if data[idx3] == ',' {
		fmt.Printf("               found trailing comma \n")
		if data[idx2] == ',' {
			fmt.Printf("             found leading comma \n")
			newdata = append(data[:idx2], data[idx3:]...)
			os.WriteFile("output/nd1test3.json", newdata, 0666)
			return newdata,err
	// newdata = data[:idx2]
	// 		newdata = append(newdata, data[idx3:]...)
		} else {
			fmt.Printf("            no leading comma \n")
			newdata = append(data[:idx2+1], data[idx3+1:]...)
			os.WriteFile("output/nd2test3.json", newdata, 0666)
			return newdata,err
			// newdata = data[:idx2+1]
			// newdata = append(newdata, data[idx3+1:]...)
		}
		//fmt.Printf(" newdata ==>[%v] \n", string(newdata[0:idx3+20]))

	} else {
		// TODO
		fmt.Printf("               no trailing comma \n")
		if data[idx2] == ',' {
			fmt.Printf("                  found leading comma \n")
			// newdata = data[:idx2-1]
			// newdata = append(newdata, data[idx3:]...)
			newdata = append(data[:idx2-1], data[idx3:]...)
			os.WriteFile("output/nd3test3.json", newdata, 0666)
		} else {
			fmt.Printf("                   no leading comma \n")
			//newdata = []byte
			fmt.Printf("---> #0 idx2 [%v] newdata ==>[%v] \n", idx2, newdata)
			// newdata = data[:idx2-1]
			// newdata = append(newdata, data[idx3+1:]...)
			newdata = append(data[:idx2-1], data[idx3+1:]...)
			os.WriteFile("output/nd4test3.json", newdata, 0666)

			//newdata =  append(newdata,data[:idx2]...)
			//fmt.Printf("---> #1 idx2 [%v] idx3 [%v] newdata ==>[%v] \n", idx2, idx3, string(newdata[0:idx2]))
			//newdata = append(newdata, data[idx3+1:]...)
		}
	}
	fmt.Printf("--->  newdata ==>[%v] \n", string(newdata[0:idx3+20]))
    value = newdata
	return newdata,err
}

//val,_ := AddItem(val, "new_poi_feeder" ,"1234" ,"assets","feeders")
func AddItem (data []byte, aofs int, vname string, vvalue string, keys ...string) (value []byte, err error) {
	dsfix := uintptr(unsafe.Pointer(&data[0]))

	sdata1,xt,_,_,err := jp.Get(data[aofs:], keys...)
    if err != nil {
        fmt.Printf("Error #1 is :[%v] \n",err)
	}
	isArray := xt == jp.Array
	isObject := xt == jp.Object
	fmt.Printf(" --->Item type = [%v] isArray [%v] isObject [%v]\n", xt, isArray, isObject)
    
	value = sdata1
	usfix := uintptr(unsafe.Pointer(&sdata1[0]))
	asfix := usfix - dsfix
	idx2 := int(asfix)

	//move past '{' 
	idx3 := ValueFindNext(data, idx2)

	idx3 += 1
	if data[idx3] == '\n' {
		idx3 += 1
	}
 
	// create new variable
    newvar := strconv.Quote(vname)
	newvar += ":"
	newvar += vvalue
	if data[idx3] == '}' {
        fmt.Printf("AddItem  added newline to an empty object :[%v] \n",vname)
		newvar = "\n"+newvar
	}
	if data[idx3] == ']' {
        fmt.Printf("AddItem  added newline to an empty array :[%v] \n",vname)
		newvar = "\n"+newvar
	}

	// we are adding to an object 
	// obj:{ item: val, item:val} ]
	if isObject {
	}
	// we are adding to an array 
	// arr: [ {...},{...} ]
	if isArray {
		avar := "{"
		avar += newvar
		avar += "}"
		newvar = avar
		// fix idx3 to be after the array opener
		idx3 = ValueFindNextArray(data, idx2)
		idx3 +=1 
		// then find next  ':' , if we find a ']' before the var then dont add a ','
		idx4 := ValueFindNextArrayVar(data, idx3)
		if data[idx4] == ':' {
			newvar = "\n" + newvar + ",\n"
		} else if data[idx4] == ']' {
			newvar = "\n" + newvar + "\n"
		} else {
			newvar += "\n"
		}

	} else {
		// find next  ':' , if we find a '}' before the var then dont add a ','
		idx4 := ValueFindNextVar(data, idx3)
		if data[idx4] == ':' {
			newvar += ",\n"
		} else {
			newvar += "\n"
		}
	}
	nvb := []byte(newvar)

	var newdata []byte
	newdata = append(data[:idx3], append (nvb, data[idx3:]...)...)
	
	fmt.Printf("--->  newdata ==>%v<== \n", string(newdata[0:idx3+40]))

	return newdata,err
}

// TODO replace with Pos "+1"  == After, "-1" == Before
//fmt.Printf("TEST::: add item new_value41 as a variable   after new_value4  in  new_array with value 233\n")   
//val,_  = AddItemAfter(val, "new_value4", "new_value41","233", "assets","feeders","new_array")
func AddItemAfter(data []byte, aofs int, vpos string, vname string, vvalue string, keys ...string) (value []byte, err error) {
	dsfix := uintptr(unsafe.Pointer(&data[0]))

	// this finds the array
	sdata1,xt,_,_,err := jp.Get(data[aofs:], keys...)
    if err != nil {
        fmt.Printf("Error #1 is :[%v] \n",err)
	}
	isArray := xt == jp.Array

	fmt.Printf(" --->Item type = [%v] isArray [%v]\n", xt, isArray)
	// if xt is an array we need to turn the item into an object and find the '[' to insert it
    
	value = sdata1
	usfix := uintptr(unsafe.Pointer(&sdata1[0]))
	asfix := usfix - dsfix
	idx2 := int(asfix)
	//move past '{' 
	// keep moving forward till we find a comma or a '{'

	idx3 := ValueFindNext(data, idx2)
	idx3 += 1
	if data[idx3] == '\n' {
		idx3 += 1
	}
 	// create new variable
    newvar := strconv.Quote(vname)
	newvar += ":"
	newvar += vvalue
 
	if isArray {
		avar := "{"
		avar += newvar
		avar += "}"
		newvar = avar
		// fix idx3 to be after the array opener
		idx3 = ValueFindNextArray(data, idx2)
		idx3 +=1 
		// now move past vpos
		sdata2,_,_,_,err := jp.Get(sdata1, vpos)
		if err != nil {
			fmt.Printf("vpos Error #1 is :[%v] \n",err)
		// } else {
		// 	fmt.Printf(" looks like we found vpos [%s] \n", vpos)
		}
		vsfix := uintptr(unsafe.Pointer(&sdata2[0]))
		asfixa := vsfix - dsfix
		idx3a := int(asfixa)
		//fmt.Printf(" looks like we found vpos [%s] idx3a [%v] \n", vpos, idx3a)
		// if we find the end of the array next before the variable 
		idx3 = ValueFindNext(data, idx3a)
		idx3 += 1
		if data[idx3] == '\n' {
			idx3 += 1
		}
		// move past ':'
		//idx3a += 1

		// then find next  ':' , if we find a ']' before the var then dont add a ','
		idx4 := ValueFindNextArrayVar(data, idx3a)
		//fmt.Printf(" looks like we found vpos [%s] idx3a [%v] idx4 [%v] data [%v]  idx4 [%v]\n", 
		//	vpos, idx3a, idx4 , string(data[idx3a:idx4+1]), string(data[idx4]))
		// idx4 is good it is the end of either the object ':' or the array']'
		if data[idx4] == ':' {
			newvar += ",\n"
		} else {
			newvar += "\n"
		}
		if data[idx4] == ']' {
			newvar = ",\n" + newvar
			idx3 = idx4
		}

	} else {
		// find next  ':' , if we find a '}' before the var then dont add a ','
		idx4 := ValueFindNextVar(data, idx3)
		if data[idx4] == ':' {
			newvar += ",\n"
		} else {
			newvar += "\n"
		}
	}
	nvb := []byte(newvar)

	var newdata []byte
	newdata = append(data[:idx3], append (nvb, data[idx3:]...)...)
	
	fmt.Printf("--->  newdata ==>%v<== \n", string(newdata[0:idx3+40]))

	return newdata,err
}

//fmt.Printf("TEST::: add item new_value40 as a variable   before new_value4  in  new_array with value 233\n")   
//val,_  = AddItemBefore(val, "new_value4", "new_value40","-33", "assets","feeders","new_array")
func AddItemBefore (data []byte, aofs int, vpos string, vname string, vvalue string, keys ...string) (value []byte, err error) {
	dsfix := uintptr(unsafe.Pointer(&data[0]))

	// this finds the array
	sdata1,xt,_,_,err := jp.Get(data[aofs:], keys...)
    if err != nil {
        fmt.Printf("Error #1 is :[%v] \n",err)
	}
	isArray := xt == jp.Array

	fmt.Printf(" --->Item type = [%v] isArray [%v]\n", xt, isArray)
	// if xt is an array we need to turn the item into an object and find the '[' to insert it
    
	value = sdata1
	usfix := uintptr(unsafe.Pointer(&sdata1[0]))
	asfix := usfix - dsfix
	idx2 := int(asfix)
	//move past '{' 
	// keep moving forward till we find a comma or a '{'

	idx3 := ValueFindNext(data, idx2)
	idx3 += 1
	if data[idx3] == '\n' {
		idx3 += 1
	}
 	// create new variable
    newvar := strconv.Quote(vname)
	newvar += ":"
	newvar += vvalue
 
	if isArray {
		avar := "{"
		avar += newvar
		avar += "}"
		newvar = avar
		// fix idx3 to be after the array opener
		idx3 = ValueFindNextArray(data, idx2)
		idx3 +=1 
		// now move past vpos
		sdata2,_,_,_,err := jp.Get(sdata1, vpos)
		if err != nil {
			fmt.Printf("vpos Error #1 is :[%v] \n",err)
		// } else {
		// 	fmt.Printf(" looks like we found vpos [%s] \n", vpos)
		}
		vsfix := uintptr(unsafe.Pointer(&sdata2[0]))
		asfixa := vsfix - dsfix
		idx3a := int(asfixa)
		//fmt.Printf(" looks like we found vpos [%s] idx3a [%v] \n", vpos, idx3a)
		// if we find the end of the array next before the variable 

		//idx3 = ValueFindNext(data, idx3a)
		// find the start of this value 
		idx3 = VarFind(data, idx3a)
		// move past '{'
		idx3 -= 1

		// then find next  ':' , if we find a ']' before the var then dont add a ','
		// idx4 := ValueFindNextArrayVar(data, idx3a)
		// //fmt.Printf(" looks like we found vpos [%s] idx3a [%v] idx4 [%v] data [%v]  idx4 [%v]\n", 
		// //	vpos, idx3a, idx4 , string(data[idx3a:idx4+1]), string(data[idx4]))
		// // idx4 is good it is the end of either the object ':' or the array']'
		// if data[idx4] == ':' {
		newvar += ",\n"
		// } else {
		// 	newvar += "\n"
		// }
		// if data[idx4] == ']' {
		// 	newvar = ",\n" + newvar
		// 	idx3 = idx4
		// }

	} else {
		// find next  ':' , if we find a '}' before the var then dont add a ','
		idx4 := ValueFindNextVar(data, idx3)
		if data[idx4] == ':' {
			newvar += ",\n"
		} else {
			newvar += "\n"
		}
	}
	nvb := []byte(newvar)

	var newdata []byte
	newdata = append(data[:idx3], append (nvb, data[idx3:]...)...)
	
	fmt.Printf("--->  newdata ==>%v<== \n", string(newdata[0:idx3+40]))

	return newdata,err
}

// TODO run something on all matches in an array
// this returns either a name/value item in an array or the  first item in an array  
//val,_ := ArrayIdx(val, 0, "id", "asset_instances" ,"assets","feeders", "asset_instances")
func ArrayIdx(data []byte, aofs int, defname string, defval string, keys ...string) (idx int, err error) {
	dsfix := uintptr(unsafe.Pointer(&data[0]))
	adata,xt,_,xy,err := jp.Get(data[aofs:], keys...)
	isArray := xt == jp.Array

	vsfix := uintptr(unsafe.Pointer(&adata[0]))
	asfix := vsfix - dsfix
	idx = int(asfix)
	fmt.Printf(" --->Item type [%v], start [%v], size [%v], isArray [%v]\n", xt, idx, xy,isArray)
    ridx := -1
	offset := 1
	for true {
		v, t, _, o, e := jp.Get(adata[offset:])
	
		if t != jp.NotExist {
			isArray := t == jp.Array

			vsfix := uintptr(unsafe.Pointer(&v[0]))
			asfix := vsfix - dsfix
			idx = int(asfix)
			//ridx = idx
			fmt.Printf(" ------>Item type [%v], idx [%v], size [%v], isArray [%v]\n", t, idx, o, isArray)
			//If this is the array item we are looking for 
			//cb(v, t, o, e)
			if (len(defname) > 0) && (len(defval) > 0) {
				v1,t1,_,_,e1 := jp.Get(v,defname)
				if e1 == nil {
					matched := false
					if string(v1) == defval {
						matched = true
						ridx = idx
					} 
					fmt.Printf(" ------>Id type [%v], Id [%v] match [%v] \n", t1, string(v1), matched)
				}
			}
		}

		if e != nil {
			break
		}

		offset += o
	}
	return ridx,err
}

func ReplaceIx (data []byte, ix int , iy int, orig *string, rep []byte ) (value []byte, err error) {
    return append(data[:ix], append(rep, data[iy:]...)...), nil
}

func ReplaceIxIf (data []byte, ix int , iy int, orig *[]byte, rep *[]byte ) (value []byte, err error) {
	if bytes.Compare(data[ix:iy],*orig) == 0 {
		fmt.Printf("found it  %v \n",string(*orig))
	}
	return append(data[:ix], append(*rep, data[iy:]...)...), nil
}

func TestJunk (data [] byte) {
	feeders,_,_,xx,err := jp.Get(data,"assets","feeders")
    if err != nil {
        fmt.Printf("Error is :[%v] \n",err) 
	}

	sfeeders,xt,xy,xx,err := jp.Get(data,"assets","feeders","sync_feeder")
    if err != nil {
        fmt.Printf("Error is :[%v] \n",err)
    } else {
        fmt.Printf("Sync :[%v] xy %v xx %v sf %p da %p sfxx [%v] \n",string(sfeeders),xy, xx, &sfeeders[0], &data[0], string(feeders[xx-3:xx]))
		fmt.Printf("st [%v] sfeed [%v] \n",xt, string(data[xy-1:xy+9]))
	}
    lxx := len(sfeeders)
	sync,_,err := jp.GetNumber(feeders, "sync_feeder")
    
	fmt.Printf(" sync_feeder get number [%v]  len [%v]\n", sync, lxx)

    if err != nil {
        fmt.Printf("Error is :[%v] \n",err)
    }
	fmt.Printf("feeders xx is %v \n",xx)
	//fmt.Println(feeders[0])
	//fmt.Println(string(feeders[0]))


	usfix := uintptr(unsafe.Pointer(&sfeeders[0]))
	dsfix := uintptr(unsafe.Pointer(&data[0]))
	asfix := usfix - dsfix
	fmt.Printf("asfix %v \n",asfix)
	// so data[:asfix] is the old data up to the change point  and data[asfix+uintptr(lxx):] is the rest 
	// drop hte replacement in between the two
	fmt.Printf("dataxx  %v \n",string(data[asfix:asfix+uintptr(lxx)]))

	fmt.Printf("replace \"my\" with \"your\" \n")
	dataix := []byte(" this is my data")
	fmt.Printf("before replace  [%v] \n",string(dataix))
	rep := []byte(string("your"))
	orig := []byte(string("my"))
	dataix,_ = ReplaceIx (dataix, 9 , 11, nil, rep )
	fmt.Printf("after replace   [%v] \n",string(dataix))
	dataix = []byte(" this is my data")
	dataix,_ = ReplaceIxIf (dataix, 9 , 11, &orig, &rep )
	fmt.Printf("after replace   [%v] \n",string(dataix))
}

func main() {
    // file, err := os.Open("./assets.json")
	data, err := os.ReadFile("config/assets.json")
    if err != nil {
        log.Fatal(err)
    }
 
    // cmd = ReplaceValue "22" with "423" in [assets.feeders.sync_feeder]
    fmt.Printf("TEST::: replace the Value, 22  of sync_feeder with 423, need to confirm the old value\n")   
	val,_ := ReplaceValue(data, 0, "22", "423","assets","feeders","sync_feeder")
	fmt.Printf("TEST::: val after ReplaceValue  [%v] \n",string(val))   
	//permissions := 0644 // or whatever you need
	//err := 
	os.WriteFile("output/test1.json", val, 0666)
	// if err != nil { 
	// 	// handle error
	// }

	// replace sync_feeder: with my_sync_feeder 
    // cmd = ReplaceVarName "sync_feeder" value 22  with "new_sync_feeder" in [assets.feeders]
	//val,_ := 
    fmt.Printf("TEST::: replace the variable name  sync_feeder value 22 with new_sync_feeder \n")   
	val,_ = ReplaceVarName(val, 0, "sync_feeder","22", "new_sync_feeder","assets","feeders")
	os.WriteFile("output/test2.json", val, 0666)

    // cmd = RemoveItem "sync_feeder"  [with value 423] in [assets.feeders]
	//val,_ := 
    fmt.Printf("TEST::: remove variable  new_sync_feeder  value 423 TODO check value \n")   
	os.WriteFile("output/pretest3.json", val, 0666)
	val,_= RemoveItem(val, 0, "new_sync_feeder","423", "assets","feeders")
	os.WriteFile("output/test3.json", val, 0666)

    // cmd = RemoveItem "poi_feeder"  in [assets.feeders]
    fmt.Printf("TEST::: remove start variable  poi_feeder  with any value \n")   
	val,_ = RemoveItem(val, 0, "poi_feeder","", "assets","feeders")
	os.WriteFile("output/test4.json", val, 0666)


	// cmd = AddItem "new_poi_feeder" value 1234  in [assets.feeders]
	fmt.Printf("TEST::: add item as a variable   new_poi_feeder  with value 1234 TODO check for object or array\n")   
	val,_ = AddItem(val, 0, "new_poi_feeder","1234", "assets","feeders")
	os.WriteFile("output/test5.json", val, 0666)

	fmt.Printf("TEST::: add item as an object   new_object\n")   
	val,_ = AddItem(val, 0, "new_object","{}", "assets","feeders")
	os.WriteFile("output/test6.json", val, 0666)

	fmt.Printf("TEST::: add item as a variable   new_value1  with value 1234\n")   
	val,_  = AddItem(val, 0, "new_value1","1234", "assets","feeders","new_object")

	fmt.Printf("TEST::: add item as a variable   new_value2  with value 12344\n")   
	val,_  = AddItem(val, 0, "new_value2","12345", "assets","feeders","new_object")

	fmt.Printf("TEST::: add item as an array   new_array\n")   
	val,_ = AddItem(val, 0, "new_array","[]", "assets","feeders")

	fmt.Printf("TEST::: add item as a variable   new_value3  in an array with value 12344\n")   
	val,_  = AddItem(val, 0, "new_value3","123456", "assets","feeders","new_array")

	fmt.Printf("TEST::: add item as a variable   new_value4  in an array with value 344\n")   
	val,_  = AddItem(val, 0, "new_value4","344", "assets","feeders","new_array")

	fmt.Printf("TEST::: add item AFTER new_value41 as a variable   after new_value4  in  new_array with value 233\n")   
	val,_  = AddItemAfter(val, 0, "new_value4", "new_value41","233", "assets","feeders","new_array")

	fmt.Printf("TEST::: add item AFTER new_value31 as a variable   after new_value3  in  new_array with value 133\n")   
	val,_  = AddItemAfter(val, 0, "new_value3", "new_value31","133", "assets","feeders","new_array")

	fmt.Printf("TEST::: add item BEFORE new_value40 as a variable   before new_value4  in  new_array with value 4444\n")   
	val,_  = AddItemBefore(val, 0, "new_value4", "new_value40","4444", "assets","feeders","new_array")
	//fmt.Printf("TEST::: remove item new_array \n")   
	//val,_  = RemoveItem(val, "new_array", "", "assets","feeders")

	fmt.Printf("TEST::: remove item new_object \n")   
	val,_  = RemoveItem(val, 0, "new_object", "", "assets","feeders")
	os.WriteFile("output/test7.json", val, 0666)

		// we want to capture the array assets_instances with an id of feed_1 
	// we can then use that as a new base
	fmt.Printf("TEST::: create arrayidx  \n")   
	//idx,_  = 
	ArrayIdx(val, 0, "id", "feed_1", "assets","feeders","asset_instances")

	//Todo 
	// TODO find a name value pair in a range
	// TODO define a command block 
	// TODO If name/value found then do command block

	// added aofs to allow it all to happen in the context of an array
	// 
	// InsideArray 
	// do all the above in an array where key == value or where indx is defined 
    // create an arrayidx
	// fmt.Printf("TEST::: create arrayidx  \n")   
	// idx,_  = ArrayIdx(val, 0, "name", "val", "assets","feeders","asset_instances")

	// Stretch  goal Auto create update from before and after objects.

	// so how do we add this to assets.json  
	//                     "absolute_power_direction_flag":
    //                      {
    //         "name": "Charging when Enabled",
    //         "ui_type": "control",
    //         "type": "enum_slider",
    //         "var_type": "Bool",
    //         "value": false,
    //         "remote_enabled": true,
    //         "options":
    //         [
    //             {
    //                 "name": "Charge",
    //                 "value": true
    //             },
    //             {
    //                 "name": "Discharge",
    //                 "value": false
    //             }
    //         ]
    //     },
	//here
	//"ess":
	//	{
	//	"asset_instances":
	//		[
	//			{
	//				"id": "ess_#",
	//	            "components":
	//	            [
	//		          {
	//			        "component_id": "flexgen_ess_#_hs",
	//			        "variables":
	//			        {
		
	fmt.Printf("TEST::: create arrayidx  \n")   
// navigate to ess asset_instances 
	idx1,_  := ArrayIdx(val, 0, "id", "ess_#", "assets","ess","asset_instances")
	fmt.Printf("TEST::: create arrayidx id = ess_#  found [%v] \n",idx1)   
// navigate to components flexgen_ess_#_hs
	idx2,_  := ArrayIdx(val, idx1, "component_id", "flexgen_ess_#_hs", "components")
	fmt.Printf("TEST::: create arrayidx  component_id flexgen_ess_#_hs  found [%v] \n",idx2)   
    // add object to variables
	val,_ = AddItem(val, idx2, "absolute_power_direction_flag","{}", "variables")
	fmt.Printf("TEST::: val after first AddItem [%v] \n",string(val))   
    
    // add name 
	val,_ = AddItem(val, idx2, "name","\"Charging when Enabled\"", "variables","absolute_power_direction_flag")
    // ui_type 
	val,_ = AddItem(val, idx2, "ui_type","\"Control\"", "variables","absolute_power_direction_flag")
    // ui_type 
	val,_ = AddItem(val, idx2, "remote_enabled","true", "variables","absolute_power_direction_flag")
    // options 
	val,_ = AddItem(val, idx2, "options","[]", "variables","absolute_power_direction_flag")
    // name charge 
	val,_ = AddItem(val, idx2, "name","\"Charge\"", "variables","absolute_power_direction_flag","options")
	idx3,_  := ArrayIdx(val, idx2, "name", "Charge", "variables","absolute_power_direction_flag","options")
	fmt.Printf("TEST::: create arrayidx  name charge  found [%v] \n",idx3)   
    // value true 
	val,_ = AddItem(val, idx3, "value","true")
    // name Discharge 
	val,_ = AddItem(val, idx2, "name","\"Discharge\"", "variables","absolute_power_direction_flag","options")

	idx4,_  := ArrayIdx(val, idx2, "name", "Discharge", "variables","absolute_power_direction_flag","options")
	fmt.Printf("TEST::: create arrayidx  name discharge  found [%v] \n",idx4)   
	fmt.Printf("TEST::: create arrayidx  name discharge  found [%v] \n",string(val[idx4-20:idx4+20]))   
    // value true 
	val,_ = AddItem(val, idx4, "value","true")
	os.WriteFile("output/test8.json", val, 0666)

	fmt.Printf("TEST::: val at the end [%v] \n",string(val))   

// TODO wrap up Array processing with 
// AddArrayItem (val, idx2, "name","\"Charge\"", pos, "variables","absolute_power_flag","options"
// where pos is the position in the array, 0 for start, -1 for end , any other will be an index.

// TODO , perhaps fix indent 
// TODO  allow oldvalue chanck ( from /to)
// TODO complete if  if not value.  // implied if not if false
//      AddItemIf(val, idx4, "value","true", "if","val", true, stuff...) 


}
