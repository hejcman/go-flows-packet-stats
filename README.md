# go-flows-packet-stats

Package go_flows_packet_stats implements some features which should be exported from the gathered flows. These features
are a 1:1 implementation of the [PSTATS](https://github.com/CESNET/ipfixprobe#pstats) and 
[PHISTS](https://github.com/CESNET/ipfixprobe#phists) from the ipfixprobe exporter by CESNET.
The exported features use the standard defined in [libfds](https://github.com/CESNET/libfds) for use with CESNET collectors.
