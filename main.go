package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func showHelp() {
	fmt.Println("Auto SSH Config - Setup SSH key dan config otomatis")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  auto-sshconfig           Setup SSH baru (interactive)")
	fmt.Println("  auto-sshconfig -h        Help")
	fmt.Println("  auto-sshconfig -r        Remove SSH config")
	fmt.Println("  auto-sshconfig -c        Copy key ke server")
	fmt.Println("  auto-sshconfig -s        Show SSH configs")
	fmt.Println("")
	fmt.Println("Contoh:")
	fmt.Println("  auto-sshconfig")
	fmt.Println("  auto-sshconfig -r vpstest")
	fmt.Println("  auto-sshconfig -c vpstest root@192.168.1.100")
	fmt.Println("  auto-sshconfig -s")
}

func showSSHConfigs() {
	homeDir, _ := os.UserHomeDir()
	sshDir := filepath.Join(homeDir, ".ssh")
	configPath := filepath.Join(sshDir, "config")

	fmt.Println("=== SSH Configs ===")

	if data, err := os.ReadFile(configPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "Host ") {
				hostName := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "Host "))
				if hostName != "*" && hostName != "" {
					keyPath := filepath.Join(sshDir, fmt.Sprintf("id_rsa_%s", hostName))
					if _, err := os.Stat(keyPath); err == nil {
						fmt.Printf("✓ %s\n", hostName)
					} else {
						fmt.Printf("- %s (key tidak ada)\n", hostName)
					}
				}
			}
		}
	} else {
		fmt.Println("SSH config tidak ditemukan")
	}
}

func deleteSSHConfig() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: auto-sshconfig -r <hostname>")
		fmt.Println("Contoh: auto-sshconfig -r vpstest")
		return
	}

	hostName := os.Args[2]
	homeDir, _ := os.UserHomeDir()
	sshDir := filepath.Join(homeDir, ".ssh")
	keyPath := filepath.Join(sshDir, fmt.Sprintf("id_rsa_%s", hostName))
	pubKeyPath := keyPath + ".pub"
	configPath := filepath.Join(sshDir, "config")

	fmt.Printf("Menghapus SSH config untuk: %s\n", hostName)

	if _, err := os.Stat(keyPath); err == nil {
		os.Remove(keyPath)
		fmt.Printf("✓ Hapus: %s\n", keyPath)
	}
	if _, err := os.Stat(pubKeyPath); err == nil {
		os.Remove(pubKeyPath)
		fmt.Printf("✓ Hapus: %s\n", pubKeyPath)
	}

	if data, err := os.ReadFile(configPath); err == nil {
		lines := strings.Split(string(data), "\n")
		var filteredLines []string
		skipBlock := false

		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "Host ") {
				hostInLine := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "Host "))
				skipBlock = (hostInLine == hostName)
			}
			if !skipBlock {
				filteredLines = append(filteredLines, line)
			}
		}

		newConfig := strings.Join(filteredLines, "\n")
		os.WriteFile(configPath, []byte(newConfig), 0600)
		fmt.Printf("✓ Hapus dari config: %s\n", configPath)
	}

	fmt.Printf("SSH config '%s' berhasil dihapus\n", hostName)
}

