package amqp

// TODO
// Inject custom logger

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/buroz/grpc-clean-example/pkg/config"
	"github.com/streadway/amqp"
)

type AmqpClient struct {
	sync.Mutex

	config *config.AmqpConfig

	UseTLS bool

	// TODO ...
	// SSLCA              string `toml:"ssl_ca"`
	// SSLCert            string `toml:"ssl_cert"`
	// SSLKey             string `toml:"ssl_key"`
	// InsecureSkipVerify bool

	Channel *amqp.Channel
}

func NewAmqpClient(conf *config.AmqpConfig) AmqpClient {
	return AmqpClient{
		config: conf,
	}
}

func (q *AmqpClient) Connect() error {
	q.Lock()
	defer q.Unlock()

	var connection *amqp.Connection
	var err error

	if q.UseTLS {
		// TODO !!
		// url := fmt.Sprintf("amqps://%s:%s@%s:%d/", q.config.User, q.config.Password, q.config.Host, q.config.Port)
		// connection, err = amqp.DialTLS(url, tls)
		return nil
	} else {
		url := fmt.Sprintf("amqp://%s:%s@%s:%d/", q.config.User, q.config.Password, q.config.Host, q.config.Port)
		connection, err = amqp.Dial(url)
	}

	if err != nil {
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	q.Channel = channel

	go func() {
		log.Printf("Closing: %s", <-connection.NotifyClose(make(chan *amqp.Error)))
		log.Printf("Trying to reconnect...")

		for err := q.Connect(); err != nil; err = q.Connect() {
			log.Println(err)
			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}
