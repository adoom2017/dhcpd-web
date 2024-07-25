## 安装DHCP服务软件
```BASH
# 基于Ubuntu
sudo apt-get install isc-dhcp-server
```

## 配置dhcp服务，这里是一个简单举例，根据自己情况进行修改

> IP地址范围是192.168.0.100到192.168.0.200  
> 子网掩码是255.255.255.0  
> 网关是192.168.0.1  
> DNS为114.114.114.114  
> 客户端的/etc/resolv.conf里面设置的search参数为pipci.com  
> 默认的IP地址租约时间1小时，最大租赁时间为2小时  
> 为硬件MAC地址0:0:c0:5d:bd:95保留IP地址为192.168.0.188 ，主机名为pipci

```conf
subnet 192.168.0.0 netmask 255.255.255.0 {
	range 192.168.0.100 192.168.0.200;
	option subnet-mask 255.255.255.0;
	option routers 192.168.0.1;
	option domain-name "pipci.com";
	option domain-name-servers 114.114.114.114;
	default-lease-time 3600;
	max-lease-time 7200;
	host pipci {
		hardware ethernet 00:0c:29:27:c6:12;
		fixed-address 192.168.0.188;
	}
}
```

## 使用dhcpd-web
将库clone之后
```
cd dhcpd-web
make build
```

编译完成之后，可执行文件将在build目录下，直接运行该文件
```
./build/dhcpd-web
```
使用浏览器访问http://ip-addr:8080
