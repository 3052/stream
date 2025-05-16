# nord

~~~
curl `
--proxy https://USERNAME:PASSWORD@pl128.nordvpn.com:89 `
--proxy-ca-native `
ip.smartproxy.com/json
~~~

- https://gist.github.com/cupcakearmy/7c30d5791280a27fadf04ffab1733ac7
- https://github.com/wenerme/wener/blob/master/notes/platform/nordvpn.md

## technology

~~~
{1 ikev2}
{3 openvpn_udp}
{5 openvpn_tcp}
{7 socks}
{15 openvpn_xor_udp}
{17 openvpn_xor_tcp}
{21 proxy_ssl}
{23 proxy_ssl_cybersec}
{35 wireguard_udp}
{42 openvpn_dedicated_udp}
{45 openvpn_dedicated_tcp}
{51 nordwhisper}
~~~

## NordVPN proxy setup for qBittorrent

https://support.nordvpn.com/hc/en-us/articles/20195967385745
