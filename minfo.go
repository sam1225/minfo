/*
Program Name    : minfo.go
Binary Name		: minfo
Build Command	: go build -o minfo
Program Purpose : Displays custom system information and stats including IPv4 information.
Description     : Minfo is a tool for Linux Systems that displays a very succint overview of the OS.
                  Instead of running multiple commands to know various aspects of the OS, a single
                  command can be used.
                  It is written in Golang and tested on most Linux variants (Red Hat, Debian, & SUSE).
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var minfo_version string = "minfo 1.0"

var master_map = make(map[string]string, 0)

func main() {

	proc_uptime()

	hostnamectl()

	proc_meminfo()

	lscpu()

	timedatectl()

	argument_parsing()

}

func argument_parsing() {
	help_text := `Minfo is a tool for displaying minimal system & network information.

Usage: minfo [OPTION]

  -h, --help          this message
  -V, --version       output version information
  -a, --all           displays all information
  -s, --sysinfo       displays system information
  -i, --ipinfo        displays ip information

`

	if len(os.Args[:]) == 1 {
		fmt.Printf("%s", help_text)
	}

	arg := os.Args[1]

	switch {
	case arg == "-h" || arg == "--help":
		fmt.Printf("%s", help_text)

	case arg == "-V" || arg == "--version":
		fmt.Printf("%s\n", minfo_version)

	case arg == "-a" || arg == "--all":
		fmt.Printf("\n")
		display_system_info()
		fmt.Printf("\n\n")
		display_ip_info()
		fmt.Printf("\n")

	case arg == "-s" || arg == "--sysinfo":
		fmt.Printf("\n")
		display_system_info()
		fmt.Printf("\n")

	case arg == "-i" || arg == "--ipinfo":
		fmt.Printf("\n")
		display_ip_info()
		fmt.Printf("\n")

	default:
		fmt.Println("Unrecognized option.")
		fmt.Printf("%s", help_text)
		os.Exit(1)
	}
}

func display_system_info() {
	fmt.Printf("System Info:\n")
	fmt.Printf("  %-40s: %s\n", "CPU(s)", master_map["CPU(s)"])
	fmt.Printf("  %-40s: %s\n", "CPU Model", master_map["Model name"])
	fmt.Printf("  %-40s: %s\n", "Architecture", master_map["Architecture"])
	fmt.Printf("  %-40s: %s\n", "Total Memory", master_map["MemTotal"])
	fmt.Printf("  %-40s: %s\n", "Uptime", master_map["Uptime"])
	fmt.Printf("  %-40s: %s\n", "Operating System", master_map["Operating System"])
	fmt.Printf("  %-40s: %s\n", "Kernel", master_map["Kernel"])
	fmt.Printf("  %-40s: %s\n", "System Type", master_map["Virtualization"])
	fmt.Printf("  %-40s: %s\n", "Hostname", master_map["Static hostname"])
	fmt.Printf("  %-40s: %s\n", "Time Zone", master_map["Time zone"])
}

func display_ip_info() {
	fmt.Printf("IPv4 Info:\n")
	fmt.Printf("  %-40s  %s\n", "interface", "ip")
	fmt.Printf("  %-40s  %s\n", "---------", "---------")
	for k, v := range IP_Details() {
		fmt.Printf("  %-40s: %s\n", k, v)
	}
}

func IP_Details() map[string]string {
	ip_map := make(map[string]string, 0)

	interfaces, err := net.Interfaces()
	checkError(err)

	for _, i := range interfaces {
		interface_name := i.Name
		if interface_name == "lo" {
			continue
		}

		byNameInterface, err := net.InterfaceByName(i.Name)
		checkError(err)

		addresses, err := byNameInterface.Addrs()
		for k, v := range addresses {
			ip, _, _ := net.ParseCIDR(v.String())
			if ip4 := ip.To4(); ip4 != nil {
				key_str1 := fmt.Sprintf("%v[%v]", interface_name, k)
				ip_map[key_str1] = v.String()
			}
		}
	}

	return ip_map
}

func proc_uptime() {
	cmd_uptime := exec.Command("cat", "/proc/uptime")
	out_proc_uptime, _ := cmd_uptime.CombinedOutput()
	uptime_str := string(out_proc_uptime)
	uptime_sec_str := strings.Split(uptime_str, " ")
	uptime_sec_float, _ := strconv.ParseFloat(uptime_sec_str[0], 64)

	master_map["Uptime"] = fmt.Sprintf("%v secs / %v mins / %v hours / %v days",
		uptime_sec_str[0],
		int(uptime_sec_float)/60,
		int(uptime_sec_float)/60/60,
		int(uptime_sec_float)/60/60/24)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func OS_Command_Execution(out []byte, output_seperator string) map[string]string {
	data_map := make(map[string]string, 0)
	scanner1 := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner1.Scan() {
		s := strings.Split(scanner1.Text(), output_seperator)
		s_0 := strings.TrimSpace(s[0])
		s_1 := strings.TrimSpace(s[1])
		data_map[s_0] = s_1
	}

	return data_map
}

func hostnamectl() {
	cmd_hostnamectl := exec.Command("hostnamectl")
	out_hostnamectl, _ := cmd_hostnamectl.CombinedOutput()
	map_hostnamectl := OS_Command_Execution(out_hostnamectl, ":")

	if value, ok := map_hostnamectl["Operating System"]; ok {
		master_map["Operating System"] = value
	} else {
		master_map["Operating System"] = "N/A"
	}

	if value, ok := map_hostnamectl["Kernel"]; ok {
		master_map["Kernel"] = value
	} else {
		master_map["Kernel"] = "N/A"
	}

	if value, ok := map_hostnamectl["Virtualization"]; ok {
		master_map["Virtualization"] = fmt.Sprintf("Virtual Machine (%v)", value)
	} else {
		master_map["Virtualization"] = "Physical Machine"
	}

	if value, ok := map_hostnamectl["Architecture"]; ok {
		master_map["Architecture"] = value
	} else {
		master_map["Architecture"] = "N/A"
	}

	if value, ok := map_hostnamectl["Static hostname"]; ok {
		master_map["Static hostname"] = value
	} else {
		master_map["Static hostname"] = "N/A"
	}

}

func lscpu() {
	cmd_lscpu := exec.Command("lscpu")
	out_lscpu, _ := cmd_lscpu.CombinedOutput()
	map_lscpu := OS_Command_Execution(out_lscpu, ":")

	if value, ok := map_lscpu["CPU(s)"]; ok {
		master_map["CPU(s)"] = value
	} else {
		master_map["CPU(s)"] = "N/A"
	}

	if value, ok := map_lscpu["Model name"]; ok {
		master_map["Model name"] = value
	} else {
		master_map["Model name"] = "N/A"
	}
}

func proc_meminfo() {
	cmd_proc_meminfo := exec.Command("cat", "/proc/meminfo")
	out_proc_meminfo, _ := cmd_proc_meminfo.CombinedOutput()
	map_proc_meminfo := OS_Command_Execution(out_proc_meminfo, ":")

	if value, ok := map_proc_meminfo["MemTotal"]; ok {
		mem_str := value
		re := regexp.MustCompile("[0-9]+")
		mem_int, _ := strconv.Atoi(re.FindAllString(mem_str, -1)[0])
		master_map["MemTotal"] = fmt.Sprintf("%v / %v MB / %v GB",
			mem_str,
			mem_int/1024,
			mem_int/1024/1024)
	} else {
		master_map["MemTotal"] = "N/A"
	}
}

func timedatectl() {
	cmd_timedatectl := exec.Command("timedatectl")
	out_timedatectl, _ := cmd_timedatectl.CombinedOutput()
	map_timedatectl := OS_Command_Execution(out_timedatectl, ":")

	if value, ok := map_timedatectl["Time zone"]; ok {
		master_map["Time zone"] = value
	} else {
		master_map["Time zone"] = "N/A"
	}
}
