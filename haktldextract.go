package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/hakluke/tldextract"
)

var (
	domainMap = make(map[string]bool) // this map will keep track of unique domains
	domainMapMutex = &sync.Mutex{} // this mutex will make sure that reading and writing to the map is thread-safe
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
		var domain string
		if subdomainsPtr && len(result.Sub) > 0 {
			domain = result.Sub + "." + result.Root + "." + result.Tld
		} else {
			domain = result.Root + "." + result.Tld
		}
		// lock the mutex to prevent concurrent read/write
		domainMapMutex.Lock()
		if _, ok := domainMap[domain]; !ok { // check if the domain is already in the map
			domainMap[domain] = true 
			fmt.Println(domain) // print the domain
		}
		// unlock the mutex
		domainMapMutex.Unlock()
	}
	wg.Done()
}
