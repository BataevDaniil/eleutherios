• Основные команды:

iptables-save

Показать только правила eleutherios:

iptables-save | grep ELEUTHERIOS

По таблицам отдельно:

iptables -t nat -S
iptables -t mangle -S
iptables -t filter -S

Счётчики пакетов:

iptables -t nat -vnL
iptables -t mangle -vnL

Маршруты и policy routing:

ip route
ip route show table 1001
ip rule

ipset-списки:

ipset -n list
ipset list ELEUTHERIOS_RU
ipset list ELEUTHERIOS_EXCLUDED

DNS-конфиги:

cat /opt/etc/dnsmasq.conf
cat /opt/etc/dnsmasq.d/eleutherios.dnsmasq

Самый полезный набор для твоего текущего случая:

iptables-save | grep ELEUTHERIOS
ip rule | grep d1000
ip route show table 1001
ipset list ELEUTHERIOS_RU
ipset list ELEUTHERIOS_EXCLUDED