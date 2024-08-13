# hack to get around no tun devices in containers
if ! [ -c /dev/net/tun ]; then
 mkdir -p /dev/net
 mknod -m 666 /dev/net/tun c 10 200
fi

# hack to keep container up forever while we attach our 
# own shell.
while true; do sleep 10; done;