# Eleutherios

Eleutherios — утилита для Keenetic/Entware, которая настраивает WireGuard split tunnel: домены `*.ru` идут через обычный интернет провайдера, а остальной трафик уходит через WireGuard. Проект управляет `dnsmasq`, `ipset`, `iptables`, policy routing и автозапуском на роутере.

## Документация

- [Routing](docs/routing.md) — как устроена маршрутизация через DNS, `ipset`, mark и routing table.
- [ELEUTHERIOS_EXCLUDED](docs/ELEUTHERIOS_EXCLUDED.md) — какие IP-сети нельзя отправлять в WireGuard.
- [Debug commands](docs/debug.md) — команды для диагностики текущего состояния.

## Требования

- Роутер Keenetic с установленным Entware (раздел `/opt`).
- `start`, `stop` и `status` запускаются от **root** (под обычным пользователем выйдут с ошибкой).
- `dnsmasq-full` (не базовый `dnsmasq` — нужна поддержка `ipset`).
- `ipset` и `iptables`.
- HTTP-доступ к локальному RCI Keenetic (`http://127.0.0.1:79`) — должен быть включён в админке.

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

## Что делает `start` / `stop`

`start` выполняет 6 шагов:

1. Поднимает WireGuard через RCI Keenetic (если он был выключен).
2. Создаёт ipset `ELEUTHERIOS_RU` (для `*.ru`) и `ELEUTHERIOS_EXCLUDED` (приватные сети, multicast и т.п.).
3. Пишет конфиг dnsmasq в `/opt/etc/dnsmasq.d/eleutherios.dnsmasq` (домены `*.ru` попадают в ipset), бэкапит и подменяет `/opt/etc/dnsmasq.conf`, перезапускает dnsmasq на порту 9753.
4. Настраивает iptables: цепочка `ELEUTHERIOS_DNS` редиректит DNS на dnsmasq, `ELEUTHERIOS_MARK` ставит fwmark на пакеты для туннеля.
5. Добавляет маршрут по умолчанию в отдельной таблице (`RouteTableID=1001`) через WG, и `ip rule` по fwmark.
6. Ставит автозапуск: `/opt/etc/init.d/S96eleutherios` и NDM hook `/opt/etc/ndm/fs.d/15-eleutherios-start.sh`, плюс netfilter hook `/opt/etc/ndm/netfilter.d/100-eleutherios` для восстановления правил после ребута/реконфига.

`stop` делает обратное: чистит iptables, ipset, маршруты, init-скрипт и hook'и.

## Откат вручную

Если `stop` не отработал, можно убрать всё руками:

```sh
rm -f /opt/etc/init.d/S96eleutherios
rm -f /opt/etc/ndm/fs.d/15-eleutherios-start.sh
rm -f /opt/etc/ndm/netfilter.d/100-eleutherios
rm -f /opt/etc/dnsmasq.d/eleutherios.dnsmasq
mv /opt/etc/dnsmasq.conf.backup /opt/etc/dnsmasq.conf 2>/dev/null
ipset destroy ELEUTHERIOS_RU
ipset destroy ELEUTHERIOS_EXCLUDED
iptables -t nat -F ELEUTHERIOS_DNS && iptables -t nat -X ELEUTHERIOS_DNS
iptables -t mangle -F ELEUTHERIOS_MARK && iptables -t mangle -X ELEUTHERIOS_MARK
ip route flush table 1001
ip rule del fwmark 0xd1000/0xd1000 table 1001 priority 1778
/opt/etc/init.d/S56dnsmasq restart
```

## Логи

Все команды пишут логи в `stdout`. Чтобы дополнительно сохранить их в файл, используйте глобальный флаг `--log-file`:

```sh
eleutherios start --wg Wireguard0 --net br0 --log-file /tmp/eleutherios.log
```

Формат логов — `slog.TextHandler`, близкий к `logfmt`: `time=... level=... msg=... component=...`. Ошибки логируются с `level=ERROR`.

Если `start` был запущен с `--log-file`, установленные NDM hook-скрипты тоже будут писать в этот файл.

## Таймаут

По умолчанию команда падает по таймауту через 3 минуты. На медленном роутере (или для длинных диагностик через `status`) можно изменить:

```sh
eleutherios --timeout 10m start --wg Wireguard0 --net br0
eleutherios --timeout 30s status
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
