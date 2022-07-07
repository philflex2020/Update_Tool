package main

import (
	"fmt"
	"strings"
	"strconv"
)

type Cmds struct {
	cmds map[string]Cmd
}

type Cmd struct {
	key, help string
	Func func (data []byte, fm map[string]string) (value []byte, err error)
}

func (c Cmds) addCmd (key,help string, f func (data []byte, fm map[string]string) (value []byte, err error) ) int {
	c.cmds[key]= Cmd{key:key, help:help, Func:f}
	return 0
}

func (c Cmds) runCmd (data []byte, fm map[string]string) (value []byte, err error)  {
	err = nil
	fun,ok := c.cmds[fm["func"]]
	if ok {
	   data,err =  fun.Func(data, fm)	
	// } else {
	// 	conn.Write([]byte( message +" not understood, try \"help\"  \n"))
	} 
	return data,err
}

func standardizeSpaces(s string) string {
    return strings.Join(strings.Fields(s), " ")
}

func RemoveDoubleWhiteSpace(str string) string {
    var b strings.Builder
    b.Grow(len(str))
    for i := range str {
        if !(str[i] == 32 && (i+1 < len(str) && str[i+1] == 32)) {
            b.WriteRune(rune(str[i]))
        }
    }
    return b.String()
}

func SplitOnWhiteSpace(str string) []string {
    var b []string
	i2 := -1
	i:=-1
    for i = range str {
		
        if (str[i] == 32) {
			b = append(b,str[i2+1:i])
			i2 = i
        }
    }
	b = append(b,str[i2+1:i+1])

    return b
}


func AddItem(data []byte, fm map[string]string) (value []byte, err error) {
	at,err := strconv.Atoi(fm["atx"])
	if err != nil {
		at = 0
	}
	fmt.Printf(" CMD  AddItem running\n fm [%v] at %d \n", fm, at)
	return data,nil
}

func main() {
	cmds := new(Cmds)
	cmds.cmds = make(map[string]Cmd)
	cmds.addCmd("AddItem", "Add an Item", AddItem)

	fo := "AddItem at:0      after:new_value4 as:new_value41 value:233  in:assets.feeders.new_array"
	data := []byte("This is some data")

	//fmt.Printf(" string 1 [%v]\n", foo2[1])
	//fmt.Printf(" string 2 [%v]\n", foo2[2])
	foo := RemoveDoubleWhiteSpace(string(fo))
	foo2 := SplitOnWhiteSpace(foo)
	fm := make(map[string]string)
    for x,i := range foo2 {
		fmt.Printf(" x [%v] i [%v] \n", x, i)
		if x > 0 {
			xx :=  strings.Split(i,":")
			fm[xx[0]] = xx[1]
		} else {
			fm["func"] = i
		}

	}
	fmt.Printf(" fm [%v] \n", fm)
	//AddItem( data, fm )
	data,err := cmds.runCmd(data, fm)
	if err != nil {
		fmt.Printf(" err [%v]\n", err)
	}


}
