package mp1

import (
	"fmt"
	"sync"
)

func mapper(in <-chan string, out chan<- map[string]int) {
	count := make(map[string]int)
	for word := range in {
		count[word] = count[word] + 1
	}

	out <- count
	close(out)
}

func reducer(in <-chan int, out chan<- float32) {
	sum, count := 0, 0
	for n := range in {
		sum += n
		count++
	}

	out <- float32(sum) / float32(count)
	close(out)
}

func inputReader(out [3]chan<- string) {
	input := [][]string{
		{"noun", "verb", "verb", "noun", "noun"},
		{"verb", "verb", "verb", "noun", "noun", "verb"},
		{"noun", "noun", "verb", "noun"},
	}

	for i := range out {
		go func(ch chan<- string, word []string) {
			for _, w := range word {
				ch <- w
			}
			close(ch)
		}(out[i], input[i])
	}
}

func shuffler(in []<-chan map[string]int, out [2]chan<- int) {
	var wg sync.WaitGroup
	wg.Add(len(in))

	for _, ch := range in {
		go func(c <-chan map[string]int) {
			for m := range c {
				nc, ok := m["noun"]
				if ok {
					out[0] <- nc
				}

				vc, ok := m["verb"]
				if ok {
					out[1] <- vc
				}
			}
			wg.Done()
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out[0])
		close(out[1])
	}()
}

func outputWriter(in []<-chan float32) {
	var wg sync.WaitGroup
	wg.Add(len(in))

	name := []string{"noun", "verb"}
	for i := 0; i < len(in); i++ {
		go func(n int, c <-chan float32) {
			for avg := range c {
				fmt.Printf("Average number of %ss per input text: %f\n", name[n], avg)
			}
			wg.Done()
		}(i, in[i])
	}

	wg.Wait()
}
