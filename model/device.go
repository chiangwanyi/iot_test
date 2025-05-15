package model

import (
	"fmt"
	"log"
	"time"

	"github.com/chiangwanyi/iot_test/db"
)

// Device 定义设备结构体
type Device struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	SerialNumber string    `json:"serial_number"`
	Type         string    `json:"type"`
	IPAddress    string    `json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
}

func (d *Device) TableName() string {
	return "devices" // 数据库表名
}

// CreateTables 创建数据库表
func CreateTables() {
	// 创建设备表
	_, err := db.SqliteConn.Exec(`
        CREATE TABLE IF NOT EXISTS devices (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            serial_number TEXT UNIQUE,
            type TEXT DEFAULT 'offline',
            ip_address TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		log.Fatalf("创建设备表失败: %v", err)
	}

	fmt.Println("数据库表初始化完成")
}

func (device *Device) Create() (int64, error) {
	result, err := db.SqliteConn.Exec(`
        INSERT INTO devices (name, serial_number, ip_address)
        VALUES (?, ?, ?, ?)
        `, device.Name, device.SerialNumber, device.Type, device.IPAddress)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
