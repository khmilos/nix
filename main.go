package main

import (
	"fmt"
	"sync"
)

// Реализовать функцию что будет принимать неогр колво каналов и читать входящие сообщение из этих чанелов
// передавать в канал

func rec(in []chan int, out chan int, wg *sync.WaitGroup) {
	for _, c := range in {
		go func(c chan int) {
			for {
				v, ok := <-c
				if !ok {
					wg.Done()
					return
				}
				out <- v
			}
		}(c)
	}
}

func main() {
	out := make(chan int)
	channels := make([]chan int, 10)

	var wg sync.WaitGroup // Нужно как-то синхронизировать. Закрыть канал нельзя, так как он 1, а отправителей
	wg.Add(len(channels)) // много. По-этому передаю им sync.WaitGroup

	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan int)
	}
	go rec(channels, out, &wg)

	for i := 0; i < len(channels); i++ {
		channels[i] <- i
		close(channels[i]) // Забыл добавить закрытие канала что получает
	}

	go func() {
		// Специальная горутина, которая будет ожидать окончания всех горутин-отправителей
		// Когда все горутины будут закрыты, она закроет канал в который они отправляют
		wg.Wait()
		close(out)
	}()

	for {
		v, ok := <-out
		if !ok {
			return
		}
		fmt.Println(v)
	}
}
