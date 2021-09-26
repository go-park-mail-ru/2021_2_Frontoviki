package repository

import (
	"errors"
	"log"
	"sync"
	"time"
	"yula/internal/config"
	"yula/internal/models"
	"yula/internal/pkg/session"

	"github.com/tarantool/go-tarantool"
)

type SessionRepository struct {
	pool          []*tarantool.Connection
	m             sync.RWMutex
	roundRobinCur uint32
}

func NewSessionRepository(cfg *config.TarantoolConfig) session.SessionRepository {
	opts := tarantool.Opts{User: cfg.TarantoolOpts.User, Pass: cfg.TarantoolOpts.Pass}
	conn, err := tarantool.Connect(cfg.TarantoolServerAddress, opts)

	// opts := tarantool.Opts{User: "admin", Pass: "pass"}
	// conn, err := tarantool.Connect("localhost:3301", opts)

	if err != nil {
		log.Fatalf("Connection refused")
	}

	var pool []*tarantool.Connection
	pool = append(pool, conn)

	return &SessionRepository{
		pool:          pool,
		m:             sync.RWMutex{},
		roundRobinCur: 0,
	}
}

func (sr *SessionRepository) AddNewConnectionToPool(cfg *config.TarantoolConfig) error {
	opts := tarantool.Opts{User: cfg.TarantoolOpts.User, Pass: cfg.TarantoolOpts.Pass}
	conn, err := tarantool.Connect(cfg.TarantoolServerAddress, opts)

	// opts := tarantool.Opts{User: "admin", Pass: "pass"}
	// conn, err := tarantool.Connect("localhost:3301", opts)

	if err != nil {
		return errors.New("connect error")
	}

	sr.pool = append(sr.pool, conn)
	return nil
}

func (sr *SessionRepository) Set(sess *models.Session) error {

	sr.m.Lock()
	conn := sr.pool[sr.roundRobinCur]
	_, err := conn.Insert("sessions", []interface{}{sess.Value, sess.UserId, sess.ExpiresAt.Unix()})
	sr.roundRobinCur = (sr.roundRobinCur + 1) % uint32(len(sr.pool))
	sr.m.Unlock()

	if err != nil {
		return errors.New("insert error")
	}

	return nil
}

func (sr *SessionRepository) Delete(sess *models.Session) error {
	sr.m.Lock()
	conn := sr.pool[sr.roundRobinCur]
	_, err := conn.Delete("sessions", "primary", []interface{}{sess.Value})
	sr.roundRobinCur = (sr.roundRobinCur + 1) % uint32(len(sr.pool))
	sr.m.Unlock()

	if err != nil {
		return errors.New("delete error")
	}

	return errors.New("error")
}

func (sr *SessionRepository) GetByValue(value string) (*models.Session, error) {
	sr.m.RLock()
	conn := sr.pool[sr.roundRobinCur]
	resp, err := conn.Select("sessions", "primary", 0, 1, tarantool.IterEq, []interface{}{value})
	sr.roundRobinCur = (sr.roundRobinCur + 1) % uint32(len(sr.pool))
	sr.m.RUnlock()

	if len(resp.Data) == 0 {
		return nil, nil
	}

	sess := models.Session{
		Value:     (((resp.Data[0]).([]interface{}))[0]).(string),
		UserId:    int64((((resp.Data[0]).([]interface{}))[1]).(uint64)),
		ExpiresAt: time.Unix(int64((((resp.Data[0]).([]interface{}))[2]).(uint64)), 0),
	}

	if err != nil {
		return nil, errors.New("select error")
	}

	return &sess, nil
}

// func main() {
// 	tar := NewSessionRepository(nil)
// 	sess := models.Session{
// 		Value:     "13441",
// 		UserId:    134134,
// 		ExpiresAt: time.Now().Add(1000 * time.Second),
// 	}
// 	tar.Set(&sess)
// 	selsess, _ := tar.GetByValue(sess.Value)
// 	fmt.Println(selsess.UserId)
// }
