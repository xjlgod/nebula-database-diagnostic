# !/bin/bash
go build -o ./diag-cli
rm -rf ./out/
./diag-cli info --config config/debug.yaml --name node0 -D 6s -P 2s -O ./o
