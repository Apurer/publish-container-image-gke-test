#!/bin/bash
iptables -A INPUT -p tcp -m tcp -m multiport ! --dports 8443 -m owner --uid-owner 5777 -j DROP
ip6tables -A INPUT -p tcp -m tcp -m multiport ! --dports 8443 -m owner --uid-owner 5777 -j DROP


