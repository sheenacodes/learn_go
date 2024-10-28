package main

import (
	"fmt"
	"math/rand"
	"time"
)

package main

import (
"fmt"
"math/rand"
"time"
)

var (
	Web    = fakeSearch("web")
	Web1   = fakeSearch("web1")
	Web2   = fakeSearch("web2")
	Image  = fakeSearch("image")
	Image1 = fakeSearch("image1")
	Image2 = fakeSearch("image2")
	Video1 = fakeSearch("video1")
	Video2 = fakeSearch("video2")
	Video  = fakeSearch("video")
)

type Result string
type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		fmt.Printf("********** start goroutine %s %q \n", kind, query)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q \n", kind, query))
	}
}

func Google(query string) (results []Result) {
	// serially append to index
	results = append(results, Web(query))
	results = append(results, Image(query))
	results = append(results, Video(query))
	return
}

func GoGoogle(query string) (results []Result) {
	//fanin pattern
	c := make(chan Result)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	for i := 0; i < 3; i++ {
		result := <-c
		results = append(results, result)
	}
	return
}

func GoGoogleTimeout(query string, timeout int) (results []Result) {
	//timeout pattern
	c := make(chan Result)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	timeoutAfter := time.After(time.Duration(timeout) * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeoutAfter:
			fmt.Println("timedout after", timeout)
			return
		}

	}
	return
}

func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) { c <- replicas[i](query) }

	for i := range replicas {
		go searchReplica(i)
	}

	return <-c
}

func ReplicaGoogle(query string) (results []Result) {
	c := make(chan Result)
	go func() { c <- First(query, Web1, Web2) }()
	go func() { c <- First(query, Image1, Image2) }()
	go func() { c <- First(query, Video1, Video2) }()
	timeoutAfter := time.After(time.Duration(80) * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeoutAfter:
			fmt.Println("timedout")
			return
		}
	}
	return
}

func main_search() {
	rand.Seed(time.Now().UTC().UnixNano())
	start := time.Now()
	results := GoGoogle("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)

	start = time.Now()
	results = GoGoogleTimeout("hello", 80)
	elapsed = time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)

	start = time.Now()
	results = ReplicaGoogle("replica")
	elapsed = time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)

}
