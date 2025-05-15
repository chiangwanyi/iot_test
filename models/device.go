package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chiangwanyi/iot_test/db"
)

// Device 定义设备结构体
type Device struct {
	ID           int            `json:"id"`
	SN           string         `json:"sn"`
	Name         string         `json:"name"`
	Description  sql.NullString `json:"description"`
	Model        string         `json:"model"`
	Type         string         `json:"type"`
	IPAddr       sql.NullString `json:"ipaddr"`
	Status       string         `json:"status"`
	LastOnlineAt sql.NullTime   `json:"last_online_at"`
	UpdatedAt    sql.NullTime   `json:"updated_at"`
	CreatedAt    sql.NullTime   `json:"created_at"`
}

// API响应模型
type DeviceResponse struct {
	ID           int     `json:"id"`
	SN           string  `json:"sn"`
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	Model        string  `json:"model"`
	Type         string  `json:"type"`
	IPAddr       *string `json:"ipaddr,omitempty"`
	Status       string  `json:"status"`
	LastOnlineAt *string `json:"last_online_at,omitempty"`
	UpdatedAt    string  `json:"updated_at"`
	CreatedAt    string  `json:"created_at"`
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
            id INTEGER PRIMARY KEY AUTOINCREMENT, -- 自增序列ID
            sn TEXT UNIQUE,
            name TEXT NOT NULL,
            description TEXT,
            model TEXT NOT NULL,
            type TEXT NOT NULL,
            ipaddr TEXT,
            status TEXT NOT NULL default "离线",
            last_online_at TIMESTAMP,
            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		log.Fatalf("创建设备表失败: %v", err)
	}

	fmt.Println("数据库表初始化完成")
}

func (m *DeviceModel) CreateDevice(device *Device) error {
	query := `INSERT INTO devices (sn, name, description, model, type) VALUES (?, ?, ?, ?, ?)`
	_, err := m.DB.Exec(query, device.SN, device.Name, device.Description, device.Model, device.Type)
	return err
}

// GetAllDevices 查询并返回所有设备列表
func (m *DeviceModel) GetAllDevices() ([]DeviceResponse, error) {
	query := `SELECT * FROM devices`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []DeviceResponse
	for rows.Next() {
		var device Device
		err = rows.Scan(&device.ID, &device.SN, &device.Name, &device.Description, &device.Model, &device.Type, &device.IPAddr, &device.Status, &device.LastOnlineAt, &device.UpdatedAt, &device.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, device.convert())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// GetDevicesWithPagination 分页查询设备列表
func (m *DeviceModel) GetDevicesWithPage(page, pageSize int) ([]DeviceResponse, error) {
	offset := (page - 1) * pageSize
	query := `SELECT * FROM devices LIMIT ? OFFSET ?`
	rows, err := m.DB.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []DeviceResponse
	for rows.Next() {
		var device Device
		err = rows.Scan(&device.ID, &device.SN, &device.Name, &device.Description, &device.Model, &device.Type, &device.IPAddr, &device.Status, &device.LastOnlineAt, &device.UpdatedAt, &device.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, device.convert())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// 将Device转换为DeviceResponse
func (device *Device) convert() DeviceResponse {
	var description *string
	if device.Description.Valid {
		description = &device.Description.String
	}

	var ipaddr *string
	if device.IPAddr.Valid {
		ipaddr = &device.IPAddr.String
	}

	var lastOnlineAt *string
	if device.LastOnlineAt.Valid {
		formattedTime := device.LastOnlineAt.Time.Format("2006-01-02 15:04:05")
		lastOnlineAt = &formattedTime
	}

	updatedAt := device.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	createdAt := device.CreatedAt.Time.Format("2006-01-02 15:04:05")

	return DeviceResponse{
		ID:           device.ID,
		SN:           device.SN,
		Name:         device.Name,
		Description:  description,
		Model:        device.Model,
		Type:         device.Type,
		IPAddr:       ipaddr,
		Status:       device.Status,
		LastOnlineAt: lastOnlineAt,
		UpdatedAt:    updatedAt,
		CreatedAt:    createdAt,
	}
}
