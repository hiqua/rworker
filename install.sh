#!/bin/sh
set -e
set -u

# TODO: make sure we're running from repo folder
# TODO: double-check that reinstallation is handled properly

CLIENT_BIN="client"
SERVER_BIN="server"
TARGET_BIN="/usr/local/bin/"


if [ ! -e "$CLIENT_BIN" ] || [ ! -e "$SERVER_BIN" ]; then
  echo "Binaries not available, please build the binaries before:"
  echo "go build cmd/client/client.go"
  echo "go build cmd/server/server.go"
  echo "Exiting..."
  exit 1
fi

sudo rm -f "$TARGET_BIN"/"$CLIENT_BIN"
sudo rm -f "$TARGET_BIN"/"$SERVER_BIN"
sudo cp "$CLIENT_BIN" "$SERVER_BIN" "$TARGET_BIN"

systemctl stop --user rworker.service || true # ignore non-existent service
sudo cp etc/linux-systemd/user/rworker.service /usr/lib/systemd/user/rworker.service
systemctl enable --user rworker.service
systemctl start --user rworker.service
