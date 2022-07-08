package main
// //package jsonparser
// // we have to always use pointers because the data area always shifts 
import (
//	"bytes"
//	"errors"
	"fmt"
//	d "runtime/debug"
//	"strconv"
	"os"
	"log"
//	"unsafe"
//	jp "jpack"
	ut "utool"
	uc "ucmd"
	"strconv"
	"bufio"
)


func tmain() {
	// TODO specify target file , command file
	cfgFile := "config/assets.json"
	cmdFile := "commands/assets_9_3.txt"
    // file, err := os.Open("./assets.json")
	data, err := os.ReadFile(cfgFile)
    if err != nil {
        log.Fatal(err)
    }
	//cmds
	cmds, err := os.Open(cmdFile)
    if err != nil {
         log.Print(err)
    }
	defer cmds.Close() 

	cmdlines := bufio.NewScanner(cmds)
	for cmdlines.Scan() {
		fmt.Println(cmdlines.Text())
	}

    // cmd = ReplaceValue "22" with "423" in [assets.feeders.sync_feeder]
    fmt.Printf("TEST::: replace the Value, 22  of sync_feeder with 423, need to confirm the old value\n")   
	val,_ := ut.ReplaceValue(data, 0, "22", "423","assets","feeders","sync_feeder")
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
	val,_ = ut.ReplaceVarName(val, 0, "sync_feeder","22", "new_sync_feeder","assets","feeders")
	os.WriteFile("output/test2.json", val, 0666)

    // cmd = RemoveItem "sync_feeder"  [with value 423] in [assets.feeders]
	//val,_ := 
    fmt.Printf("TEST::: remove variable  new_sync_feeder  value 423 TODO check value \n")   
	os.WriteFile("output/pretest3.json", val, 0666)
	val,_= ut.RemoveItem(val, 0, "new_sync_feeder","423", "assets","feeders")
	os.WriteFile("output/test3.json", val, 0666)

    // cmd = RemoveItem "poi_feeder"  in [assets.feeders]
    fmt.Printf("TEST::: remove start variable  poi_feeder  with any value \n")   
	val,_ = ut.RemoveItem(val, 0, "poi_feeder","", "assets","feeders")
	os.WriteFile("output/test4.json", val, 0666)


	// cmd = AddItem "new_poi_feeder" value 1234  in [assets.feeders]
	fmt.Printf("TEST::: add item as a variable   new_poi_feeder  with value 1234 TODO check for object or array\n")   
	val,_ = ut.AddItem(val, 0, "new_poi_feeder","1234", "assets","feeders")
	os.WriteFile("output/test5.json", val, 0666)

	fmt.Printf("TEST::: add item as an object   new_object\n")   
	val,_ = ut.AddItem(val, 0, "new_object","{}", "assets","feeders")
	os.WriteFile("output/test6.json", val, 0666)

	fmt.Printf("TEST::: add item as a variable   new_value1  with value 1234\n")   
	val,_  = ut.AddItem(val, 0, "new_value1","1234", "assets","feeders","new_object")

	fmt.Printf("TEST::: add item as a variable   new_value2  with value 12344\n")   
	val,_  = ut.AddItem(val, 0, "new_value2","12345", "assets","feeders","new_object")

	fmt.Printf("TEST::: add item as an array   new_array\n")   
	val,_ = ut.AddItem(val, 0, "new_array","[]", "assets","feeders")

	fmt.Printf("TEST::: add item as a variable   new_value3  in an array with value 12344\n")   
	val,_  = ut.AddItem(val, 0, "new_value3","123456", "assets","feeders","new_array")

	fmt.Printf("TEST::: add item as a variable   new_value4  in an array with value 344\n")   
	val,_  = ut.AddItem(val, 0, "new_value4","344", "assets","feeders","new_array")

	fmt.Printf("TEST::: add item AFTER new_value41 as a variable   after new_value4  in  new_array with value 233\n")   
	val,_  = ut.AddItemAfter(val, 0, "new_value4", "new_value41","233", "assets","feeders","new_array")

	fmt.Printf("TEST::: add item AFTER new_value31 as a variable   after new_value3  in  new_array with value 133\n")   
	val,_  = ut.AddItemAfter(val, 0, "new_value3", "new_value31","133", "assets","feeders","new_array")

	fmt.Printf("TEST::: add item BEFORE new_value40 as a variable   before new_value4  in  new_array with value 4444\n")   
	val,_  = ut.AddItemBefore(val, 0, "new_value4", "new_value40","4444", "assets","feeders","new_array")
	//fmt.Printf("TEST::: remove item new_array \n")   
	//val,_  = RemoveItem(val, "new_array", "", "assets","feeders")

	fmt.Printf("TEST::: remove item new_object \n")   
	val,_  = ut.RemoveItem(val, 0, "new_object", "", "assets","feeders")
	os.WriteFile("output/test7.json", val, 0666)

		// we want to capture the array assets_instances with an id of feed_1 
	// we can then use that as a new base
	fmt.Printf("TEST::: create arrayidx  \n")   
	//idx,_  = 
	ut.ArrayIdx(val, 0, "id", "feed_1", "assets","feeders","asset_instances")

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
	idx1,_  := ut.ArrayIdx(val, 0, "id", "ess_#", "assets","ess","asset_instances")
	fmt.Printf("TEST::: create arrayidx id = ess_#  found [%v] \n",idx1)   
// navigate to components flexgen_ess_#_hs
	idx2,_  := ut.ArrayIdx(val, idx1, "component_id", "flexgen_ess_#_hs", "components")
	fmt.Printf("TEST::: create arrayidx  component_id flexgen_ess_#_hs  found [%v] \n",idx2)   
    // add object to variables
	val,_ = ut.AddItem(val, idx2, "absolute_power_direction_flag","{}", "variables")
	fmt.Printf("TEST::: val after first AddItem [%v] \n",string(val))   
    
    // add name 
	val,_ = ut.AddItem(val, idx2, "name","\"Charging when Enabled\"", "variables","absolute_power_direction_flag")
    // ui_type 
	val,_ = ut.AddItem(val, idx2, "ui_type","\"Control\"", "variables","absolute_power_direction_flag")
    // ui_type 
	val,_ = ut.AddItem(val, idx2, "remote_enabled","true", "variables","absolute_power_direction_flag")
    // options 
	val,_ = ut.AddItem(val, idx2, "options","[]", "variables","absolute_power_direction_flag")
    // name charge 
	val,_ = ut.AddItem(val, idx2, "name","\"Charge\"", "variables","absolute_power_direction_flag","options")
	idx3,_  := ut.ArrayIdx(val, idx2, "name", "Charge", "variables","absolute_power_direction_flag","options")
	fmt.Printf("TEST::: create arrayidx  name charge  found [%v] \n",idx3)   
    // value true 
	val,_ = ut.AddItem(val, idx3, "value","true")
    // name Discharge 
	val,_ = ut.AddItem(val, idx2, "name","\"Discharge\"", "variables","absolute_power_direction_flag","options")

	idx4,_  := ut.ArrayIdx(val, idx2, "name", "Discharge", "variables","absolute_power_direction_flag","options")
	fmt.Printf("TEST::: create arrayidx  name discharge  found [%v] \n",idx4)   
	fmt.Printf("TEST::: create arrayidx  name discharge  found [%v] \n",string(val[idx4-20:idx4+20]))   
    // value true 
	val,_ = ut.AddItem(val, idx4, "value","true")
	os.WriteFile("output/test8.json", val, 0666)

	fmt.Printf("TEST::: val at the end [%v] \n",string(val))   
	//testCmd()
// TODO wrap up Array processing with 
// AddArrayItem (val, idx2, "name","\"Charge\"", pos, "variables","absolute_power_flag","options"
// where pos is the position in the array, 0 for start, -1 for end , any other will be an index.

// TODO , perhaps fix indent 
// TODO  allow oldvalue chanck ( from /to)
// TODO complete if  if not value.  // implied if not if false
//      AddItemIf(val, idx4, "value","true", "if","val", true, stuff...) 


}

