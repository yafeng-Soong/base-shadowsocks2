package statistic

import (
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	MAX_CONNECTS = 15
	MAX_INTERVAL = 3 * time.Second
)

type Metadata struct {
	SrcIP   net.IP `json:"source_IP"`
	DstIP   net.IP `json:"destination_IP"`
	SrcPort string `json:"source_port"`
	DstPort string `json:"destination_port"`
	Host    string `json:"destination_host"`
}

func (m *Metadata) LocalAndRemoteAddress() string {
	return strings.Join([]string{m.SrcIP.String(), m.DstString(), m.DstPort}, "-")
}

func (m *Metadata) RemoteAddress() string {
	return net.JoinHostPort(m.DstString(), m.DstPort)
}

func (m *Metadata) SourceAddress() string {
	return net.JoinHostPort(m.SrcIP.String(), m.SrcPort)
}

func (m *Metadata) DstString() string {
	if m.Host != "" {
		return m.Host
	} else if m.DstIP != nil {
		return m.DstIP.String()
	} else {
		return "<nil>"
	}
}

type Packet struct {
	Length    int   `json:"length"`
	Timestamp int64 `json:"timestamp"`
	Direction int   `json:"direction"`
}

type InfoRecord struct {
	Id            int    `gorm:"column:id;primaryKey"`
	UUID          string `gorm:"column:uuid"`
	StartTime     int64  `gorm:"column:start_time"`
	EndTime       int64  `gorm:"column:end_time"`
	Duration      int64  `gorm:"column:duration"`
	SrcIP         string `gorm:"column:src_ip"`
	DstIP         string `gorm:"column:dst_ip"`
	SrcPort       int    `gorm:"column:src_port"`
	DstPort       int    `gorm:"column:dst_port"`
	Hostname      string `gorm:"column:hostname"`
	UploadBytes   int64  `gorm:"column:upload_bytes"`
	DownloadBytes int64  `gorm:"column:download_bytes"`
	UploadNums    int64  `gorm:"column:upload_nums"`
	DownloadNums  int64  `gorm:"column:download_nums"`
	FilePath      string `gorm:"column:file_path"`
}

func (ir *InfoRecord) TableName() string {
	return "track_info"
}

func NewInfoRecord(tf *TrackerInfo, filename string) *InfoRecord {
	srcPort, _ := strconv.Atoi(tf.Metadata.SrcPort)
	dstPort, _ := strconv.Atoi(tf.Metadata.DstPort)
	return &InfoRecord{
		UUID:          tf.UUID.String(),
		StartTime:     tf.StartTime,
		EndTime:       tf.EndTime,
		Duration:      int64(tf.Duration),
		SrcIP:         tf.Metadata.SrcIP.String(),
		SrcPort:       srcPort,
		DstIP:         tf.Metadata.DstIP.String(),
		DstPort:       dstPort,
		Hostname:      tf.Metadata.Host,
		UploadBytes:   tf.UploadTotal,
		DownloadBytes: tf.DownloadTotal,
		UploadNums:    tf.UploadNums,
		DownloadNums:  tf.DownloadNums,
		FilePath:      filename,
	}
}

// type JobTimer struct {
// 	Timer *time.Timer
// 	Job   func(v ...interface{})
// }

// func NewFileTimer(interval int, job func(v ...interface{})) *JobTimer {
// 	return &JobTimer{
// 		Timer: time.NewTimer(time.Duration(interval) * time.Second),
// 		Job:   job,
// 	}
// }

// func (jt *JobTimer) Start(v ...interface{}) {
// 	go func() {
// 		<-jt.Timer.C
// 		jt.Job(v...)
// 	}()
// }

// type TrakerInfoMap struct {
// 	mutex sync.Mutex
// 	data  map[string]*TrakerInfoList
// }

// func (tm *TrakerInfoMap) Store(info *trackerInfo) {
// 	tm.mutex.Lock()
// 	defer tm.mutex.Unlock()
// 	key := info.Metadata.LocalAndRemoteAddress()
// 	list, ok := tm.data[key]
// 	if !ok {
// 		list = NewInfoQueue(MAX_INTERVAL, tm.SaveInfoList)
// 	}
// 	list.AddInfo(info)
// 	if list.Len() == MAX_CONNECTS {
// 		tm.saveInfoList(list)
// 	}
// }

// func (tm *TrakerInfoMap) saveInfoList(list *TrakerInfoList) {
// 	for _, info := range list.List {
// 		fmt.Println(info.UUID, "saved into file")
// 	}
// }

// func (tm *TrakerInfoMap) SaveInfoList(list *TrakerInfoList) {
// 	tm.mutex.Lock()
// 	defer tm.mutex.Unlock()
// 	tm.saveInfoList(list)
// }

// type TrakerInfoList struct {
// 	Timer    *time.Timer
// 	List     []*trackerInfo
// 	Job      func(list *TrakerInfoList)
// 	interval time.Duration
// }

// func NewInfoQueue(interval time.Duration, job func(list *TrakerInfoList)) *TrakerInfoList {
// 	queue := &TrakerInfoList{
// 		Timer:    time.NewTimer(interval),
// 		List:     make([]*trackerInfo, 0, 10),
// 		Job:      job,
// 		interval: interval,
// 	}
// 	queue.listen()
// 	return queue
// }

// func (tl *TrakerInfoList) listen() {
// 	go func() {
// 		<-tl.Timer.C
// 		// tl.Output <- tl.List
// 		tl.Job(tl)
// 	}()
// }

// func (tl *TrakerInfoList) update_timer() {
// 	tl.Timer.Reset(tl.interval)
// }

// func (tl *TrakerInfoList) AddInfo(info *trackerInfo) {
// 	tl.List = append(tl.List, info)
// }

// func (tl *TrakerInfoList) Len() int {
// 	return len(tl.List)
// }
