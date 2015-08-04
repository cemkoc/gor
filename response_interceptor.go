package main

import (
	"os"
	"io"
	"io/ioutil"
	"log"
	"sync/atomic"
	"time"
	"net/http"
	"net/url"
	"regexp"
	"fmt"
)

type ResponseInterceptor struct {
	ReqUrl               string
	ReqUserAgent         string
	ReqGorFlag	     string
	ReqPlacementId       string
	RespStatus           string
	RespStatusCode       int
	RespCompetingPlacements string
	RespWinningPlacement string
}

func (respInter *ResponseInterceptor) CalculateStatistics() bool {
	// Implement the statistical sampling here
	return false
}

func (respInter *ResponseInterceptor) ResponseAnalyze(req *http.Request, resp *http.Response) {
	if resp == nil {
		// nil http response - skipped elasticsearch export for this request
		return
	}
	t := time.Now()
	
	irr := ResponseInterceptor{
		ReqUrl:               req.URL.String(),
		ReqUserAgent:         req.UserAgent(),
		ReqGorFlag:	      req.Header.Get("Is-Coming-From-Gor"),
		ReqPlacementId:	      req.Header.Get("Placement-Id"),
		RespStatus:           resp.Status,
		RespStatusCode:       resp.StatusCode,
		Timestamp:            t,
		RespCompetingPlacements: resp.Header.Get("All-Competing-Placement-Ids"),
		RespWinningPlacement: resp.Header.Get("Winning-Placement-Id"),
	}
	
	isEnough = CalculateStatistics()
	if isEnough == true {
		//kill the GOR Replay process
	}
}
