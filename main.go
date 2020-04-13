package main

import (
        "github.com/joeguo/tldextract"
        "bufio"
        "flag"
	"fmt"
	"os"
)


func main() {
        concurrencyPtr := flag.Int("t", 16, "number of threads to utilise")
        subdomainsPtr := flag.Bool("s", false, "dump subdomains instead of base domains") 
        flag.Parse()

        concurrency := *concurrencyPtr 
        sem := make(chan bool, concurrency)
	scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            sem <- true
		go func(url string) {
                    defer func() { <-sem }() 
                    cache := "/tmp/tld.cache"
                    extract,err := tldextract.New(cache,false)
                    if err != nil{
                        return
                    }
                    result:=extract.Extract(url)
                    if err == nil{
                        if *subdomainsPtr {
                            fmt.Println(result.Sub + "." + result.Root + "." + result.Tld)
                        } else {
                            fmt.Println(result.Root + "." + result.Tld)
                        }
                    }
		}(scanner.Text())
	}
}
