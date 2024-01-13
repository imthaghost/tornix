<p align="center">
    <img alt="net" src="docs/media/network.png"> 
</p>
<p align="center">
Tornix is a user-friendly tool that enhances privacy and security on the Tor network by providing easy-to-manage stream isolation, 
enabling users to route their internet traffic securely through separate Tor circuits for discrete, protected online activities.

</p>
<p align="center">
   <a href="https://goreportcard.com/report/github.com/imthaghost/tornix"><img src="https://goreportcard.com/badge/github.com/imthaghost/tornix"></a>
   <a href="https://travis-ci.org/imthaghost/tornix.svg?branch=master"><img src="https://travis-ci.org/imthaghost/tornix.svg?branch=master"></a>

</p>
<br>


# Usage
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

