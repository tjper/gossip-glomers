package faulttolerantbroadcast

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/tjper/gossip-glomers/internal/safe"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var (
	errTopologyMissingNode  = errors.New("topology is missing node")
	errNodeNeighborsInvalid = errors.New("node neighbors invalid")
	errNodeNeighborInvalid  = errors.New("node neighbor invalid")
)

func NewServer() *Server {
	srv := &Server{
		node:      maelstrom.NewNode(),
		neighbors: safe.NewSlice[string](),
		messages:  safe.NewSet[int](),
	}
	srv.node.Handle("broadcast", srv.broadcast)
	srv.node.Handle("read", srv.read)
	srv.node.Handle("topology", srv.topology)

	return srv
}

type Server struct {
	node *maelstrom.Node

	neighbors *safe.Slice[string]
	messages  *safe.Set[int]
}

func (s Server) Run() error {
	return s.node.Run()
}

type broadcastRequestBody struct {
	maelstrom.MessageBody
	Message int `json:"message"`
}

func (s Server) broadcast(req maelstrom.Message) error {
	var body broadcastRequestBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return err
	}

	if s.messages.Add(body.Message) {
		if err := s.broadcastToNeighbors(body); err != nil {
			return err
		}
	}

	reply := body.MessageBody
	reply.Type = "broadcast_ok"
	if err := s.node.Reply(req, reply); err != nil {
		return err
	}

	return nil
}

func (s Server) broadcastToNeighbors(msg broadcastRequestBody) error {
	destinations := s.neighbors.Contents()
	for _, destination := range destinations {
		go func(destination string) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := s.broadcastToNeighbor(ctx, destination, msg); err != nil {
				log.Printf("while broadcasting to neighbor: %s", err)
			}
		}(destination)
	}
	return nil
}

func (s Server) broadcastToNeighbor(ctx context.Context, destination string, msg broadcastRequestBody) error {
	ackc := make(chan struct{})
	acknowledge := func(resp maelstrom.Message) error {
		close(ackc)
		return nil
	}
	broadcast := func() error {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		if err := s.node.RPC(destination, msg, acknowledge); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ackc:
			return nil
		}
	}

	for {
		err := broadcast()
		if errors.Is(err, context.DeadlineExceeded) {
			continue
		}
		if err != nil {
			return err
		}
		return nil
	}
}

type readResponseBody struct {
	maelstrom.MessageBody
	Messages []int `json:"messages"`
}

func (s Server) read(req maelstrom.Message) error {
	var body maelstrom.MessageBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return err
	}

	messages := s.messages.Contents()

	body.Type = "read_ok"
	resp := readResponseBody{
		MessageBody: body,
		Messages:    messages,
	}

	return s.node.Reply(req, resp)
}

type topologyRequestBody struct {
	maelstrom.MessageBody
	Topology map[string]any `json:"topology"`
}

func (s Server) topology(req maelstrom.Message) error {
	var body topologyRequestBody
	if err := json.Unmarshal(req.Body, &body); err != nil {
		return err
	}

	field, ok := body.Topology[s.node.ID()]
	if !ok {
		return errTopologyMissingNode
	}

	values, ok := field.([]any)
	if !ok {
		return errNodeNeighborsInvalid
	}

	var neighbors []string
	for _, value := range values {
		neighbor, ok := value.(string)
		if !ok {
			return errNodeNeighborInvalid
		}
		neighbors = append(neighbors, neighbor)
	}

	s.neighbors.Set(neighbors...)

	resp := body.MessageBody
	resp.Type = "topology_ok"
	return s.node.Reply(req, resp)
}
