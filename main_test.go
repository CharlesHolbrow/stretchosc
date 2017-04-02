package oscstretch

import (
	"fmt"
	"testing"
	"time"
)

func Test_main(t *testing.T) {

	fmt.Println("OK")
	tsc := MakeTimeStretchControl("localhost", 12340)
	tsc.Activate(3)
	tsc.StretchAmount(3, 20)
	time.Sleep(time.Second * 5)
	tsc.StretchAmount(3, 1)
	time.Sleep(time.Second * 5)
	tsc.Deactivate(3)
	tsc.WaitGroup.Wait()
	fmt.Println("Done")

}
