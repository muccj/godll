package main

import (
	"fmt"
	"runtime"
	"unsafe"
)
import "C"

const MAX_CALL_LEN = 65535

var globalBytes []byte = nil //共享内存与安卓或者ios原生底层沟通的buf

func main() {

}

var nativeCallChan chan string = make(chan string, 1024) //调用管道
var nativeBackChan chan string = make(chan string, 1024) //返回管道

//由外部调用 将一块内存共享到go
//export SetGlobalBytes
func SetGlobalBytes(in *C.char, num int32) {
	defer func() {
		if err := recover(); err != nil {
			_ = fmt.Errorf("原生调用 SetGlobalBytes 异常 %v\n", err)
		}
	}()
	globalBytes = (*[1024 * 1024]byte)(unsafe.Pointer(in))[:num:num]
	fmt.Printf("获得了原生传入的字节流地址 %p", globalBytes)
}

//开启管道异步处理
func StartChanLoop() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				caller, ok := <-nativeCallChan
				if ok {
					GoFunsEntry(caller)
				} else {
					break
				}
			}
		}()
	}
}

//这里可以去调用任意go端的代码了
func GoFunsEntry(caller string) {

}

//推送数据返回原生代码
func GoPushBackToNative(str string) {
	nativeBackChan <- str
}

//原生调用Go的主入口 将调用的数据推入处理管道
//export CallGo
func CallGo(in *C.char, num int32) {
	defer func() {
		if err := recover(); err != nil {
			_ = fmt.Errorf("原生调用 CallGo 异常 %v\n", err)
		}
	}()
	bts := (*[MAX_CALL_LEN]byte)(unsafe.Pointer(in))[:num:num]
	rbts := make([]byte, num)
	copy(rbts, bts[:num])
	nativeCallChan <- string(rbts)
}

//原生调用Go返回数据的主入口 由管道自动阻塞 有数据自动复制到共享内存中并返回
//export GetGoBack
func GetGoBack() int32 {
	defer func() {
		if err := recover(); err != nil {
			_ = fmt.Errorf("原生调用 GetGoBack 异常 %v\n", err)
		}
	}()
	buf := <-nativeBackChan
	bts := []byte(buf)
	copy(globalBytes, bts)
	return int32(len(bts))
}
