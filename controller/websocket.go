package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	dao "github.com/sztu/mutli-table/DAO"
	"github.com/sztu/mutli-table/DTO"
	"github.com/sztu/mutli-table/pkg/code"
	"github.com/sztu/mutli-table/service"
	"go.uber.org/zap"
)

func init() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			draggingLocks.Range(func(key, value interface{}) bool {
				if time.Since(value.(LockInfo).Timestamp) > 30*time.Second {
					draggingLocks.Delete(key)
				}
				return true
			})
		}
	}()
}

// 新增拖放操作消息结构
type DragItemMoveMessage struct {
	Type       string `json:"type"`
	SheetID    int64  `json:"sheet_id"`
	DragItemID int64  `json:"drag_item_id"`
	TargetRow  int    `json:"target_row"`
	TargetCol  int    `json:"target_col"`
}

type DragItemMovedMessage struct {
	Type       string `json:"type"`
	SheetID    int64  `json:"sheet_id"`
	DragItemID int64  `json:"drag_item_id"`
	FromRow    int    `json:"from_row"`
	FromCol    int    `json:"from_col"`
	TargetRow  int    `json:"target_row"`
	TargetCol  int    `json:"target_col"`
	MovedBy    int64  `json:"moved_by"`
	IsPlaced   bool   `json:"is_placed"`
}

type Client struct {
	conn     *websocket.Conn // WebSocket连接实例
	userID   int64           // 用户唯一标识
	sheetID  int64           // 当前操作的表格ID
	send     chan []byte     // 发送消息的缓冲通道
	username string          // 用户名
}

var (
	// WebSocket连接升级配置
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024, // 读缓冲区大小
		WriteBufferSize: 1024, // 写缓冲区大小
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有跨域请求
		},
	}
	rooms         = sync.Map{} // sheetID => map[*Client]bool 按表格分组的客户端集合
	clients       = sync.Map{} // *websocket.Conn => *Client 所有连接的客户端
	onlineUsers   = sync.Map{} // sheetID => map 在线用户状态跟踪
	draggingLocks = sync.Map{} // dragItemID => 正在拖动的用户ID
)

type OnlineUser struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	Connections int    `json:"connections"` // 同一用户可能有多个连接
}

type DragItemLockResponse struct {
	Type       string `json:"type"`
	Success    bool   `json:"success"`
	DragItemID int64  `json:"drag_item_id"`
	Message    string `json:"message"`
}

type DragItemLockMessage struct {
	Type       string `json:"type"`
	DragItemID int64  `json:"drag_item_id"`
}

// 锁信息结构
type LockInfo struct {
	UserID    int64
	Timestamp time.Time
}

// 定义基础消息结构
type BaseMessage struct {
	Type string `json:"type"`
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type CellUpdateMessage struct {
	Type    string `json:"type"`
	SheetID int64  `json:"sheet_id"`
	Row     int    `json:"row"`
	Column  int    `json:"column"`
	Content string `json:"content"`
	Version int64  `json:"version"`
}

type CellUpdatedMessage struct {
	Type      string `json:"type"`
	SheetID   int64  `json:"sheet_id"`
	Row       int    `json:"row"`
	Column    int    `json:"column"`
	Content   string `json:"content"`
	UpdatedBy int64  `json:"updated_by"`
	Version   int64  `json:"version"`
}

type ConflictMessage struct {
	Type        string `json:"type"`
	SheetID     int64  `json:"sheet_id"`
	Row         int    `json:"row"`
	Column      int    `json:"column"`
	YourContent string `json:"your_content"`
	NewContent  string `json:"new_content"`
}

func WebSocketHandler(c *gin.Context) {
	// 用户认证（复用现有中间件逻辑）
	userIDValue, exists := c.Get("user_id")
	if !exists {
		ResponseErrorWithMsg(c, code.InvalidAuth, "用户未登录")
		return
	}
	currentUserID, ok := userIDValue.(int64)
	if !ok {
		ResponseErrorWithMsg(c, code.ServerError, "用户ID解析错误")
		return
	}

	// 获取sheetID并验证权限
	sheetID, err := strconv.ParseInt(c.Param("sheet_id"), 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, "invalid sheet_id")
		return
	}

	// 升级WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.L().Error("WebSocket升级失败", zap.Error(err))
		ResponseErrorWithMsg(c, code.InvalidParam, "WebSocket升级失败")
		return
	}

	user, err := dao.FindUserByID(context.Background(), currentUserID)
	if err != nil {
		zap.L().Error("查找用户失败", zap.Error(err))
		ResponseErrorWithMsg(c, code.InvalidParam, "查找用户失败")
		return
	}
	username := user.Username
	client := &Client{
		conn:     conn,
		userID:   currentUserID,
		sheetID:  sheetID,
		send:     make(chan []byte, 256),
		username: username,
	}
	updateOnlineUsers(sheetID, client.userID, client.username, true)
	// 注册客户端
	clients.Store(conn, client)
	room, _ := rooms.LoadOrStore(sheetID, &sync.Map{})
	room.(*sync.Map).Store(client, true)

	// 启动goroutine处理读写
	go client.writePump()
	go client.readPump()
}

