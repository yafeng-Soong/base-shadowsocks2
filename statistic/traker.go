package statistic

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shadowsocks/go-shadowsocks2/socks"
)

type TrackerInfo struct {
	UUID            uuid.UUID     `json:"id"`
	Metadata        *Metadata     `json:"metadata"`
	Start           time.Time     `json:"start"`
	End             time.Time     `json:"end"`
	StartTime       int64         `json:"start_time"`
	EndTime         int64         `json:"end_time"`
	Duration        time.Duration `json:"duration"`
	UploadTotal     int64         `json:"upload_bytes"`
	DownloadTotal   int64         `json:"download_bytes"`
	UploadNums      int64         `json:"upload_nums"`
	DownloadNums    int64         `json:"download_nums"`
	UploadPackets   []*Packet     `json:"upload_packets"`
	DownloadPackets []*Packet     `json:"download_packets"`
}

func (tf *TrackerInfo) OutPath() string {
	date := tf.Start.Format("2006-01-02")
	dst := tf.Metadata.DstString()
	src := tf.Metadata.SrcIP.String()
	// 算与开始时间最接近的 以gap秒分割的前一个时间点
	second := int(tf.Start.Unix())
	a := second / CONFIG.Gap
	times := a * CONFIG.Gap
	return path.Join(CONFIG.OutPath, date, dst, src, strconv.Itoa(times))
}

type TcpTracker struct {
	net.Conn
	*TrackerInfo
	outChan chan<- *TrackerInfo
}

func (t *TcpTracker) Read(b []byte) (int, error) {
	n, err := t.Conn.Read(b)
	// log.Printf(
	// 	"%s<--->%s upload %d bytes:	",
	// 	t.Metadata.SourceAddress(),
	// 	t.Metadata.RemoteAddress(),
	// 	n,
	// )
	t.UploadTotal += int64(n)
	t.UploadNums += 1
	t.UploadPackets = append(
		t.UploadPackets,
		&Packet{Length: n, Timestamp: time.Now().UnixNano(), Direction: 1},
	)
	return n, err
}

func (t *TcpTracker) Write(b []byte) (int, error) {
	n, err := t.Conn.Write(b)
	// log.Printf(
	// 	"%s<--->%s download %d bytes:	",
	// 	t.Metadata.SourceAddress(),
	// 	t.Metadata.RemoteAddress(),
	// 	n,
	// )
	t.DownloadTotal += int64(n)
	t.DownloadNums += 1
	t.DownloadPackets = append(
		t.DownloadPackets,
		&Packet{Length: n, Timestamp: time.Now().UnixNano(), Direction: -1},
	)
	return n, err
}

// 收尾工作，将这条连接的统计信息放到chan中
func (t *TcpTracker) WindUp() {
	t.End = time.Now()
	t.Duration = t.End.Sub(t.Start)
	t.StartTime = t.Start.UnixNano()
	t.EndTime = t.End.UnixNano()
	// b, err := json.Marshal(t.trackerInfo)
	// if err != nil {
	// 	return
	// }
	// log.Println(string(b))
	// log.Println(t.OutPath())
	t.outChan <- t.TrackerInfo
}

func NewTcpTracker(conn net.Conn, out chan<- *TrackerInfo, address string) *TcpTracker {
	target := socks.ParseAddr(address)
	metadata := parseSocksAddr(target)
	if ip, port, err := parseAddr(conn.RemoteAddr().String()); err == nil {
		metadata.SrcIP = ip
		metadata.SrcPort = port
	}
	uuid, _ := uuid.NewV4()
	// log.Printf("%s %s", metadata.SrcIP.String(), metadata.SrcPort)
	return &TcpTracker{
		outChan: out,
		Conn:    conn,
		TrackerInfo: &TrackerInfo{
			UUID:            uuid,
			Start:           time.Now(),
			Metadata:        metadata,
			UploadTotal:     0,
			DownloadTotal:   0,
			UploadNums:      0,
			DownloadNums:    0,
			UploadPackets:   make([]*Packet, 0, 25),
			DownloadPackets: make([]*Packet, 0, 25),
		},
	}
}

// 从chan中读取每条的统计信息
func HandleMetric(ch <-chan *TrackerInfo) {
	for info := range ch {
		outpath := info.OutPath()
		if ok, _ := PathExists(outpath); !ok {
			err := os.MkdirAll(outpath, os.ModePerm)
			if err != nil {
				log.Println("创建文件夹", outpath, "出错")
				continue
			}
		}
		name := info.UUID.String()[:8] + ".json"
		filename := path.Join(outpath, name)
		go writeToFile(name, filename, info)
		go writeToDB(info, filename, name)
	}
}

func writeToFile(name, filename string, info *TrackerInfo) {

	f, err := os.Create(filename)
	if err != nil {
		log.Println("创建文件", name, "出错")
		return
	}
	defer f.Close()
	b, err := json.Marshal(info)
	if err != nil {
		log.Println("序列化", info.UUID, "出错")
		return
	}
	_, err = f.Write(b)
	if err != nil {
		log.Println("写入文件", name, "出错")
		return
	}
	log.Println(name, "写入完成")
}

func writeToDB(info *TrackerInfo, filename string, name string) {
	record := NewInfoRecord(info, filename)
	if err := DB.Create(record).Error; err != nil {
		log.Println(err.Error())
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
