# usps

Example of calculating shipping rate for one package.

```
package main

import (
	"encoding/xml"
	"fmt"
	"github.com/dougkirkley/usps"
	"github.com/dougkirkley/usps/rate"
	"log"
)

func main() {
	u := &usps.USPS{
		Username: "USERID",
	}
	rates := []rate.Request{
		{
			XMLName:        xml.Name{Local: "Package"},
			ID:             "1st",
			Service:        "PRIORITY",
			ZipOrigination: "21212",
			ZipDestination: "20759",
			Pounds:         "0",
			Ounces:         "12",
			Container:      "VARIABLE",
		},
	}
	cost, err := rate.NewRate(*u).Calculate(rates)
	if err != nil {
		log.Print(err)
	}
	fmt.Print(cost)
}
```
