go build -o cmdry .
go build -o cmdry-ports ./plugins/ports/cmd/cmdry-ports
sudo install -Dm755 cmdry /usr/local/bin/cmdry
sudo install -Dm755 cmdry-ports /opt/cmdry/plugins/cmdry-ports
CMDRY_PLUGIN_DIR="$PWD" CMDRY_DATA_DIR="$PWD/.cmdry-data" ./cmdry serve