package rabbitmq

import (
	"github.com/qazaqpyn/bookCRUD/model"
	audit "github.com/qazaqpyn/crud-audit-log/pkg/domain"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewServer(uri string) (*Server, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Server{
		conn: conn,
		ch:   ch,
	}, nil
}

func (s *Server) closeConnection() error {
	return s.conn.Close()
}

func (s *Server) closeChannel() error {
	return s.ch.Close()
}

func (s *Server) CloseServer() error {
	err := s.closeChannel()
	if err != nil {
		return err
	}

	err = s.closeConnection()
	return err
}

func (s *Server) SendToQueue(req model.Msg) error {
	action, err := audit.ToPbAction(req.Action)
	if err != nil {
		return err
	}

	entity, err := audit.ToPbEntity(req.Entity)
	if err != nil {
		return err
	}
	data, err := proto.Marshal(&audit.LogRequest{
		Action:    action,
		Entity:    entity,
		EntityId:  req.EntityID,
		Timestamp: timestamppb.New(req.Timestamp),
	})
	if err != nil {
		return err
	}

	q, err := s.ch.QueueDeclare(
		"messageQueue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	if err := s.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		}); err != nil {
		return err
	}
	return nil
}