func copyKeyToServer() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: auto-sshconfig -c <hostname> <user@ip>")
		fmt.Println("Contoh: auto-sshconfig -c vpstest root@192.168.1.100")
		return
	}

	hostName := os.Args[2]
	userHost := os.Args[3]
	parts := strings.Split(userHost, "@")
	if len(parts) != 2 {
		fmt.Println("Format salah. Gunakan: user@ip")
		fmt.Println("Contoh: root@192.168.1.100")
		return
	}

	username := parts[0]
	ipAddress := parts[1]

	fmt.Printf("Copy key '%s' ke %s@%s\n", hostName, username, ipAddress)

	homeDir, _ := os.UserHomeDir()
	sshDir := filepath.Join(homeDir, ".ssh")
	keyPath := filepath.Join(sshDir, fmt.Sprintf("id_rsa_%s", hostName))
	pubKeyPath := keyPath + ".pub"

	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		fmt.Printf("Public key tidak ditemukan: %s\n", pubKeyPath)
		fmt.Println("Jalankan setup dulu: auto-sshconfig")
		return
	}

	fmt.Printf("Copy public key ke %s@%s...\n", username, ipAddress)
	cmd := exec.Command("ssh-copy-id", "-i", pubKeyPath, "-o", "StrictHostKeyChecking=no", userHost)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\nSSH-copy-id gagal. Copy manual:\n")
		pubKeyData, _ := os.ReadFile(pubKeyPath)
		fmt.Printf("%s\n", string(pubKeyData))
		fmt.Printf("Jalankan di server: echo '%s' >> ~/.ssh/authorized_keys\n", strings.TrimSpace(string(pubKeyData)))
	} else {
		fmt.Println("✓ Public key berhasil di-copy")
	}
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-h", "--help", "help":
			showHelp()
			return
		case "-r", "delete":
			deleteSSHConfig()
			return
		case "-c", "copy":
			copyKeyToServer()
			return
		case "-s", "show":
			showSSHConfigs()
			return
		default:
			fmt.Printf("Command tidak dikenal: %s\n", os.Args[1])
			fmt.Println("Gunakan -h untuk help")
			return
		}
	}

	fmt.Println("=== Auto SSH Config ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Nama host: ")
	hostName, _ := reader.ReadString('\n')
	hostName = strings.TrimSpace(hostName)

	fmt.Print("IP address: ")
	ipAddress, _ := reader.ReadString('\n')
	ipAddress = strings.TrimSpace(ipAddress)

	fmt.Print("Username (default: root): ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if username == "" {
		username = "root"
	}

	keyName := fmt.Sprintf("id_rsa_%s", hostName)
	homeDir, _ := os.UserHomeDir()
	sshDir := filepath.Join(homeDir, ".ssh")
	keyPath := filepath.Join(sshDir, keyName)
	pubKeyPath := keyPath + ".pub"

	fmt.Println("\n=== Setup SSH ===")

	os.MkdirAll(sshDir, 0700)

	fmt.Println("1. Cleaning known_hosts...")
	exec.Command("ssh-keygen", "-R", ipAddress).Run()

	fmt.Println("2. Generate SSH key...")
	os.Remove(keyPath)
	os.Remove(pubKeyPath)

	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", keyPath, "-N", "", "-C", fmt.Sprintf("%s@%s", username, ipAddress))
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
	os.Chmod(keyPath, 0600)
	fmt.Printf("✓ Key: %s\n", keyPath)

	fmt.Println("3. Copy public key...")
	cmd = exec.Command("ssh-copy-id", "-i", pubKeyPath, "-o", "StrictHostKeyChecking=no", fmt.Sprintf("%s@%s", username, ipAddress))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\nManual copy:\n")
		pubKeyData, _ := os.ReadFile(pubKeyPath)
		fmt.Printf("%s\n", string(pubKeyData))
		fmt.Printf("Jalankan di server: echo '%s' >> ~/.ssh/authorized_keys\n", strings.TrimSpace(string(pubKeyData)))
		fmt.Print("Enter jika sudah...")
		reader.ReadString('\n')
	}

	fmt.Println("4. Update SSH config...")
	configPath := filepath.Join(sshDir, "config")

	newConfig := fmt.Sprintf(`
Host %s
    HostName %s
    User %s
    IdentityFile %s
    StrictHostKeyChecking no
`, hostName, ipAddress, username, keyPath)

	var existingConfig string
	if data, err := os.ReadFile(configPath); err == nil {
		lines := strings.Split(string(data), "\n")
		var filteredLines []string
		skipBlock := false

		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "Host ") {
				hostInLine := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "Host "))
				skipBlock = (hostInLine == hostName)
			}
			if !skipBlock {
				filteredLines = append(filteredLines, line)
			}
		}
		existingConfig = strings.Join(filteredLines, "\n")
	}

	finalConfig := existingConfig + newConfig
	os.WriteFile(configPath, []byte(finalConfig), 0600)
	fmt.Printf("✓ Config: %s\n", configPath)

	fmt.Println("5. Test SSH...")
	cmd = exec.Command("ssh", "-o", "ConnectTimeout=5", hostName, "echo 'OK'")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Test gagal: %v\n", err)
	} else {
		fmt.Printf("✓ Test: %s", string(output))
	}

	fmt.Printf("\n=== Selesai ===\n")
	fmt.Printf("SSH: ssh %s\n", hostName)
}
