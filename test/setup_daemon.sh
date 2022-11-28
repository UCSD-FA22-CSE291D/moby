hack/make.sh binary
make install
echo "{\"experimental\":true}" > /etc/docker/daemon.json
dockerd &
