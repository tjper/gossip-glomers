package idgen

import (
	"encoding/json"
	"strconv"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func NewServer() *Server {
	srv := &Server{
		node:        maelstrom.NewNode(),
    mutex: new(sync.Mutex),
		incrementor: 0,
	}
	srv.node.Handle("generate", srv.idgen)

	return srv
}

type Server struct {
	node        *maelstrom.Node

  mutex *sync.Mutex
	incrementor int
}

func (s Server) Run() error {
	return s.node.Run()
}

func (s *Server) idgen(req maelstrom.Message) error {
	var body maelstrom.MessageBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return err
	}

  s.mutex.Lock()
	id := s.incrementor
	s.incrementor++
  s.mutex.Unlock()

	body.Type = "generate_ok"
	resp := IdGenBody{
		MessageBody: body,
		Id:          s.node.ID() + "-" + strconv.Itoa(id),
	}

	return s.node.Reply(req, resp)
}

type IdGenBody struct {
	maelstrom.MessageBody
	Id string `json:"id"`
}
