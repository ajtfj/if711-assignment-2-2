package main

import (
	"fmt"
	"os"
	"sync"
)

const (
	CAP = 10
)

type Client int

type Barbershop struct {
	queue []Client
	cap   int
	cond  *sync.Cond
}

func NewBarbershop(cap int) *Barbershop {
	barbershop := Barbershop{
		cap: cap,
		cond: &sync.Cond{
			L: &sync.Mutex{},
		},
	}

	return &barbershop
}

func (b *Barbershop) GoToBarbershop(client Client) error {
	b.cond.L.Lock()
	if len(b.queue) == b.cap {
		b.cond.L.Unlock()
		return fmt.Errorf("Barbershop is at full capacity")
	}

	b.queue = append(b.queue, client)
	b.cond.Broadcast()
	b.cond.L.Unlock()

	return nil
}

func (b *Barbershop) CutHair() *Client {
	b.cond.L.Lock()
	for len(b.queue) == 0 {
		fmt.Println("No client, barber will sleep")
		b.cond.Wait()
	}

	client := b.queue[0]
	b.queue = b.queue[1:]

	b.cond.Broadcast()
	b.cond.L.Unlock()

	return &client

}

func main() {
	barbershop := NewBarbershop(CAP)
	var event string
	var id int
	for {
		fmt.Scanf("%s %d", &event, &id)
		switch event {
		case "GoToBarbershop":
			fmt.Printf("Client %d is going to barbershop\n", id)
			client := Client(id)
			if err := barbershop.GoToBarbershop(client); err != nil {
				fmt.Println(err)
			}
		case "CutHair":
			go func() {
				client := barbershop.CutHair()
				fmt.Printf("Client %v had a haircut\n", *client)
			}()
		default:
			fmt.Println("Undefined event")
			os.Exit(1)
		}
	}
}
