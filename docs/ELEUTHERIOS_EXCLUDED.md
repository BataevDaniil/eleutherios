• ELEUTHERIOS_EXCLUDED это список IP-сетей, которые нельзя отправлять в WireGuard. В ELEUTHERIOS_MARK он проверяется первым:

  -A ELEUTHERIOS_MARK -m set --match-set ELEUTHERIOS_EXCLUDED dst -j RETURN

  То есть если destination IP попал в этот список, трафик не маркируется и идёт обычным маршрутом.

  Сейчас туда добавляются:

  0.0.0.0/8          служебная “this network”
  10.0.0.0/8         private LAN
  100.64.0.0/10      CGNAT, часто WAN у провайдера
  127.0.0.0/8        localhost
  169.254.0.0/16     link-local
  172.16.0.0/12      private LAN
  192.168.0.0/16     private LAN
  224.0.0.0/4        multicast
  240.0.0.0/4        reserved
  78.47.125.180      служебный IP Keenetic/my.keenetic

  Почему это нужно:

  - чтобы не отправлять локальную сеть через VPN;
  - чтобы не ломать доступ к роутеру, DHCP/DNS/локальным устройствам;
  - чтобы не отправлять CGNAT/служебные адреса провайдера в туннель;
  - чтобы multicast/link-local не уходили в WireGuard;
  - чтобы служебные адреса Keenetic оставались доступны напрямую.

  Посмотреть список:

  ipset list ELEUTHERIOS_EXCLUDED
