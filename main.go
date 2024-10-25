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
		for i := 0; i < 5; i++ {
			text := fmt.Sprintf(" * boring...%s %d \n", msg, i)
			fmt.Println(text)
			c <- Message{text, waitForIt}
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			<-waitForIt
		}

	}()
	return c
}

func fanIn_old(inputs ...<-chan Message) <-chan Message {
	c := make(chan Message)
	for _, input := range inputs {
		go func() {
			for msg := range input {
				fmt.Printf(" * * multiplex %s \n", msg.str)
				c <- msg
			}
		}()
	}
	return c
}

func fanIn_select(input1, input2 <-chan Message) <-chan Message {
	c := make(chan Message)

	// only one go routine
	go func() {
		for {
			select {
			case s := <-input1:
				{
					fmt.Printf(" * * multiplex %s \n", s.str)
					c <- s
				}
			case s := <-input2:
				{
					fmt.Printf(" * * multiplex %s \n", s.str)
					c <- s
				}
			}
		}
	}()
	return c
}

func main() {

	c := fanIn_select(boring("JOE"), boring("ANN"))
	fmt.Println("main listening")
	/*	for i := 0; i < 5; i++ {
		msg1 := <-c
		fmt.Printf(" * * * main reading %s \n", msg1.str)
		msg2 := <-c
		fmt.Printf(" * * * main reading %s \n", msg2.str)
		msg1.wait <- true
		msg2.wait <- true
	}*/

	var time1 time.Time

	for {
		select {
		case msg := <-c:
			fmt.Printf(" * * * main select %s \n", msg.str)
			msg.wait <- true
			time1 = currentTime()
		case <-time.After(3 * time.Second):
			fmt.Println("main timeout")
			time2 := currentTime()
			fmt.Println(time2.Sub(time1))
			return
		}
	}

	// goroutine is not a thread. in fact they are much lighter and multiplexed onto threads.
	//they stacks grow and shrink as needed

	//dont communicate by sharing memory. share memory to communicate.

	// concurrency patterns
	//generator - func return a channel - ie a function returns a channel
	//through which we can communicate to the service it provides

	//fan in or multiplexer - channel reads from multiple channels.

	//restoring sequencing - pass a channel on a channel

	//select statememt is a control statement - like a switch - controls program execution based on what you receive
	//reason why goroutine is not a library because hard to do control statements that depend on libraries

	//each case for this switch instead of being an expression is a communication.
	// it is something you receive from a channel

	// select blocks until one communication can proceed - if multiple channels communicate -
	//select chooses pseudo ramdomly
	// a default will execute immediately if not channel is ready and makes select - non blocking!
	// like djikstra's guarded commands - several if statements - which one is chosen at runtime

	//select can be used to timeout a communication

}

func currentTime() time.Time {
	dt := time.Now()
	fmt.Println("Current date and time is: ", dt.String())
	return dt
}
