cd $HOME/Development/go
git checkout feature/lp2
cd $HOME/Development/SAKE-PSK-Test
GOROOT=$HOME/Development/go
GOPATH=$HOME/go

$GOROOT/bin/go run client/client.go 127.0.0.1:2208 10 -r
