package main

import "fmt"

func main() {
	usb:=PhoneConnecter{"PhoneConnecter"}
	Disconnect(usb)
}

type USB interface {
	Name() string
	Connect()
}

type PhoneConnecter struct {
	name string
}

func (pc PhoneConnecter) Name() string{
	return pc.name
}

func (pc PhoneConnecter) Connect() {
	fmt.Println("Connect:",pc.name)
}

func Disconnect(usb USB)  {
	fmt.Println("Disconnected.")
}