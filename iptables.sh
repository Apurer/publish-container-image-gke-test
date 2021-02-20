#!/bin/bash
iptables -A INPUT -p tcp -m tcp -m multiport ! --dports 8443 -m owner --uid-owner hellouser -j DROP
ip6tables -A INPUT -p tcp -m tcp -m multiport ! --dports 8443 -m owner --uid-owner hellouser -j DROP


