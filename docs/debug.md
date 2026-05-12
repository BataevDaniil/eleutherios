# Debug commands

## Основные команды

```sh
iptables-save
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
