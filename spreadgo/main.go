package main

/*
//Spread API Reference http://www.spread.org/docs/guide/users_guide.pdf
#cgo CFLAGS: -I../spread-src-4.4.0/include/
#cgo LDFLAGS: -L../spread-src-4.4.0/install/lib -lm -lspread-core -ldl
//-static
#include <stdio.h>
#include <sp.h>
typedef struct {
	int     mver, miver, pver;
}version ;

static  char    Private_group[MAX_GROUP_NAME];
static  mailbox Mbox;
static	char	User[80];
static  char    Spread_name[80];

int cgo_SP_connect_timeout(){
	int ret;

	sp_time test_timeout;
	test_timeout.sec = 5;
	test_timeout.usec = 0;

	sprintf( User, "go" );
	sprintf( Spread_name, "4803");

    ret = SP_connect_timeout( Spread_name, User, 0, 1, &Mbox, Private_group, test_timeout );
	return ret;
}

int cgo_SP_disconnect(){
	return SP_disconnect( Mbox );
}

int cgo_SP_multicast() {
	int ret;

	ret= SP_multicast( Mbox, SAFE_MESS, "0", 1, 13, "jiangjqian" );
	return ret;
}



*/
import "C"
import "fmt"

//Test
//spread-src-4.4.0/install$ ./sbin/spread
//spread-src-4.4.0/install$ ./bin/spuser   #inpput>j 0  //join group 0
//spread-src-4.4.0/install$ ./bin/spuser -u root  #input>s 0 //send message to group 0
func main() {
	v := C.version{}
	C.SP_version(&v.mver, &v.miver, &v.pver)
	fmt.Println("spread version: ", v.mver, v.miver, v.pver)

	//ret != ACCEPT_SESSION

	ret := C.cgo_SP_connect_timeout()
	if ret != C.ACCEPT_SESSION {
		fmt.Println("Connect failed ", ret)
		return
	}
	fmt.Println("Connected to spread daemon")

	C.cgo_SP_multicast()

	ret = C.cgo_SP_disconnect()
	fmt.Println("Disconnect ", ret)
}
