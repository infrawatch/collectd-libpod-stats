package host

/*
#include "globals.h"
*/
//import "C"
import "fmt"

//GlobalHostname retrieve collectd global hostname setting (hostname_g)
func GlobalHostname() {
	//hostname := C.GoString(C.hostname_g)
	hostname := "hellp"
	fmt.Printf("This is ma global host name!!! %s\n", hostname)
}
