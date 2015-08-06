package main

import (
	//"io"
	//"io/ioutil"
	"log"
	"net/http"
	"regexp"
)
var TOTAL_AUCTION_NUMBER int = 0
var TARGETED_AUCTIONS_NUMBER int = 0
var BIDDED_AUCTIONS_NUMBER int = 0

type ResponseInterceptor struct {
	//Dummy struct don't delete
}

func CalculateStatistics(req *http.Request, resp *http.Response) bool {
	// Implement the statistical sampling here
	replayed_placement := req.Header.Get("Placement-Id")
	competing_placements_list := resp.Header.Get("All-Competing-Placement-Ids")
	bidding_placement := resp.Header.Get("Winning-Placement-Id")
	
		

	
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
