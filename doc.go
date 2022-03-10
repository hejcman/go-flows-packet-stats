/*
Package go_flows_packet_stats implements some features which should be exported from the gathered flows. These features
are a 1:1 implementation of the PSTATS and PHISTS from the ipfixprobe exporter by CESNET.

PStats: https://github.com/CESNET/ipfixprobe#pstats

PHists: https://github.com/CESNET/ipfixprobe#phists

The exported features use the standard defined in libfds for use with CESNET collectors.

LibFDS: https://github.com/CESNET/libfds

CESNET feature definitions: https://github.com/CESNET/libfds/blob/master/config/system/elements/cesnet.xml
*/
package go_flows_packet_stats
