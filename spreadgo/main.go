package main

/*
#cgo CFLAGS: -I../spread-src-4.4.0/include/
#cgo LDFLAGS: -static -L../spread-src-4.4.0/libspread -lm -lspread-core -ldl
#include <sp.h>
typedef struct {
	int     mver, miver, pver;
}version ;
*/
import "C"
import "fmt"

func main() {
	v := C.version{}
	C.SP_version(&v.mver, &v.miver, &v.pver)
	fmt.Println("spread version: ", v.mver, v.miver, v.pver)
}
