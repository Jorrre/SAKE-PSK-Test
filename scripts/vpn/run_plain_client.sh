cd $HOME/Development/go
git checkout feature/plain-psk
cd $HOME/Development/SAKE-PSK-Test
GOROOT=$HOME/Development/go
GOPATH=$HOME/go

$GOROOT/bin/go run client/client.go 10.212.136.193:2208 60 -r
