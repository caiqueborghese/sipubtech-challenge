module github.com/caiqueborghese/sipubtech-challenge/movies

go 1.25

// usa o código local do módulo proto
replace github.com/caiqueborghese/sipubtech-challenge/proto => ../proto

require (
	github.com/caiqueborghese/sipubtech-challenge/proto v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.6.0
	github.com/nats-io/nats.go v1.45.0
	github.com/stretchr/testify v1.9.0
	go.mongodb.org/mongo-driver v1.17.4
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.7
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