func updateOnlineUsers(sheetID, userID int64, username string, add bool) {
	key := sheetID
	value, _ := onlineUsers.LoadOrStore(key, &sync.Map{})

	usersMap := value.(*sync.Map)
	if add {
		var count int
		if v, ok := usersMap.Load(userID); ok {
			count = v.(*OnlineUser).Connections + 1
		} else {
			count = 1
		}
		usersMap.Store(userID, &OnlineUser{
			UserID:      userID,
			Username:    username,
			Connections: count,
		})
	} else {
		if v, ok := usersMap.Load(userID); ok {
			user := v.(*OnlineUser)
			if user.Connections <= 1 {
				usersMap.Delete(userID)
			} else {
				user.Connections--
				usersMap.Store(userID, user)
			}
		}
	}
}

// writePump 消息写入循环
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		clients.Delete(c.conn)
		if room, ok := rooms.Load(c.sheetID); ok {
			room.(*sync.Map).Delete(c)
		}
		updateOnlineUsers(c.sheetID, c.userID, c.username, false)

		// 清理该用户的所有锁
		draggingLocks.Range(func(key, value interface{}) bool {
			if value.(LockInfo).UserID == c.userID {
				draggingLocks.Delete(key)
			}
			return true
		})
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 通道关闭时发送关闭消息
				c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}

			// 带重试的消息发送机制
			for i := 0; i < 3; i++ {
				if err := c.conn.WriteMessage(websocket.TextMessage, message); err == nil {
					break
				}
				time.Sleep(time.Duration(i*i) * time.Second) // 指数退避重试
			}
		case <-ticker.C:
			// 定时发送心跳ping
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump 消息读取循环
func (c *Client) readPump() {
	defer c.conn.Close()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// 网络错误计数器
	errorCount := 0
	maxRetries := 3

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				zap.L().Error("非预期关闭", zap.Error(err))
			}

			// 网络波动重试逻辑
			if errorCount < maxRetries {
				errorCount++
				time.Sleep(time.Duration(errorCount) * time.Second)
				continue
			}
			break
		}
		errorCount = 0 // 重置计数器

		// 先解析基础消息，提取 type 字段
		var base BaseMessage
		if err := json.Unmarshal(message, &base); err != nil {
			zap.L().Warn("基础消息解析失败", zap.Error(err))
			continue
		}

		switch base.Type {
		case "GET_USERS":
			// 获取并发送当前在线用户列表
			if u, ok := onlineUsers.Load(c.sheetID); ok {
				var users []*OnlineUser
				u.(*sync.Map).Range(func(key, value interface{}) bool {
					users = append(users, value.(*OnlineUser))
					return true
				})
				if userList, err := json.Marshal(users); err == nil {
					select {
					case c.send <- userList:
					default:
						zap.L().Warn("GET_USERS响应发送通道已满")
					}
				} else {
					zap.L().Error("用户列表序列化失败", zap.Error(err))
				}
			}
		case "CELL_UPDATE":
			var updateMsg CellUpdateMessage
			if err := json.Unmarshal(message, &updateMsg); err != nil {
				zap.L().Warn("单元格更新消息解析失败", zap.Error(err))
				continue
			}

			// 获取当前单元格最新状态
			currentCell, err := dao.GetCellWithVersion(context.Background(),
				updateMsg.SheetID,
				updateMsg.Row,
				updateMsg.Column)
			if err != nil {
				zap.L().Error("获取单元格状态失败", zap.Error(err))
				continue
			}

			// 版本检查
			if currentCell.Version != updateMsg.Version {
				// 发送冲突通知
				conflictMsg := ConflictMessage{
					Type:        "CELL_CONFLICT",
					SheetID:     updateMsg.SheetID,
					Row:         updateMsg.Row,
					Column:      updateMsg.Column,
					YourContent: updateMsg.Content,
					NewContent:  currentCell.Content,
				}
				conflictBytes, _ := json.Marshal(conflictMsg)
				c.send <- conflictBytes
				continue
			}

			// 更新数据库（需要修改 DAO 方法支持版本更新）
			newVersion := currentCell.Version + 1
			err = dao.UpdateCellWithVersion(context.Background(),
				updateMsg.SheetID,
				updateMsg.Content,
				updateMsg.Row,
				updateMsg.Column,
				newVersion,
				c.userID)
			if err != nil {
				zap.L().Error("单元格更新失败", zap.Error(err))
				continue
			}

			// 广播更新（携带新版本号）
			updatedMsg := CellUpdatedMessage{
				Type:      "CELL_UPDATED",
				SheetID:   updateMsg.SheetID,
				Row:       updateMsg.Row,
				Column:    updateMsg.Column,
				Content:   updateMsg.Content,
				UpdatedBy: c.userID,
				Version:   newVersion,
			}

			msgBytes, _ := json.Marshal(updatedMsg)
			if room, ok := rooms.Load(updateMsg.SheetID); ok {
				room.(*sync.Map).Range(func(k, _ interface{}) bool {
					client := k.(*Client)
					select {
					case client.send <- msgBytes:
					default:
						close(client.send)
						room.(*sync.Map).Delete(client)
					}
					return true
				})
			}
		case "DRAG_ITEM_MOVE":
			var moveMsg DragItemMoveMessage
			if err := json.Unmarshal(message, &moveMsg); err != nil {
				zap.L().Warn("拖放操作消息解析失败", zap.Error(err))
				continue
			}

			// 带超时的锁验证
			if lockInfo, ok := draggingLocks.Load(moveMsg.DragItemID); ok {
				info := lockInfo.(LockInfo)
				if time.Since(info.Timestamp) > 30*time.Second || info.UserID != c.userID {
					draggingLocks.Delete(moveMsg.DragItemID)
					ok = false
				}
				if !ok {
					errMsg := DragItemLockResponse{
						Type:       "DRAG_ITEM_LOCK",
						DragItemID: moveMsg.DragItemID,
						Success:    false,
						Message:    "操作失败：锁已超时或失效",
					}
					if errBytes, err := json.Marshal(errMsg); err == nil {
						c.send <- errBytes
					}
					continue
				}
				// 更新锁时间戳
				draggingLocks.Store(moveMsg.DragItemID, LockInfo{
					UserID:    c.userID,
					Timestamp: time.Now(),
				})
			} else {
				errMsg := DragItemLockResponse{
					Type:       "DRAG_ITEM_LOCK",
					DragItemID: moveMsg.DragItemID,
					Success:    false,
					Message:    "操作失败：未获取元素锁",
				}
				if errBytes, err := json.Marshal(errMsg); err == nil {
					c.send <- errBytes
				}
				continue
			}

			// 原有服务调用和广播逻辑
			err := service.MoveDragItem(context.Background(), c.userID, moveMsg.SheetID,
				moveMsg.DragItemID, &DTO.MoveDragItemRequest{
					TargetRow: moveMsg.TargetRow,
					TargetCol: moveMsg.TargetCol,
				})
			if err != nil {
				draggingLocks.Delete(moveMsg.DragItemID)
				zap.L().Error("拖放操作失败", zap.Error(err))
				continue
			}

			movedMsg := DragItemMovedMessage{
				Type:       "DRAG_ITEM_MOVED",
				SheetID:    moveMsg.SheetID,
				DragItemID: moveMsg.DragItemID,
				TargetRow:  moveMsg.TargetRow,
				TargetCol:  moveMsg.TargetCol,
				MovedBy:    c.userID,
				IsPlaced:   true,
			}

			msgBytes, _ := json.Marshal(movedMsg)
			if room, ok := rooms.Load(moveMsg.SheetID); ok {
				room.(*sync.Map).Range(func(k, _ interface{}) bool {
					client := k.(*Client)
					select {
					case client.send <- msgBytes:
					default:
						close(client.send)
						room.(*sync.Map).Delete(client)
					}
					return true
				})
			}

		case "GET_DRAG_ITEM":
			var lockMsg DragItemLockMessage
			if err := json.Unmarshal(message, &lockMsg); err != nil {
				zap.L().Warn("锁定拖拽元素请求解析失败", zap.Error(err))
				continue
			}

			resp := DragItemLockResponse{
				Type:       "DRAG_ITEM_LOCK",
				DragItemID: lockMsg.DragItemID,
			}

			if currentLock, ok := draggingLocks.Load(lockMsg.DragItemID); ok {
				info := currentLock.(LockInfo)
				if time.Since(info.Timestamp) > 30*time.Second {
					draggingLocks.Delete(lockMsg.DragItemID)
				} else {
					resp.Success = false
					resp.Message = fmt.Sprintf("元素已被用户 %d 锁定", info.UserID)
				}
			}

			if !resp.Success {
				// 存储带时间戳的锁信息
				draggingLocks.Store(lockMsg.DragItemID, LockInfo{
					UserID:    c.userID,
					Timestamp: time.Now(),
				})
				resp.Success = true
				resp.Message = "锁定成功"
			}

			if respBytes, err := json.Marshal(resp); err == nil {
				select {
				case c.send <- respBytes:
				default:
					zap.L().Warn("拖拽元素锁定响应发送通道已满")
				}
			}

		default:
			zap.L().Warn("未知消息类型", zap.String("type", base.Type))
		}
	}
}
