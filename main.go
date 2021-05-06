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
	var wg sync.WaitGroup

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		host := scanner.Text()
		wg.Add(1)
		go printContentLengthOfResp(host, &wg)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	wg.Wait()
}

func printContentLengthOfResp(host string, wg *sync.WaitGroup) {
	defer wg.Done()
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
