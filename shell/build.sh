
git pull

go build -o shell/chain/chain -gcflags "-N -l" ./chain

go build -o shell/midapi/midapi -gcflags "-N -l" ./midapi/cmd
