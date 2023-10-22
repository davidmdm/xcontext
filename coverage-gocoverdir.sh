rm -rf ./coverage && mkdir ./coverage
go test -count 1 -cover -test.gocoverdir=./coverage ./...
go tool covdata percent -i ./coverage
go tool covdata textfmt -i ./coverage -o ./coverage.out