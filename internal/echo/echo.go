package echo

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func NewServer() *Server {
	srv := &Server{
		node: maelstrom.NewNode(),
	}
	srv.node.Handle("echo", srv.echo)

	return srv
}

type Server struct {
	node *maelstrom.Node
}

func (s Server) Run() error {
	return s.node.Run()
}

func (s Server) echo(req maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return err
	}

	body["type"] = "echo_ok"

	return s.node.Reply(req, body)
}
