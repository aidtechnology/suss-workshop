module github.com/aidtechnology/suss-workshop

go 1.12

require (
	github.com/bryk-io/x v0.0.0-20190614052234-0398d942366b
	github.com/chzyer/readline v0.0.0-20160729034951-a0c5244a21f4
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/google/certificate-transparency-go v1.0.21 // indirect
	github.com/gorilla/mux v1.7.2
	github.com/gorilla/websocket v1.4.0
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/kisielk/sqlstruct v0.0.0-20150923205031-648daed35d49 // indirect
	github.com/lib/pq v1.1.1 // indirect
	github.com/logrusorgru/aurora v0.0.0-20190428105938-cea283e61946
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
)

replace (
	github.com/cloudflare/cfssl => github.com/bryk-io/cfssl v0.0.0-20190614051308-96819d845a26
	github.com/dgraph-io/badger v1.5.5 => github.com/bryk-io/badger v1.5.5
	github.com/grpc-ecosystem/go-grpc-middleware => github.com/bryk-io/go-grpc-middleware v1.0.1-0.20190419153159-d28668ee9f4e
)
