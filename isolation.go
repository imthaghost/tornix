package tornix

import (
	"fmt"
	"io/ioutil"
	"math/rand"

	"github.com/fatih/color"
	"github.com/imroc/req"
)

// StreamIsolation Separate streams across circuits by connection metadata
//		 When a stream arrives at Tor, we have the following data to examine:
//		 1) The destination address
//		 2) The destination port (unless this a DNS lookup)
//		 3) The protocol used by the application to send the stream to Tor:
//			SOCKS4, SOCKS4A, SOCKS5, or whatever local "transparent proxy"
//			mechanism the kernel gives us.
//		 4) The port used by the application to send the stream to Tor --
//			that is, the SOCKSListenAddress or TransListenAddress that the
//			application used, if we have more than one.
//		 5) The SOCKS username and password, if any.
//		 6) The source address and port for the application.

//	   We propose to use 3, 4, and 5 as a backchannel for applications to
//	   tell Tor about different sessions.  Rather than running only one
//	   SOCKSPort, a Tor user who would prefer better session isolation should
//	   run multiple SOCKSPorts/TransPorts, and configure different
//	   applications to use separate ports. Applications that support SOCKS
//	   authentication can further be separated on a single port by their
//	   choice of username/password.  Streams sent to separate ports or using
//	   different authentication information should never be sent over the
//	   same circuit.  We allow each port to have its own settings for
//	   isolation based on destination port, destination address, or both.
func Create() {

	// random integer
	num := rand.Intn(0x7fffffff-10000) + 10000
	// base url
	proxybase := "socks5://%d:x@199.241.139.7:9050"
	// proxy url with random credentials
	proxyURL := fmt.Sprintf(proxybase, num)
	// set proxy url
	err := req.SetProxyUrl(proxyURL)
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s Failed to set proxy url  %s\n", red("[-]"), err)
	}
	// check ip with aws
	url := "https://checkip.amazonaws.com"

	r, err := req.Get(url)

	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(r.Response().Body)
	body := string(bodyBytes)
	fmt.Print(body)

	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s Failed to make request  %s\n", red("[-]"), err)
	}
	// response
	resp := r.Response()
	// dsiplay
	if resp.StatusCode == 200 {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s Response Code: %s \n", green("[+]"), resp.Status)
	} else {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s Response code  %s \n", red("[-]"), resp.Status)
	}

}
