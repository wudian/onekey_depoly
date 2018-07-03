.PHONY: ann api all
all: ann api

ann:
	go build -ldflags "-X gitlab.zhonganinfo.com/tech_bighealth/za-delos/chain/version.commitVer=`git rev-parse HEAD`" -o ./build/ann ./chain
api:
	go build -o ./build/ann-api ./api/cmd

test:
	 # todo go test all case
	go test ./ann/mempool
	go test ./ann/consensus
	go test ./ann/blockchain
