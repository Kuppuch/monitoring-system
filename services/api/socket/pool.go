package socket

type Pool struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewPool() *Pool {
	return &Pool{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (p *Pool) Run() {
	for {
		select {
		case client := <-p.register:
			p.clients[client] = true
		case client := <-p.unregister:
			if _, ok := p.clients[client]; ok {
				delete(p.clients, client)
				close(client.send)
			}
		case message := <-p.broadcast:
			for client := range p.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(p.clients, client)
				}
			}
		}
	}
}
