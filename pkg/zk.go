package pkg

import (
	"time"

	"github.com/go-zookeeper/zk"
)

// connectZk sets up the zookeeper connection
func ConnectZk(zkAddr string) (*zk.Conn, error) {
	conn, _, err := zk.Connect([]string{zkAddr}, 1*time.Second)
	return conn, err
}
