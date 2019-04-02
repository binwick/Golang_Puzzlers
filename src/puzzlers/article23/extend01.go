package main

import (
	"fmt"
	"sync"
)

type Button struct {
	//1
	Clicked *sync.Cond
}

func main() {

	lock := &sync.Mutex{}
	button := Button{Clicked: sync.NewCond(lock)}
	var box uint
	subscribe := func(c *sync.Cond, fn func()) { //2
		go func() {
			lock.Lock()
			for box == 0 {
				c.Wait()
			}
			lock.Unlock()
			fn()
		}()
	}

	var wg sync.WaitGroup //3
	wg.Add(3)
	subscribe(button.Clicked, func() { //4
		fmt.Println("Maximizing window.")
		wg.Done()
	})
	subscribe(button.Clicked, func() { //5
		fmt.Println("Displaying annoying dialog box!")
		wg.Done()
	})
	subscribe(button.Clicked, func() { //6
		fmt.Println("Mouse clicked.")
		wg.Done()
	})

	fmt.Println("before click.")
	lock.Lock()
	box = 1
	lock.Unlock()
	button.Clicked.Broadcast() //7

	wg.Wait()
}
