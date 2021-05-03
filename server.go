package main

import (
	atmi "github.com/endurox-dev/endurox-go"
        "fmt"
        "os"
)

//Service func
//Here svc contains the caller infos
func TPDISCOVER(ac *atmi.ATMICtx, svc *atmi.TPSVCINFO) {


        ac.TpReturn(atmi.TPSUCCESS, 0, &svc.Data, 0)

}

//Server boot/init
func Init(ac *atmi.ATMICtx) int {

        //Advertize TPDISCOVER

	startReadWrite()
	
        if err := ac.TpAdvertise("TPDISCOVER", "TPDISCOVER", TPDISCOVER); err != nil {
                fmt.Println(err)
                return atmi.FAIL
        }

        return atmi.SUCCEED
}

//Server shutdown
func Uninit(ac *atmi.ATMICtx) {
        fmt.Println("Server shutting down...")
}


//Server main
func main() {
        //Have some context
        ac, err := atmi.NewATMICtx()

        if nil != err {
                fmt.Errorf("Failed to allocate cotnext!", err)
                os.Exit(atmi.FAIL)
        } else {
                //Run as server
                ac.TpRun(Init, Uninit)
        }
}

