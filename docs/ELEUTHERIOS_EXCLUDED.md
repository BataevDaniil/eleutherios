# ELEUTHERIOS_EXCLUDED

`ELEUTHERIOS_EXCLUDED` — список IP-сетей, которые нельзя отправлять в WireGuard.

В `ELEUTHERIOS_MARK` этот список проверяется первым:

```iptables
-A ELEUTHERIOS_MARK -m set --match-set ELEUTHERIOS_EXCLUDED dst -j RETURN
```

Если destination IP попал в этот список, трафик не маркируется и идёт обычным маршрутом.

## Текущий состав списка

| Сеть или адрес | Назначение |
| --- | --- |
| `0.0.0.0/8` | Служебная сеть `this network` |
| `10.0.0.0/8` | Private LAN |
| `100.64.0.0/10` | CGNAT, часто WAN у провайдера |
| `127.0.0.0/8` | Localhost |
| `169.254.0.0/16` | Link-local |
| `172.16.0.0/12` | Private LAN |
| `192.168.0.0/16` | Private LAN |
| `224.0.0.0/4` | Multicast |
| `240.0.0.0/4` | Reserved |
| `78.47.125.180` | Служебный IP Keenetic/my.keenetic |

## Зачем это нужно

- Не отправлять локальную сеть через VPN.
- Не ломать доступ к роутеру, DHCP, DNS и локальным устройствам.
- Не отправлять CGNAT и служебные адреса провайдера в туннель.
- Не отправлять multicast и link-local в WireGuard.
- Оставлять служебные адреса Keenetic доступными напрямую.

## Проверка

Посмотреть список:

```sh
ipset list ELEUTHERIOS_EXCLUDED
```
