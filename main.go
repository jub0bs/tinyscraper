package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func main() {
	hosts := make(chan string)
	var wg sync.WaitGroup

	const nbWorkers = 2
	wg.Add(nbWorkers)
	for i := 0; i < nbWorkers; i++ {
		go worker(hosts, &wg)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		host := scanner.Text()
		hosts <- host
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	close(hosts)
	wg.Wait()
}

func worker(hosts <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for host := range hosts {
		printContentLengthOfResp(host)
	}
}

func printContentLengthOfResp(host string) {
	resp, err := http.Get("http://" + host)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer resp.Body.Close()
	contentLength, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(contentLength)
}
