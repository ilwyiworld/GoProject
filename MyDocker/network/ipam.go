package network

import (
	"net"
	"os"
	"path"
	"strings"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

const ipamDefaultAllocatorPath = "/var/run/mydocker/network/ipam/subnet.json"

type IPAM struct {
	//分配文件存放位置
	SubnetAllocatorPath string
	//网段和位图算法的数组map，key是网段，value是分配的位图数组
	Subnets *map[string]string
}

//初始化一个IPAM的对象，默认使用ipamDefaultAllocatorPath
var ipAllocator = &IPAM{
	SubnetAllocatorPath: ipamDefaultAllocatorPath,
}

func (ipam *IPAM) load() error {
	//检查存储文件状态，如果不存在，则说明之前没有分配，不需要加载
	if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	//打开并读取存储文件
	subnetConfigFile, err := os.Open(ipam.SubnetAllocatorPath)
	defer subnetConfigFile.Close()
	if err != nil {
		return err
	}
	subnetJson := make([]byte, 2000)
	n, err := subnetConfigFile.Read(subnetJson)
	if err != nil {
		return err
	}
	//将文件中的内容反序列化出IP的分配信息
	err = json.Unmarshal(subnetJson[:n], ipam.Subnets)
	if err != nil {
		log.Errorf("Error dump allocation info, %v", err)
		return err
	}
	return nil
}

//存储网段地址分配信息
func (ipam *IPAM) dump() error {
	//path.Split能够分隔目录和文件
	ipamConfigFileDir, _ := path.Split(ipam.SubnetAllocatorPath)
	if _, err := os.Stat(ipamConfigFileDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(ipamConfigFileDir, 0644)
		} else {
			return err
		}
	}
	//os.O_TRUNC 表示存在则清空  os.O_CREATE表示不存在则创建
	subnetConfigFile, err := os.OpenFile(ipam.SubnetAllocatorPath, os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0644)
	defer subnetConfigFile.Close()
	if err != nil {
		return err
	}

	//序列化ipam对象到json串
	ipamConfigJson, err := json.Marshal(ipam.Subnets)
	if err != nil {
		return err
	}

	_, err = subnetConfigFile.Write(ipamConfigJson)
	if err != nil {
		return err
	}

	return nil
}

func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
	// 存放网段中地址分配信息的数组
	ipam.Subnets = &map[string]string{}

	// 从文件中加载已经分配的网段信息
	err = ipam.load()
	if err != nil {
		log.Errorf("Error dump allocation info, %v", err)
	}

	//将网段的字符串转换成net.IPNet的对象
	_, subnet, _ = net.ParseCIDR(subnet.String())

	one, size := subnet.Mask.Size()		//返回网段前面固定位的长度和子网掩码的总长度
	//例如192.168.0.1/24 返回 24	32

	//如果之前没有分配过这个网段，则初始化网段的分配配置
	//用"0"填满这个网段的配置，1 << uint8(size - one)表示这个网段中的可用地址数量
	if _, exist := (*ipam.Subnets)[subnet.String()]; !exist {
		(*ipam.Subnets)[subnet.String()] = strings.Repeat("0", 1 << uint8(size - one))
	}

	//遍历网段的位图数组
	for c := range((*ipam.Subnets)[subnet.String()]) {
		//找到数组中为"0"的项和数字序号，即可分配的IP
		if (*ipam.Subnets)[subnet.String()][c] == '0' {
			//设置这个为"0"的序号值为"1"，即分配这个ip
			//go的字符串，创建之后不能修改，所以通过转换成byte数组，修改后再转换成字符串赋值
			ipalloc := []byte((*ipam.Subnets)[subnet.String()])
			ipalloc[c] = '1'
			(*ipam.Subnets)[subnet.String()] = string(ipalloc)
			//这里的ip为初始ip，比如对于网段192.168.0.0/16，这里就是192.168.0.0
			ip = subnet.IP
			/*
			通过网段的IP与上面的偏移相加计算出分配的IP地址，由于IP地址是一个uint的一个数组
			需要通过数组中的每一项加所需要的值，比如网段是172.16.0.0/12 数组序号是65555
			那么在[172,16,0,0]上依次加上[unit8(65555>>24)，unit8(65555>>16)，unit8(65555>>8)，unit8(65555>>0)]
			即[0,1,0,19] 那么获得的IP就是172.17.0.19
			 */
			for t := uint(4); t > 0; t-=1 {
				[]byte(ip)[4-t] += uint8(c >> ((t - 1) * 8))
			}
			//由于此处IP是从1开始分配的，所以最后再加1，最终得到分配的IP是172.17.0.20
			ip[3]+=1
			break
		}
	}
	//将分配结果保存到文件中
	ipam.dump()
	return
}

func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnets = &map[string]string{}

	_, subnet, _ = net.ParseCIDR(subnet.String())
	//从文件中加载网段的分配信息
	err := ipam.load()
	if err != nil {
		log.Errorf("Error dump allocation info, %v", err)
	}

	//计算IP地址再网段位图数组中的索引位置
	c := 0
	releaseIP := ipaddr.To4()
	//将IP地址转换成4字节的表示方式
	//由于IP是从1开始分配的,所以转换成索引应减1
	releaseIP[3]-=1
	for t := uint(4); t > 0; t-=1 {
		//与分配IP相反，释放IP获得索引的方式是IP地址的每一位相减后分别左移将对应的数值相加到索引上
		//subnet.IP是第一个IP 即网关IP 如172.17.0.0
		c += int(releaseIP[t-1] - subnet.IP[t-1]) << ((4-t) * 8)
	}

	//将分配的位图数组中索引位置的值置为0
	ipalloc := []byte((*ipam.Subnets)[subnet.String()])
	ipalloc[c] = '0'
	(*ipam.Subnets)[subnet.String()] = string(ipalloc)

	//保存IP分配信息
	ipam.dump()
	return nil
}