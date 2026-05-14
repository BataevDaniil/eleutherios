# Debug commands

## Основные команды

```sh
iptables-save
```

## Логи Eleutherios

Продублировать вывод команды в файл:

```sh
eleutherios start --wg Wireguard0 --net br0 --log-file /tmp/eleutherios.log
```

После такого запуска установленные NDM hook-скрипты также будут писать в этот файл.

Посмотреть ошибки в сохранённом логе:

```sh
grep 'level=ERROR' /tmp/eleutherios.log
```

## Правила Eleutherios

Показать только правила Eleutherios:

```sh
iptables-save | grep ELEUTHERIOS
```

По таблицам отдельно:

```sh
iptables -t nat -S
iptables -t mangle -S
iptables -t filter -S
```

## Счётчики пакетов

```sh
iptables -t nat -vnL
iptables -t mangle -vnL
```

## Маршруты и policy routing

```sh
ip route
ip route show table 1001
ip rule
```

## ipset-списки

```sh
ipset -n list
ipset list ELEUTHERIOS_RU
ipset list ELEUTHERIOS_EXCLUDED
```

## DNS-конфиги

```sh
cat /opt/etc/dnsmasq.conf
cat /opt/etc/dnsmasq.d/eleutherios.dnsmasq
```

## Полезный набор для диагностики

```sh
iptables-save | grep ELEUTHERIOS
ip rule | grep d1000
ip route show table 1001
ipset list ELEUTHERIOS_RU
ipset list ELEUTHERIOS_EXCLUDED
```
