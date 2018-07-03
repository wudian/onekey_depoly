#sh
git pull


go build -o shell/chain/chain -gcflags "-N -l" ./chain

# go build -o shell/api/api -gcflags "-N -l" ./api/cmd

# go build -o shell/browser/browser -gcflags "-N -l" ./browser

# go build -o shell/midware/midware -gcflags "-N -l" ./midware/cmd

go build -o shell/midapi/midapi -gcflags "-N -l" ./midapi/cmd
