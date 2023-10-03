package main

import (
	"fmt"
	"sync"
	"time"
)



func main(){

	tick := make(chan int, 2)
	done := make(chan interface{})
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func ()  {
		defer wg.Done()
		
		for {
			select {
			case c := <- tick:
				fmt.Printf("time: %v Count: %d\n",time.Now(), c)
				time.Sleep(time.Second * 3)
				
			case <- done:
				fmt.Println("channel closed")
				return
				

			}

			
			

			
		}
		
	}()


	for i := 0; i < 10; i++{
		fmt.Printf("Sleeping loop: %d\n", i)
		time.Sleep(time.Second * 2)
		tick<- i
	}

	close(done)

	wg.Wait()


	
}