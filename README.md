----
### Minfo (Minimal info)

Minfo is a tool for Linux Systems that displays a very succint overview of the OS.
Instead of running multiple commands to know various aspects of the OS, a single 
command can be used.

***It is written in Golang and tested on most Linux variants (Red Hat, Debian, & SUSE).***

[DOWNLOAD](https://github.com/sam1225/minfo/raw/master/minfo)

[Repository Link](https://github.com/sam1225/minfo)

```markdown
$ ./minfo
Minfo is a tool for displaying minimal system & network information.

Usage: minfo [OPTION]

  -h, --help          this message
  -V, --version       output version information
  -a, --all           displays all information
  -s, --sysinfo       displays system information
  -i, --ipinfo        displays ip information

$ ./minfo -a
System Info:
  CPU(s)                                  : 1
  CPU Model                               : Intel(R) Core(TM) i5 CPU M 480  @ 2.67GHz
  Architecture                            : x86-64
  Total Memory                            : 997956 kB / 974 MB / 0 GB
  Uptime                                  : 141.94 secs / 2 mins / 0 hours / 0 days
  Operating System                        : CentOS Linux 7 (Core)
  Kernel                                  : Linux 3.10.0-862.el7.x86_64
  System Type                             : Virtual Machine (vmware)
  Hostname                                : cos7
  Time Zone                               : America/Chicago (CDT, -0500)

IPv4 Info:
  interface                                 ip
  ---------                                 ---------
  ens33[0]                                : 192.168.229.20/24
  docker0[0]                              : 172.17.0.1/16

```
