package main

import (
	"fmt"
	"github.com/etsme-com/ssf/v2/base/envcfg"
	"time"
)

//var log *logrus.Logger

func init() {

}

func main() {
	// fmt.Println("ssferr.ErrEmpty = ", ssferr.ErrEmpty)
	a, b := envcfg.GetRunEnv()
	for {
		time.Sleep(time.Duration(3) * time.Second)
		//ssf.Logger.Infoln("test--------- log:", ssferr.ErrIsNil)
		// fmt.Println("ssferr.ErrIsNil = ", ssferr.ErrIsNil)
		fmt.Println("ssferr.ErrEmpty = ", a, b)
	}
}
