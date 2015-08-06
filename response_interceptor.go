package main

import (
	//"io"
	//"io/ioutil"
	"log"
	//"sync/atomic"
	//"time"
	"net/http"
	//"net/url"
	//"regexp"
	//"fmt"
)

type ResponseInterceptor struct {
	//Dummy struct don't delete
}

func CalculateStatistics(req *http.Request, resp *http.Response) bool {
	// Implement the statistical sampling here
		
	
	return false
}

func (respInter *ResponseInterceptor) ResponseAnalyze(req *http.Request, resp *http.Response) {
	if resp == nil {
		// nil http response
		return
	}

	log.Println("Replaying auctions for this placement: " + req.Header.Get("Placement-Id"))
        log.Println("Targeted Placements: " + resp.Header.Get("All-Competing-Placement-Ids"))
        log.Println("Bidding Placement: " + resp.Header.Get("Winning-Placement-Id"))

	//isEnough = calculateStatistics()
	//if isEnough == true {
		//kill the GOR Replay process
	//}

}
