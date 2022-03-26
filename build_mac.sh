#!/bin/zsh

make build-darwin-server

cp ./bin/server/Spellbook-Server-darwin test_server/
cd test_server
./Spellbook-Server-darwin start