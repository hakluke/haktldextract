package main

import (
        "github.com/joeguo/tldextract"
        "bufio"
	"fmt"
	"os"
)


func main() {
        concurrency := 16 
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
                        fmt.Println(result.Root + "." + result.Tld)
                    }
		}(scanner.Text())
	}
}
