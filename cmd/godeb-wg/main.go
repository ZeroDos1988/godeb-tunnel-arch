package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ZeroDos1988/godeb-tunnel-arch/internal/cidr"
	"github.com/ZeroDos1988/godeb-tunnel-arch/internal/trusted"
	"github.com/ZeroDos1988/godeb-tunnel-arch/internal/tunnel"
)

const version = "1.0.0-godeb"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "version":
		fmt.Println(version)
	case "list":
		err = cmdList()
	case "import":
		err = cmdImport(os.Args[2:])
	case "up":
		err = requireName(os.Args[2:], tunnel.Up)
	case "down":
		err = requireName(os.Args[2:], tunnel.Down)
	case "exclude":
		err = cmdExclude(os.Args[2:])
	case "monitor":
		err = cmdMonitor()
	case "show":
		err = cmdShow(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `godeb-wg %s — GoDeb Tunnel client for Arch / Linux

Usage:
  godeb-wg list
  godeb-wg import <name> <file.conf>
  godeb-wg show <name>
  godeb-wg up <name>
  godeb-wg down <name>
  godeb-wg exclude <allowed> <excluded>   # print computed AllowedIPs
  godeb-wg monitor                        # trusted-endpoint auto-disconnect
  godeb-wg version

Tunnel configs live in ~/.config/godeb-wg/tunnels/.
Trusted endpoint: add  # GoDeb-Trusted-Endpoint: http://host:8080/api/trusted
Excluded IPs:      add  # GoDeb-ExcludedIPs: 172.16.1.0/24  (and computed AllowedIPs)
`, version)
}

func cmdList() error {
	names, err := tunnel.List()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Println("(no tunnels)")
		return nil
	}
	for _, n := range names {
		state := "down"
		if tunnel.IsUp(n) {
			state = "up"
		}
		fmt.Printf("%-24s %s\n", n, state)
	}
	return nil
}

func cmdImport(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: godeb-wg import <name> <file.conf>")
	}
	data, err := os.ReadFile(args[1])
	if err != nil {
		return err
	}
	if err := tunnel.Import(args[0], string(data)); err != nil {
		return err
	}
	fmt.Printf("imported %s\n", args[0])
	return nil
}

func cmdShow(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: godeb-wg show <name>")
	}
	text, err := tunnel.Read(args[0])
	if err != nil {
		return err
	}
	fmt.Print(text)
	if !strings.HasSuffix(text, "\n") {
		fmt.Println()
	}
	return nil
}

func requireName(args []string, fn func(string) error) error {
	if len(args) != 1 {
		return fmt.Errorf("tunnel name required")
	}
	return fn(args[0])
}

func cmdExclude(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: godeb-wg exclude \"0.0.0.0/0, ::/0\" \"172.16.1.0/24\"")
	}
	allowed, excluded, err := cidr.ApplyExcludedIPs(args[0], args[1])
	if err != nil {
		return err
	}
	fmt.Println("# GoDeb-AllowedIPs-Base:", args[0])
	if len(excluded) > 0 {
		fmt.Println("# GoDeb-ExcludedIPs:", strings.Join(excluded, ", "))
	}
	fmt.Println("AllowedIPs =", strings.Join(allowed, ", "))
	return nil
}

func cmdMonitor() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	(&trusted.Monitor{}).Run(ctx)
	return nil
}
