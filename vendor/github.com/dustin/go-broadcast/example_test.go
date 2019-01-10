package broadcast

import "log"

// Example of a simple broadcaster sending numbers to two workers.
//
// Five messages are sent.  The first worker prints all five.  The second worker prints the first and then unsubscribes.
func Example() {
	b := NewBroadcaster(100)

	workerOne(b)
	workerTwo(b)

	for i := 0; i < 5; i++ {
		log.Printf("Sending %v", i)
		b.Submit(i)
	}
	b.Close()
}

func workerOne(b Broadcaster) {
	ch := make(chan interface{})
	b.Register(ch)
	defer b.Unregister(ch)

	// Dump out each message sent to the broadcaster.
	go func() {
		for v := range ch {
			log.Printf("workerOne read %v", v)
		}
	}()
}

func workerTwo(b Broadcaster) {
	ch := make(chan interface{})
	b.Register(ch)
	defer b.Unregister(ch)
	defer log.Printf("workerTwo is done\n")

	go func() {
		log.Printf("workerTwo read %v\n", <-ch)
	}()
}
