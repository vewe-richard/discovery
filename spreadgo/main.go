package main

/*
//Spread API Reference http://www.spread.org/docs/guide/users_guide.pdf
#cgo CFLAGS: -I../spread-src-4.4.0/include/
#cgo LDFLAGS: -static -L../spread-src-4.4.0/install/lib -lm -lspread-core -ldl

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

int cgo_SP_multicast(char *msg, int msglen) {
	int ret;

	ret= SP_multicast( Mbox, SAFE_MESS, "0", 1, msglen, msg);
	return ret;
}



*/
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/mdns"
	"log"
	"os"
	"time"
)

/* How to test spread?
spread-src-4.4.0/install$ ./sbin/spread
spread-src-4.4.0/install$ ./bin/spuser   #inpput>j 0  //join group 0
spread-src-4.4.0/install$ ./bin/spuser -u root  #input>s 0 //send message to group 0
spreadgo$ go build main.go
spreadgo$ LD_LIBRARY_PATH=/home/richard/work/2022/edgecompute/tmp/spread/spread-src-4.4.0/install/lib/ ./main
*/
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify a command\n", usage())
	}
	if os.Args[1] == "testSpread" {
		testSpread()
	} else if os.Args[1] == "mdnsTestService" {
		mdnsTestService()
	} else if os.Args[1] == "mdnsTestLookup" {
		mdnsTestLookup()
	} else if os.Args[1] == "run" {
		fmt.Println("run")
		run()
	} else {
		fmt.Println("Please specify a valid command")
		fmt.Println(usage())
	}
}

func usage() string {
	return "mdnsDiscovery [testSpread | mdnsTestService | mdnsTestLookup | run]"
}

func run() {
	ret := C.cgo_SP_connect_timeout()
	if ret != C.ACCEPT_SESSION {
		fmt.Println("Connect failed ", ret)
		return
	}
	fmt.Println("Connected to spread daemon")

	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for {
			entry := <-entriesCh
			sendToSpreadBus(entry)
		}
	}()

	// Start the lookup
	fmt.Println("Lookup")
	mdns.Lookup("_foobar._tcp", entriesCh)
	time.Sleep(1000000000)
	close(entriesCh)
	C.cgo_SP_disconnect()
}

func sendToSpreadBus(entry *mdns.ServiceEntry) {
	b, _ := json.Marshal(entry)
	msg := string(b)
	fmt.Printf("Got new entry: %v\n", msg)
	C.cgo_SP_multicast(C.CString(msg), C.int(len(msg)))
}

func mdnsTestLookup() {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry)
		}
	}()

	// Start the lookup
	mdns.Lookup("_foobar._tcp", entriesCh)
	close(entriesCh)
}

func mdnsTestService() {
	// Setup our service export
	host, _ := os.Hostname()
	info := []string{"My awesome service"}
	service, _ := mdns.NewMDNSService(host, "_foobar._tcp", "", "", 8000, nil, info)

	// Create the mDNS server, defer shutdown
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()
	for true {
		fmt.Println("sleep")
		time.Sleep(1000000000)
	}
}

func testSpread() {
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

	msg := "hello world"
	C.cgo_SP_multicast(C.CString(msg), C.int(len(msg)))

	ret = C.cgo_SP_disconnect()
	fmt.Println("Disconnect ", ret)
}
