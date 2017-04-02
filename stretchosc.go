package stretchosc

import (
	"fmt"
	"sync"

	"github.com/hypebeast/go-osc/osc"
)

// Message represents a message that we send to the time stretcher
type Message struct {
	Addr string
	Arg  interface{}
}

// TimeStretchControl wraps an osc client. Messages on the send channel will be
// sent to the host via the `client` member.
type TimeStretchControl struct {
	send      chan Message
	client    *osc.Client
	WaitGroup sync.WaitGroup
}

// Send a message. Safe for concurrent calls
func (tsc *TimeStretchControl) Send(addr string, arg interface{}) {
	tsc.WaitGroup.Add(1)
	tsc.send <- Message{
		Addr: addr,
		Arg:  arg,
	}
}

// Close stops listening for messages.
func (tsc *TimeStretchControl) Close() {
	close(tsc.send)
}

// MakeTimeStretchControl creates a TimeStretchControl struct that listens for
// Messages on Send channel, and sends those the the host specified by ip, port.
func MakeTimeStretchControl(ip string, port int) *TimeStretchControl {
	tsc := &TimeStretchControl{
		client: osc.NewClient(ip, port),
		send:   make(chan Message, 4),
	}

	go func() {
		for {
			select {
			case msg, ok := <-tsc.send:
				if !ok {
					return
				}
				oscMsg := osc.NewMessage(msg.Addr)
				oscMsg.Append(msg.Arg)
				tsc.client.Send(oscMsg)
				tsc.WaitGroup.Done()
			}
		}
	}()
	return tsc
}

func (tsc *TimeStretchControl) setToggle(stretchNum int, enable bool) {
	if stretchNum < 1 {
		fmt.Println("TimeStretchControl received invalid toggle request")
		return
	}
	oscAddr := fmt.Sprintf("/1/toggle%d", stretchNum)

	val := int32(0)
	if enable {
		val = 1
	}
	tsc.Send(oscAddr, val)
}

// Activate a stretcher
func (tsc *TimeStretchControl) Activate(stretchNum int) {
	tsc.setToggle(stretchNum, true)
}

// Deactivate a stretcher
func (tsc *TimeStretchControl) Deactivate(stretchNum int) {
	tsc.setToggle(stretchNum, false)
}

// StretchAmount sets the stretch amount for stretcher indexed by stretchNum
func (tsc *TimeStretchControl) StretchAmount(stretchNum int, amt float32) {
	if stretchNum < 1 || amt <= 0 {
		fmt.Println("TimeStretchControl received invalid stretch request")
		return
	}
	oscAddr := fmt.Sprintf("/1/fader%d", stretchNum)
	tsc.Send(oscAddr, amt)
}
