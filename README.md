# Eleutherios

Eleutherios — утилита для Keenetic/Entware, которая настраивает WireGuard split tunnel: домены `*.ru` идут через обычный интернет провайдера, а остальной трафик уходит через WireGuard. Проект управляет `dnsmasq`, `ipset`, `iptables`, policy routing и автозапуском на роутере.

## Документация

- [Routing](docs/routing.md) — как устроена маршрутизация через DNS, `ipset`, mark и routing table.
- [ELEUTHERIOS_EXCLUDED](docs/ELEUTHERIOS_EXCLUDED.md) — какие IP-сети нельзя отправлять в WireGuard.
- [Debug commands](docs/debug.md) — команды для диагностики текущего состояния.

## Зависимости

На роутере с Entware установите зависимости:

```sh
opkg update
opkg install dnsmasq-full ipset iptables
```

## Запуск

Сначала посмотрите доступные WireGuard-интерфейсы и bridge-сети:

```sh
eleutherios vpn ls
eleutherios net ls
```

Запустите обход, указав Keenetic-имя WireGuard и сеть:

```sh
eleutherios start --wg Wireguard0 --net br0
```

Проверить состояние:

```sh
eleutherios status
```

Остановить обход и вернуть настройки:

```sh
eleutherios stop
```

## Сборка

Собрать бинарники для поддерживаемых архитектур Keenetic:

```sh
./helpers/build.sh
```

Собрать `.ipk`-пакет, например для `mipsel`:

```sh
./helpers/package.sh mipsel
```
