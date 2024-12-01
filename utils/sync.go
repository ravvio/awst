package utils

import "sync"

func WithSemaphore(i int, wg *sync.WaitGroup, semaphore chan struct{}, callback func()) {
	defer wg.Done()
	callback()
	// Free channel
	<- semaphore
}
