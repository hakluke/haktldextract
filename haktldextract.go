package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/hakluke/tldextract"
)

func main() {
	concurrencyPtr := flag.Int("t", 8, "Number of threads to utilise. Default is 8.")
	subdomainsPtr := flag.Bool("s", false, "dump subdomains instead of base domains")
	flag.Parse()

	cache := "/tmp/tld.cache"
	extract, err := tldextract.New(cache, false)
	if err != nil {
		fmt.Println(err)
	}

	numWorkers := *concurrencyPtr
	work := make(chan string)
	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			work <- s.Text()
		}
		close(work)
	}()

	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go doWork(work, wg, *subdomainsPtr, extract)
	}
	wg.Wait()
}

func doWork(work chan string, wg *sync.WaitGroup, subdomainsPtr bool, extract *tldextract.TLDExtract) {
	for url := range work {
		result := extract.Extract(url)
		if subdomainsPtr && len(result.Sub) > 0 {
			fmt.Println(result.Sub + "." + result.Root + "." + result.Tld)
		} else {
			fmt.Println(result.Root + "." + result.Tld)
		}
	}
	wg.Done()
}
