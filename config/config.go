package config

//ServerConfig manager server confog
type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

//EventBusConfig rabbitmq server confog
type EventBusConfig struct {
	AMQPHost           string `json:"amqpHost"`
	AMQPPort           int    `json:"amqPort"`
	ManagementHost     string `json:"managementHost"`
	ManagementPort     int    `json:"managementPort"`
	ManagementUserName string `json:"managementUserName"`
	ManagementPassword string `json:"managementPassword"`
}
