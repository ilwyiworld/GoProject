package main
import (
	"fmt"
	"sync"
	"time"
)
func main() {
	rw := new(sync.RWMutex)
	for i := 0; i < 2; i++ {   // 建立两个写者
		go func(i int) {
			for j := 0; j < 3; j++ {
				rw.Lock()
				fmt.Printf("第%v个写\n", i)
				rw.Unlock()
			}
		}(i)
	}
	for i := 0; i < 2; i++ {    // 建立两个读者
		go func(i int) {
			for j := 0; j < 3; j++ {
				rw.RLock()
				fmt.Printf("第%v个读\n", i)
				rw.RUnlock()
			}
		}(i)
	}
	time.Sleep(time.Second)
	fmt.Println("Done")
}