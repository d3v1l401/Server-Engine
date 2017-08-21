package main

import (
	"fmt"
	"io/ioutil"

	"./common"
	"./server"
	"./server/session"
	"./server/ssettings"
)

func OnNewClient(s *server.Server, c *session.Instance) {
	fmt.Printf("New client %s (%s) @ %s\n", c.GetID(), c.GetIP(), s.GetName())
	c.SetData(make(map[string]string))
	lel := c.AccessData()
	lel.(map[string]string)["Auth"] = "no"
}

func OnClientDeath(s *server.Server, c *session.Instance) {
	//fmt.Printf("Client %s (%s) @ %s being removed\n", c.GetID(), c.GetIP(), "server#1")
}

func BeforePacketSend(client *session.Instance, buffer []byte) {
	fmt.Printf("Packet being sent to %s...\n", client.GetID())
}

func OnTick(c *session.Instance) {
	//fmt.Printf("Client %s (%s) @ %s is dead\n", c.GetID(), c.GetIP(), "server#1")
}

func OnHeartbeat(s *server.Server) {
	//fmt.Printf("Online clients: %d.\n", s.OnlineClients())
}

func OnPacketRecv(client *session.Instance, buffer []byte) bool {
	if buffer != nil {
		fmt.Printf("Request received from %s: %s\n", client.GetID(), string(buffer))

		response, er := ioutil.ReadFile(string(buffer[:len(buffer)-2]))
		if er != nil {
			client.Send([]byte(er.Error()))
		}
		client.Send([]byte(string(response) + "\n"))

		if buffer[0] == 0x1B {
			client.Destroy()
		}

		return true
	}
	return false
}

func OnPacketSent(client *session.Instance) {
	fmt.Printf("Packet sent to %s\n", client.GetID())
}

func OnDestroy(client *session.Instance) {
	fmt.Printf("Client %s destroyed\n", client.GetID())
}

func main() {
	common.Print("Starting server...", common.MTYPE_NORMAL)
	ssettings.Initialize(1024, 10000000, true)

	cbSets := &session.CallbackSet{
		BeforePacketSend: BeforePacketSend,
		OnPacketReceive:  OnPacketRecv,
		OnPacketSent:     OnPacketSent,
		OnDestroy:        OnDestroy,
		OnClientTick:     OnTick,
	}

	tcpServer := &server.Server{}
	if tcpServer.Initialize("d3vPC", "127.0.0.1", 1337, server.PROTOCOL_TCP) {
		tcpServer.OnNewClient = OnNewClient
		tcpServer.OnClientDeath = OnClientDeath
		tcpServer.OnHeartbeat = OnHeartbeat
		tcpServer.SetClientCallbacks(cbSets)
		common.Print("Server started.", common.MTYPE_NORMAL)
		tcpServer.Start()
		for {

		}
	}

	common.Print("Server stopped.", common.MTYPE_WARNING)
}