func AddItemCmd(data []byte, fm map[string]string) (value []byte, err error) {
	at,err := strconv.Atoi(fm["at"])
	if err != nil {
		at = 0
		fmt.Printf("error [%v]\n", err)
	}
	// note that the command returns idx and res 
	data,_ = ut.AddItem(data, 0, fm["name"],fm["value"])
	fm["res"]= "OK"
	fm["idx"]= strconv.Itoa(at)
	fmt.Printf(" CMD  AddItem running\n fm [%v] at %d \n", fm, at)
	return data,nil
}

func main () {
	cfgFile := "config/basic.json"
	cmdFile := "commands/basic.txt"
    // file, err := os.Open("./assets.json")
	data, err := os.ReadFile(cfgFile)
    if err != nil {
        log.Fatal(err)
    }
	//cmds
	cmds, err := os.Open(cmdFile)
    if err != nil {
         log.Print(err)
    }
	defer cmds.Close() 

	cmdlines := bufio.NewScanner(cmds)
	for cmdlines.Scan() {
    	// Task #3 read all the cmdlines into command structure
		// Task #4 read all the variables into the variable structure  
		fmt.Println(cmdlines.Text())
	}
    // this is the raw command
	data,_ = ut.AddItem(data, 0, "test","\"A Test Value\"")

    // this is the command from the command file

	// task #1 modify ucmd.go to allow strings with spaces in the value
	// task #2 modify Get to use "assets|feeders|variables"
    
	// task #5 put a newline in front of the first object
	//TEST::: data at the end [
	//	{
	//	"test2":A_Test_Value,
    //  "test":"A Test Value"
    //  }
	// ]


	// cmd := "AddItem at:0 name:test value:\\\"A Test Value\\\""
	cmd := "AddItem at:0 name:test2 value:A_Test_Value"
	
	// we need to create a function map fm in ucmd but values cannot have spaces yet .. fix it
	//uc.Tmain()
	fm,_ := uc.Make_fm(cmd)

	// now we need the AddItemCmd to run the ut.AddItem funtion
	data,_ = AddItemCmd(data, fm)

	
	os.WriteFile("output/test1.json", data, 0666)
	fmt.Printf("TEST::: command was [%s] \n",cmd)   
	fmt.Printf("TEST::: fm is [%v] \n",fm)   

	fmt.Printf("TEST::: data at the end [%v] \n",string(data))   

}