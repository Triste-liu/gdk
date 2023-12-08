package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/triste-liu/gdk/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/soft_delete"
	"runtime"
	"time"
)

type DefaultModel struct {
	ID        int                   `gorm:"primarykey;comment:主键" json:"id"`
	CreatedAt UnixTime              `gorm:"TYPE:TIMESTAMP;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt UnixTime              `gorm:"TYPE:TIMESTAMP;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	IsDeleted soft_delete.DeletedAt `gorm:"default:0;COMMENT:删除时间;softDelete:flag" json:"-"`
}

// PagePayload 分页查询
type PagePayload struct {
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	SearchKey string `json:"search_key"`
}

// PageData 分页响应体
type PageData struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
}

func PageQuery(db *gorm.DB, page PagePayload, data interface{}, model interface{}) (p PageData) {
	if page.Limit > 100 {
		page.Limit = 100
	}
	if page.Limit == 0 {
		page.Limit = 10
	}
	if model != nil {
		db = db.Model(model)
	}
	db.Count(&p.Total)
	if r := db.Limit(page.Limit).Offset(page.Offset).Find(&data); r.RowsAffected == 0 {
		p.Data = make([]interface{}, 0)
		return
	}
	p.Data = data
	return
}

type ClientConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DB       string
}

// 自定义logger
//type Interface interface {
//	LogMode(LogLevel) Interface
//	Info(context.Context, string, ...interface{})
//	Warn(context.Context, string, ...interface{})
//	Error(context.Context, string, ...interface{})
//	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
//}

type LoggerConfig struct {
	SlowThreshold             time.Duration // 慢查询阈值
	IgnoreRecordNotFoundError bool          // 忽略未找到错误
	ParameterizedQueries      bool          // sql是否打印参数
	Level                     gormLogger.LogLevel
}

func (l LoggerConfig) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.Level = level
	return l
}

func (l LoggerConfig) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Info(msg, data...)
}

// Warn print warn messages
func (l LoggerConfig) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Warning(msg, data...)
}

// Error print error messages
func (l LoggerConfig) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Error(msg, data...)
}

// Trace print sql message
func (l LoggerConfig) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	pc, file, line, _ := runtime.Caller(3)
	fn := runtime.FuncForPC(pc)
	caller := fmt.Sprintf("%s:%s:%d", file, fn.Name(), line)
	var logText string
	switch {
	case err != nil && (!errors.Is(err, gormLogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):

		if rows == -1 {
			logText = fmt.Sprintf("[%s] %s [%.3f ms] [%s]", caller, err, float64(elapsed.Microseconds())/1e3, sql)
		} else {
			logText = fmt.Sprintf("[%s] %s [%.3f ms] [rows:%d] [%s]", caller, err, float64(elapsed.Microseconds())/1e3, rows, sql)
		}
		log.Error(logText)
	case elapsed >= l.SlowThreshold && l.SlowThreshold != 0:
		if rows == -1 {
			logText = fmt.Sprintf("[%s] SLOW SQL >= %v [%.3f ms]  [%s]", caller, l.SlowThreshold, float64(elapsed.Microseconds())/1e3, sql)
		} else {
			logText = fmt.Sprintf("[%s] SLOW SQL >= %v [%.3f ms]  [rows:%d] [%s]", caller, l.SlowThreshold, float64(elapsed.Microseconds())/1e3, rows, sql)
		}
		log.Warning(logText)
	case l.Level == gormLogger.Info:
		if rows == -1 {
			logText = fmt.Sprintf("[%s] [%.3f ms] [%s]", caller, float64(elapsed.Microseconds())/1e3, sql)

		} else {
			logText = fmt.Sprintf("[%s] [%.3f ms] [rows:%d] [%s]", caller, float64(elapsed.Microseconds())/1e3, rows, sql)

		}
		log.Info(logText)
	}
}

func (l LoggerConfig) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

var session *gorm.DB

func Connect(clientConfig ClientConfig, loggerConfig LoggerConfig) {
	log.Info("init database")
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", clientConfig.User,
		clientConfig.Password, clientConfig.Host, clientConfig.Port, clientConfig.DB)
	s, err := gorm.Open(mysql.Open(url), &gorm.Config{Logger: loggerConfig})
	if err != nil {
		log.Error("database open error：%v", err)
	}
	sqlDB, _ := s.DB()
	err = sqlDB.Ping()
	if err != nil {
		log.Error("database connection error：%v", err)
		return
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	session = s
	log.Info("init database success")
}
func Session() *gorm.DB {
	if session == nil {
		log.Panic("no database connection,execute the \"Connect\" function")
	}
	return session
}
