/* Simple Capture 
* 
* Opens a PF ring interface and prints raw packet data to the screen. 
* Should only be used for testing or the basis of a more complete program.
*
*/

package main

import (
	"flag"
	"fmt"
	"github.com/marv2097/gopfring"
	"log"
	"os"
	"strings"
)


// Command Line Args
var iface = flag.String("i", "eth0", "Interface to read packets from")
var snaplen = flag.Int("s", 65536, "Snap length (number of bytes max to read per packet")
var maxcount = flag.Int("c", -1, "Only grab this many packets, then exit")

func main() {
    flag.Parse()

    fmt.Printf("Starting capture on %s\n",*iface)
	var ring *pfring.Ring
	var err error
    var ci pfring.CaptureInfo

    // Create a new pfring Ring    
	if ring, err = pfring.NewRing(*iface, uint32(*snaplen), pfring.FlagPromisc); err != nil {
		log.Fatalln("pfring ring creation error:", err)
	}
    defer ring.Close()
    
    // Set BPF if needed
	if len(flag.Args()) > 0 {
		bpffilter := strings.Join(flag.Args(), " ")
		fmt.Fprintf(os.Stderr, "Using BPF filter %q\n", bpffilter)
		if err = ring.SetBPFFilter(bpffilter); err != nil {
			log.Fatalln("BPF filter error:", err)
		}
	}                   
    
    // Set some options
	if err = ring.SetSocketMode(pfring.ReadOnly); err != nil {
		log.Fatalln("pfring SetSocketMode error:", err)
	} else if err = ring.Enable(); err != nil {
		log.Fatalln("pfring Enable error:", err)
	}
    
    // Make a slice to use for data
    data := make([]byte, 1518)
    count := 0
    
    for {
        // Get the packet from the ring
        ci, err  = ring.ReadPacketDataTo(data)
        count++
        
        // Print the Capture Info and Raw data        
        fmt.Println("-----------------------------------------------------------------")
        fmt.Println(ci)
        fmt.Println(data[:ci.CaptureLength])
        
        // Check if over the packet limit
        limit := *maxcount > 0 && count >= *maxcount
        if limit {
            break
        }
    }
}

        
