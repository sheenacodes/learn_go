package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	str  string
	wait chan bool
}

func boring(msg string) <-chan Message {
	waitForIt := make(chan bool)
	c := make(chan Message)
	go func() {
		for i := 0; i < 10; i++ {
			text := fmt.Sprintf("boring...%s %d", msg, i)
			fmt.Println(text)
			c <- Message{text, waitForIt}
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			<-waitForIt
		}

	}()
	return c
}

func fanIn(inputs ...<-chan Message) <-chan Message {
	c := make(chan Message)
	for _, input := range inputs {
		go func() {
			for msg := range input {
				fmt.Printf("multiplex %s \n", msg.str)
				c <- msg
			}
		}()
	}
	return c
}

func main() {

	c := fanIn(boring("JOE"), boring("ANN"))
	fmt.Println("main listening")
	for i := 0; i < 10; i++ {
		msg1 := <-c
		fmt.Printf("main reading %s \n", msg1.str)
		msg2 := <-c
		fmt.Printf("main reading %s \n", msg2.str)
		msg1.wait <- true
		msg2.wait <- true
	}
	time.Sleep(time.Second)

	// goroutine is not a thread. infact they are much lighter and mulitplexed onto threads.
	//they stacks grow and shrink as needed

	//dont communicate by sharing memory. share memory to comunicate.

	// concurrency patterns
	//generator - func return a channel - ie a function returns a channel
	//through which we can communicate to the service it provides

	//fan in or multiplexer - channel reads from multiple channels.

	//restoring sequencing - pass a channel on a channel

}
