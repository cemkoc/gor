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

func CalculateStatistics(req *http.Request, resp *http.Response) {
	// Implement the statistical sampling here
	TOTAL_AUCTION_NUMBER++
	replayed_placement := req.Header.Get("Placement-Id")
	competing_placements_list := resp.Header.Get("All-Competing-Placement-Ids")
	bidding_placement := resp.Header.Get("Winning-Placement-Id")
	
	re := regexp.MustCompile("\\b" + replayed_placement + "\\b")

	if len(re.FindAllString(competing_placements_list, -1)) >= 1 {
		TARGETED_AUCTIONS_NUMBER++
	}
	if len(re.FindAllString(bidding_placement, -1)) >= 1 {
		BIDDED_AUCTIONS_NUMBER++
	}
		
	targeted_ratio_sofar := float64(TARGETED_AUCTIONS_NUMBER)/float64(TOTAL_AUCTION_NUMBER)
	bidded_ratio_sofar := float64(BIDDED_AUCTIONS_NUMBER)/float64(TOTAL_AUCTION_NUMBER)
	
	log.Println("Targeted Ratio for placement " + replayed_placement + " is: %d percent", targeted_ratio_sofar)
	log.Println("Bidded Ratio for placement " + replayed_placement + " is: %d percent", bidded_ratio_sofar)
	 
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
