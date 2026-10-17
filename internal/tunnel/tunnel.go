package tunnel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Dir is the directory of tunnel .conf files (default ~/.config/godeb-wg/tunnels).
func Dir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "godeb-wg", "tunnels")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}

func pathFor(name string) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name+".conf"), nil
}

// List returns tunnel names (without .conf).
func List() ([]string, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".conf") {
			continue
		}
		names = append(names, strings.TrimSuffix(e.Name(), ".conf"))
	}
	return names, nil
}

// Import copies confText into the tunnel store as name.conf.
func Import(name, confText string) error {
	if name == "" || strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("invalid tunnel name")
	}
	p, err := pathFor(name)
	if err != nil {
		return err
	}
	return os.WriteFile(p, []byte(confText), 0o600)
}

// Read returns the conf file contents.
func Read(name string) (string, error) {
	p, err := pathFor(name)
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Up brings the tunnel up via wg-quick (requires root / polkit / CAP_NET_ADMIN).
func Up(name string) error {
	p, err := pathFor(name)
	if err != nil {
		return err
	}
	cmd := exec.Command("wg-quick", "up", p)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Down brings the tunnel down.
func Down(name string) error {
	p, err := pathFor(name)
	if err != nil {
		return err
	}
	cmd := exec.Command("wg-quick", "down", p)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsUp reports whether interface name is present in `wg show`.
func IsUp(name string) bool {
	out, err := exec.Command("wg", "show", "interfaces").Output()
	if err != nil {
		return false
	}
	for _, iface := range strings.Fields(string(out)) {
		if iface == name {
			return true
		}
	}
	// wg-quick uses the conf basename as interface name
	return false
}
