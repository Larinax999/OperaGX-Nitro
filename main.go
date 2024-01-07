package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func GenToken() string {
	for {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		req.Header.SetMethod("POST")
		req.SetRequestURI("https://api.discord.gx.games/v1/direct-fulfillment")
		req.Header.SetContentType("application/json")
		req.SetBody([]byte(fmt.Sprintf(`{"partnerUserId":"%s"}`, uuid.New().String())))
		fasthttp.Do(req, resp)
		a := string(resp.Body())
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)
		if resp.StatusCode() == 429 {
			fmt.Println("[!] Got Rate limit. Change ur fucking ip")
			continue
		}
		ac := strings.Split(a, `"`)
		if len(ac) == 1 {
			fmt.Println("[!] unknown err")
			continue
		}
		return strings.Split(ac[3], `"`)[0]
	}
}

func main() {
	fmt.Println("Niga")
	file, err := os.OpenFile("promos.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	tokens := make(chan string)
	defer file.Close() // Close the file at the end of the function
	go func() {
		for {
			_, err = file.Write([]byte(<-tokens))
			if err != nil {
				fmt.Println("Error Write file:", err)
				return
			}
		}
	}()
	for i := 0; i < 100; i++ {
		go func() {
			for {
				tokens <- GenToken() + "\n"
			}
		}()
	}
	select {}
}
