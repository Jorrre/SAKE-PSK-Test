cd $HOME/Development/go
git checkout feature/plain-psk
cd $HOME/Development/SAKE-PSK-Test
GOROOT=$HOME/Development/go
GOPATH=$HOME/go

$GOROOT/bin/go run $PWD/server/server.go
