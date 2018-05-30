package main

import (
	"github.com/sakjur/telegraf/pkg/smsgw"
	"github.com/sakjur/telegraf/pkg/smsgw/elks"
)

func main() {
	message := smsgw.Message{From: "Telegraf", To: "+46700000000", Message: "Hello, upper east side"}
	resp, err := elks.Send(message)

	if err != nil {
		println(err.Error())
	}

	println("Sent message for " + resp.CostToString() + ".")
}
