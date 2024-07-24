package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Lease represents a single DHCP lease record
type Lease struct {
	IP                    string `json:"ip"`
	BindingState          string `json:"state"`
	HardwareEthernet      string `json:"hardware"`
	UID                   string `json:"uid"`
	VendorClassIdentifier string `json:"vendor"`
	Hostname              string `json:"host"`
}

// convertOctalStringToHex converts a string with octal escape sequences to a hex string
func convertOctalStringToHex(octalStr string) string {
	var hexStr strings.Builder
	octalStr = strings.Trim(octalStr, "\"")
	for i := 0; i < len(octalStr); i++ {
		if octalStr[i] == '\\' && i+3 < len(octalStr) && isOctalDigit(octalStr[i+1]) && isOctalDigit(octalStr[i+2]) && isOctalDigit(octalStr[i+3]) {
			octalVal, _ := strconv.ParseInt(octalStr[i+1:i+4], 8, 64)
			hexStr.WriteString(fmt.Sprintf("%02x ", octalVal))
			i += 3
		} else {
			hexStr.WriteString(fmt.Sprintf("%02x ", octalStr[i]))
		}
	}
	return hexStr.String()
}

// isOctalDigit checks if a byte is an octal digit
func isOctalDigit(c byte) bool {
	return c >= '0' && c <= '7'
}

// parseLeases parses the content of a DHCP leases file and returns a slice of Lease structs
func parseLeases(content string) []Lease {
	leases := []Lease{}
	leaseRegex := regexp.MustCompile(`lease\s+([\d\.]+)\s+{\n([^}]+)\n}`)
	matches := leaseRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		ip := match[1]
		leaseBody := match[2]
		lease := Lease{IP: ip}

		for _, line := range strings.Split(leaseBody, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, " ", 2)
			if len(parts) != 2 {
				continue
			}
			key := parts[0]
			value := strings.Trim(parts[1], ";")
			switch key {
			case "binding":
				lease.BindingState = strings.TrimPrefix(value, "state ")
			case "hardware":
				lease.HardwareEthernet = strings.TrimPrefix(value, "ethernet ")
			case "uid":
				lease.UID = convertOctalStringToHex(value)
			case "set":
				if strings.HasPrefix(value, "vendor-class-identifier = ") {
					lease.VendorClassIdentifier = strings.TrimPrefix(value, "vendor-class-identifier = ")
				}
			case "client-hostname":
				lease.Hostname = value
			}
		}
		leases = append(leases, lease)
	}
	return leases
}

func DoDhcpLeasePage(c *gin.Context) {
	c.HTML(http.StatusOK, "page.html", nil)
}

func DoDhcpLeaseApi(c *gin.Context) {
	// Read the file content
	file, err := os.Open("/var/lib/dhcp/dhcpd.leases")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		content.WriteString(line + "\n")
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Parse the leases
	leases := parseLeases(content.String())

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取结果成功",
		"data": leases,
	})
}

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	// 初始化路由
	r := gin.New()

	r.StaticFS("/css", http.Dir("./web/static/css"))
	r.StaticFS("/js", http.Dir("./web/static/js"))
	r.StaticFile("/favicon.ico", "./web/static/favicon.ico")

	r.LoadHTMLGlob("./web/template/*.html")

	r.GET("/api/leases", DoDhcpLeaseApi)
	r.GET("/", DoDhcpLeasePage)

	return r
}

func main() {

	// 初始化服务器
	server := http.Server{
		Addr:    ":8080",
		Handler: InitRouter(),
	}

	quit := make(chan os.Signal, 1)

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			quit <- syscall.SIGTERM
		}
	}()

	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Waiting for shutdown finishing...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Shutdown server err: %v.", err)
	}
	fmt.Println("Server shutdown succeed.")
}
