package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/chiangwanyi/iot_test/db"
)

// Device 定义设备结构体
type Device struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	SN        string    `json:"sn"`
	Type      string    `json:"type"`
	IPAddress string    `json:"ip_address"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// DeviceModel 设备模型操作
type DeviceModel struct {
	DB *sql.DB
}

// CreateTables 创建数据库表
func (m *DeviceModel) CreateTables() {
	// 创建设备表
	_, err := db.SqliteConn.Exec(`
        CREATE TABLE IF NOT EXISTS devices (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            sn TEXT UNIQUE,
            type TEXT NOT NULL,
            ip_address TEXT,
            status TEXT DEFAULT 'offline',
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		log.Fatalf("创建设备表失败: %v", err)
	}

	fmt.Println("数据库表初始化完成")
}

func (m *DeviceModel) CreateDevice(device *Device) error {
	query := `INSERT INTO devices (name, sn, type, ip_address) VALUES (?, ?, ?, ?)`
	_, err := m.DB.Exec(query, device.Name, device.SN, device.Type, device.IPAddress)
	return err
}
