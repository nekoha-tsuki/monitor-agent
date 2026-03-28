package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int    `yaml:"port" validate:"gte=1,lte=65535"`
		Host string `yaml:"host" validate:"ip"`
	} `yaml:"server"`
	Auth struct {
		Token string `yaml:"token"`
	} `yaml:"auth"` // Added missing tag
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Initialize with defaults
	config := &Config{}
	config.Server.Port = 8080
	config.Server.Host = "127.0.0.1"

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	// 1. Load the configuration
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	fmt.Printf("Starting server on %s:%d...\n", cfg.Server.Host, cfg.Server.Port)

	// 2. Monitoring Loop
	for {
		fmt.Println("--- System Stats ---")

		// CPU Usage
		cpuPercent, _ := cpu.Percent(time.Second, false)
		if len(cpuPercent) > 0 {
			fmt.Printf("CPU Usage:     %.2f%%\n", cpuPercent[0])
		}

		// Memory Usage
		vMem, _ := mem.VirtualMemory()
		fmt.Printf("Memory:        %.2f%% (Used: %vMB / Total: %vMB)\n",
			vMem.UsedPercent, vMem.Used/1024/1024, vMem.Total/1024/1024)

		// Network Usage
		netStats, _ := net.IOCounters(false)
		if len(netStats) > 0 {
			fmt.Printf("Net Sent:      %v KB\n", netStats[0].BytesSent/1024)
			fmt.Printf("Net Recv:      %v KB\n", netStats[0].BytesRecv/1024)
		}

		fmt.Println("--------------------\n")
		time.Sleep(2 * time.Second)
	}
}
