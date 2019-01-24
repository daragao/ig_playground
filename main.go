package main

import (
	"fmt"

	"github.com/daragao/ig_trade/exchange"
	util "github.com/daragao/ig_trade/utils"
)

func main() {

	conf, err := util.ReadConfig("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	ig := exchange.NewIGClient(conf.Username, conf.Password, conf.APIKey)
	fmt.Printf("%#v\n", ig)

	accounts := ig.Accounts()
	fmt.Printf("%#v\n", accounts)
}
