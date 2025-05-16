package tcp_mgr

import (
	"bufio"
	"net"
	"sync"
	"time"
)

// ClientConn 客户端连接信息结构体
type ClientConn struct {
	Conn       net.Conn   // 底层连接
	CreateTime time.Time  // 连接建立时间
	CloseTime  *time.Time // 连接断开时间（nil表示未断开）
	Messages   []string   // 客户端发送的消息列表
	addr       string     // 客户端地址（缓存避免重复调用RemoteAddr）
}

// TcpMgr TCP连接管理器
type TcpMgr struct {
	listener net.Listener           // 监听对象
	conns    map[string]*ClientConn // 连接表（键：客户端地址）
	mu       sync.Mutex             // 并发安全锁
}

// NewTcpMgr 创建新的TCP管理器
func NewTcpMgr() *TcpMgr {
	return &TcpMgr{
		conns: make(map[string]*ClientConn),
	}
}

// Start 启动TCP服务监听（端口9009）
func (m *TcpMgr) Start() error {
	listener, err := net.Listen("tcp", ":9309")
	if err != nil {
		return err
	}
	m.listener = listener

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				// 监听关闭时退出循环
				return
			}
			m.handleNewConn(conn)
		}
	}()
	return nil
}

// handleNewConn 处理新连接
func (m *TcpMgr) handleNewConn(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	now := time.Now()

	m.mu.Lock()
	m.conns[addr] = &ClientConn{
		Conn:       conn,
		CreateTime: now,
		addr:       addr,
	}
	m.mu.Unlock()

	// 启动独立goroutine处理客户端消息
	go func() {
		defer func() {
			// 连接关闭时记录断开时间
			closeTime := time.Now()
			m.mu.Lock()
			if c, ok := m.conns[addr]; ok {
				c.CloseTime = &closeTime
			}
			m.mu.Unlock()
			conn.Close()
		}()

		reader := bufio.NewReader(conn)
		for {
			msg, err := reader.ReadString('\n') // 按行读取消息（假设客户端用\n结尾）
			if err != nil {
				return // 读取失败/连接关闭
			}
			m.mu.Lock()
			m.conns[addr].Messages = append(m.conns[addr].Messages, msg)
			m.mu.Unlock()
		}
	}()
}

// GetAllConnections 获取所有连接信息（用于外部查询）
func (m *TcpMgr) GetAllConnections() map[string]*ClientConn {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 返回副本避免外部修改原始数据
	result := make(map[string]*ClientConn, len(m.conns))
	for addr, conn := range m.conns {
		result[addr] = &ClientConn{
			CreateTime: conn.CreateTime,
			CloseTime:  conn.CloseTime,
			Messages:   append([]string(nil), conn.Messages...),
			addr:       conn.addr,
		}
	}
	return result
}

// Stop 停止服务（用于优雅关闭）
func (m *TcpMgr) Stop() error {
	if m.listener != nil {
		return m.listener.Close()
	}
	return nil
}
