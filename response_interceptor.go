package main

import (
	"os"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)
var TOTAL_AUCTION_NUMBER int = 0
var TARGETED_AUCTIONS_NUMBER int = 1
var BIDDED_AUCTIONS_NUMBER int = 1
var KILL_TIME = time.Now().Add(time.Second*5)

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
	targeted_ratio_sofar := float64(TARGETED_AUCTIONS_NUMBER)/float64(TOTAL_AUCTION_NUMBER) * 100
	bidded_ratio_sofar := float64(BIDDED_AUCTIONS_NUMBER)/float64(TOTAL_AUCTION_NUMBER) * 100	

	log.Println("Total number of auctions: " + strconv.Itoa(TOTAL_AUCTION_NUMBER))
        log.Println("Targeted Ratio for placement " + replayed_placement + " is: "+ strconv.Itoa(TARGETED_AUCTIONS_NUMBER))
	log.Println("Bidded Ratio for placement " + replayed_placement + " is: "+ strconv.Itoa(BIDDED_AUCTIONS_NUMBER))

	if time.Now().After(KILL_TIME) {
		log.Println("Writing to file before killing the process")
		outfile, err := os.Create("placement_stats_" + replayed_placement + ".txt")
		if err != nil {
			log.Println("THERE I HAVE SEEN THE ERROR")
			log.Fatal(err)
		}
		if _, err := outfile.WriteString("Total number of auctions: " + strconv.Itoa(TOTAL_AUCTION_NUMBER) + "\n"); err != nil {
			panic(err)
		}
		if _, err := outfile.WriteString("Targeted Auctions for placement " + replayed_placement + " is: "+ strconv.Itoa(TARGETED_AUCTIONS_NUMBER) + "\n"); err != nil {
			panic(err)
		}
		if _,err := outfile.WriteString("Bidded Auctions for placement " + replayed_placement + " is: "+ strconv.Itoa(BIDDED_AUCTIONS_NUMBER) + "\n"); err != nil {
		    panic(err)
		}
		if _,err := outfile.WriteString("Replay Completion Timestamp:  " + time.Now().Format("20060102150405") + "\n"); err != nil {
		        panic(err)
		}
		DAILY_AUCTIONS_TOTAL := float64(45000000000)
                SCALED_TARGETED := DAILY_AUCTIONS_TOTAL * targeted_ratio_sofar
                SCALED_BIDDING := DAILY_AUCTIONS_TOTAL * bidded_ratio_sofar
		
		if _, err := outfile.WriteString("Daily Scaled Total number of auctions: " + strconv.FormatFloat(DAILY_AUCTIONS_TOTAL, 'f', 0, 64) + "\n"); err != nil {
			panic(err)
		}
		if _, err := outfile.WriteString("Daily Scaled Targeted number of auctions: " + strconv.FormatFloat(SCALED_TARGETED, 'f', 0, 64) + "\n"); err != nil {
			panic(err)
		}
		if _, err := outfile.WriteString("Daily Scaled Bidded number of auctions: " + strconv.FormatFloat(SCALED_BIDDING, 'f', 0, 64) + "\n"); err != nil {
			panic(err)
		}
	
		outfile.Close()
		log.Println("Killing the process due to overtime")
		
		//log.Println("Total number of auctions: " + strconv.Itoa(TOTAL_AUCTION_NUMBER))
		//log.Println("Targeted Ratio for placement " + replayed_placement + " is: "+ strconv.Itoa(TARGETED_AUCTIONS_NUMBER))
		//log.Println("Bidded Ratio for placement " + replayed_placement + " is: "+ strconv.Itoa(BIDDED_AUCTIONS_NUMBER))
		os.Exit(0)
	}
}

func (respInter *ResponseInterceptor) ResponseAnalyze(req *http.Request, resp *http.Response) {
	if resp == nil {
		// nil http response
		return
	}
	
	CalculateStatistics(req, resp)
}
