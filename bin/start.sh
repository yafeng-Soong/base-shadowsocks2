./shadowsocks2-linux -c 'ss://AEAD_CHACHA20_POLY1305:123456@127.0.0.1:8488' -verbose -socks :1080
./shadowsocks2-linux -s 'ss://AEAD_CHACHA20_POLY1305:123456@:8488' -verbose
./shadowsocks2-linux -c 'ss://AES-256-GCM:OshtAgryp@205.185.125.213:9983' -verbose -socks :1080