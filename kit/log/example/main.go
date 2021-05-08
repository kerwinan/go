package main

import (
	"fmt"
	"kerwinan/go/kit/log"
	"time"
)

func main() {
	err := log.Init("../config/log.yaml")
	if err != nil {
		fmt.Printf("init log failed: %v\n", err)
		return
	}
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		log.X.Debugf("hello debug")
		log.X.Infof("hello info")
	}

}
