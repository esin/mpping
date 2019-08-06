# mpping
mpping - is a MiningPool (ₜᵣᵤₑ)Ping utility

#### Пинг до пула
Чем ближе находится сервер пула до майнера, тем быстрее майнер будет получать информацию о смене блока и всегда считать именно новый блок

В основном майнинг-пулы рекомендуют проверять пинг с помощью стандартной утилиты ping, что несколько неправильно

Сейчас я нахожусь в Москве и пропингую 3 разных пула flypool: азиатский, европейский и американский:

Азиатский:
```
$ ping asia1-zcash.flypool.org -c5
PING cee82a731d6d4ab6ab3d61786aefa5ce.pacloudflare.com (172.65.255.87) 56(84) bytes of data.
64 bytes from 172.65.255.87: icmp_seq=1 ttl=58 time=25.4 ms
64 bytes from 172.65.255.87: icmp_seq=2 ttl=58 time=25.4 ms
64 bytes from 172.65.255.87: icmp_seq=3 ttl=58 time=25.4 ms
64 bytes from 172.65.255.87: icmp_seq=4 ttl=58 time=25.4 ms
64 bytes from 172.65.255.87: icmp_seq=5 ttl=58 time=25.3 ms

--- cee82a731d6d4ab6ab3d61786aefa5ce.pacloudflare.com ping statistics ---
5 packets transmitted, 5 received, 0% packet loss, time 4005ms
rtt min/avg/max/mdev = 25.372/25.423/25.462/0.105 ms
```

Европейский:
```
$ ping eu1-zcash.flypool.org -c5
PING 2e9f2d5e5c8f4a96905af2066d1a569a.pacloudflare.com (172.65.200.16) 56(84) bytes of data.
64 bytes from 172.65.200.16: icmp_seq=1 ttl=58 time=25.4 ms
64 bytes from 172.65.200.16: icmp_seq=2 ttl=58 time=25.3 ms
64 bytes from 172.65.200.16: icmp_seq=3 ttl=58 time=25.3 ms
64 bytes from 172.65.200.16: icmp_seq=4 ttl=58 time=25.3 ms
64 bytes from 172.65.200.16: icmp_seq=5 ttl=58 time=25.3 ms

--- 2e9f2d5e5c8f4a96905af2066d1a569a.pacloudflare.com ping statistics ---
5 packets transmitted, 5 received, 0% packet loss, time 4005ms
rtt min/avg/max/mdev = 25.319/25.348/25.405/0.146 ms
```

Американский:
```
$ ping us1-zcash.flypool.org -c5
PING f5f5671e1d4641c78920647d677fc127.pacloudflare.com (172.65.245.55) 56(84) bytes of data.
64 bytes from 172.65.245.55: icmp_seq=1 ttl=58 time=24.8 ms
64 bytes from 172.65.245.55: icmp_seq=2 ttl=58 time=24.7 ms
64 bytes from 172.65.245.55: icmp_seq=3 ttl=58 time=24.8 ms
64 bytes from 172.65.245.55: icmp_seq=4 ttl=58 time=24.8 ms
64 bytes from 172.65.245.55: icmp_seq=5 ttl=58 time=24.7 ms

--- f5f5671e1d4641c78920647d677fc127.pacloudflare.com ping statistics ---
5 packets transmitted, 5 received, 0% packet loss, time 4007ms
rtt min/avg/max/mdev = 24.731/24.802/24.849/0.203 ms
```

Видно, что все 3 пула находятся примерно в одинаковом расстоянии от меня, хотя географически это невозможно (Европейский самый близкий должен быть, но пинг до него меньше чем до Американского)

Такое получатеся, из-за того, что пул использует Cloudflare Spectrum. Это защита от DDoS и поэтому icmp пакет останавливается на граничных серверах Cloudflare

Разумеется, пакет от граничного сервера до самого сервера пула идёт по кратчайшему пути, но этот путь не покажет утилита ping

Для более качественной проверки связи до пула и выбора пула для майнига была написана утилита `mpping`