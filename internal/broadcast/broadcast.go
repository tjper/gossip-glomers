package broadcast

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/tjper/gossip-glomers/internal/safe"
)

func NewServer() *Server {
	srv := &Server{
		node:     maelstrom.NewNode(),
    messages: safe.NewSet[int](),
	}
	srv.node.Handle("broadcast", srv.broadcast)
	srv.node.Handle("read", srv.read)
	srv.node.Handle("topology", srv.topology)

	return srv
}

type Server struct {
	node *maelstrom.Node
  messages *safe.Set[int]
}

func (s Server) Run() error {
	return s.node.Run()
}

type broadcastRequestBody struct {
	maelstrom.MessageBody
	Message int `json:"message"`
}

func (s *Server) broadcast(req maelstrom.Message) error {
	var body broadcastRequestBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return err
	}

  _ = s.messages.Add(body.Message)

	body.MessageBody.Type = "broadcast_ok"
	return s.node.Reply(req, body.MessageBody)
}

type readRequestBody struct {
	maelstrom.MessageBody
	Messages []int `json:"messages"`
}

func (s *Server) read(req maelstrom.Message) error {
	var body maelstrom.MessageBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return err
	}

  messages := s.messages.Contents()

	body.Type = "read_ok"
	resp := readRequestBody{
		MessageBody: body,
		Messages:    messages,
	}

	return s.node.Reply(req, resp)
}

func (s Server) topology(req maelstrom.Message) error {
	var body maelstrom.MessageBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return err
	}

	body.Type = "topology_ok"
	return s.node.Reply(req, body)
}
