# Routing

Маршрутизация построена вокруг DNS-перехвата, `ipset`, маркировки трафика и отдельной routing table для WireGuard.

## Поток обработки

1. Клиент из выбранной сети, например `br0`, делает DNS-запрос.
2. `iptables nat PREROUTING` перехватывает DNS на порт `53`.
3. DNS-запрос отправляется в Entware `dnsmasq` на `127.0.0.1:9753`.
4. `dnsmasq` резолвит домен.
5. Если домен `*.ru`, `dnsmasq` добавляет полученные IP в `ipset` `ELEUTHERIOS_RU`.
6. Когда клиент идёт на сайт, `iptables mangle PREROUTING` проверяет destination IP.
7. Если IP есть в `ELEUTHERIOS_RU`, трафик не маркируется и идёт обычным маршрутом ISP.
8. Если IP нет в `ELEUTHERIOS_RU`, трафик получает mark `0xd1000`.
9. `ip rule` отправляет marked-трафик в таблицу `1001`.
10. Таблица `1001` ведёт default route через WireGuard `nwg0`.

## Схема

```text
             Клиент в br0
                 |
                 | DNS query :53
                 v
    iptables nat PREROUTING
    -i br0 -j ELEUTHERIOS_DNS
                 |
                 v
    DNAT -> 127.0.0.1:9753
                 |
                 v
              dnsmasq
      ipset=/.ru/ELEUTHERIOS_RU
                 |
      +----------+----------+
      |                     |
  домен *.ru            не *.ru
      |                     |
      v                     v
IP добавляются        IP не попадают
в ELEUTHERIOS_RU      в ELEUTHERIOS_RU
```

```text
             Клиентский трафик к IP
                       |
                       v
        iptables mangle PREROUTING
        -i br0 -j ELEUTHERIOS_MARK
                       |
      +----------------+----------------+
      |                                 |
dst IP в ELEUTHERIOS_RU      dst IP не в ELEUTHERIOS_RU
      |                                 |
      v                                 v
RETURN, без mark                 MARK 0xd1000
      |                                 |
      v                                 v
main routing table        ip rule fwmark 0xd1000
обычный интернет          lookup table 1001
ISP                                 |
                                    v
                              default dev nwg0
                              WireGuard
```

## Текущее состояние

| Компонент | Назначение |
| --- | --- |
| `ELEUTHERIOS_RU` | Список IP для `*.ru` |
| `ELEUTHERIOS_MARK` | Правило выбора маршрута |
| `ip rule` | Связь mark `0xd1000` -> table `1001` |
| `table 1001` | Default route через `nwg0` |

## Главная идея

Заранее неизвестно, какие адреса относятся к `не .ru`. Поэтому `.ru` помечается как исключение через DNS, а всё остальное автоматически считается трафиком через WireGuard.
