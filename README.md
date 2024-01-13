<p align="center">
    <img alt="net" src="docs/media/network.png"> 
</p>
<p align="center">
Tornix, an advanced yet user-friendly tool, revolutionizes the way users interact with the Tor network by offering sophisticated stream isolation capabilities.
Designed for both privacy enthusiasts and security professionals, Tornix simplifies the process of securely routing various types of internet traffic through distinct Tor circuits. 
This method not only enhances privacy by preventing the correlation of different activities but also bolsters security by isolating each data stream. 
With Tornix, users gain the power to efficiently manage and protect their online presence on the Tor network, ensuring each action remains discreet and secure.

</p>
<p align="center">
   <a href="https://goreportcard.com/report/github.com/imthaghost/tornix"><img src="https://goreportcard.com/badge/github.com/imthaghost/tornix"></a>
   <a href="https://travis-ci.org/imthaghost/tornix.svg?branch=master"><img src="https://travis-ci.org/imthaghost/tornix.svg?branch=master"></a>

</p>
<br>


# Usage
```go

import (
	"github.com/imthaghost/tornix"
)

```
```bash
go get "github.com/imthaghost/tornix"

```

# Examples

```go
package main

import (
	"io"
	"log"

	"github.com/imthaghost/tornix"
)

func main() {
	// include the max concurrent sessions
	manager := tornix.NewManager(10)
	// start a new session
	session, client, err := manager.StartNewSession()
	if err != nil {
		log.Fatal(err)
	}

	// check your client's IP with AWS
	resp, err := client.Get("https://checkip.amazonaws.com")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
	// 185.220.100.254
}
```

