package main

import (
	"context"
	"os"
	"sort"
	"time"

	"github.com/grimdork/foreman/api"
	"github.com/grimdork/foreman/clients"
	ll "github.com/grimdork/loglines"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (srv *Server) openDB() error {
	var err error
	srv.db, err = pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	_, err = srv.db.Exec(context.Background(), initialTables)
	if err != nil {
		return err
	}

	// srv.LoadConfig()
	ll.Msg("Opened database connection.")
	return nil
}

func (srv *Server) closeDB() {
	srv.db.Close()
	ll.Msg("Closed database connection.")
}

// GetKeys returns all keys from the database.
func (srv *Server) GetKeys() ([]api.Key, error) {
	srv.Lock()
	defer srv.Unlock()

	rows, err := srv.db.Query(context.Background(), "select name,key,admin from keys")
	if err != nil {
		return nil, err
	}

	var keys []api.Key
	for rows.Next() {
		var k api.Key
		err = rows.Scan(&k.ID, &k.Key, &k.Admin)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].ID < keys[j].ID
	})

	return keys, nil
}

// GetKey returns a key and expiry time from the key table.
func (srv *Server) GetKey(name string) (api.Key, error) {
	srv.Lock()
	defer srv.Unlock()

	var key api.Key
	err := srv.db.QueryRow(context.Background(),
		"select key,expiry,admin from keys where name=$1",
		name).Scan(&key.Key, &key.Expiry, &key.Admin)
	if err != nil {
		return api.Key{}, err
	}

	key.ID = name
	return key, nil
}

// SetKey upserts a key.
func (srv *Server) SetKey(name, key string, admin bool) error {
	srv.Lock()
	defer srv.Unlock()
	sql := `insert into keys (name,key,admin) values ($1,$2,$3) on conflict(name) do update set key=$2,admin=$3`
	_, err := srv.db.Exec(context.Background(), sql, name, key, admin)
	return err
}

// DeleteKey from the database.
func (srv *Server) DeleteKey(name string) error {
	srv.Lock()
	defer srv.Unlock()
	_, err := srv.db.Exec(context.Background(), "delete from keys where name=$1", name)
	return err
}

// GetScouts returns a list of all scouts in the database.
func (srv *Server) GetScouts() (api.ScoutList, error) {
	srv.Lock()
	defer srv.Unlock()
	var list []api.ScoutListEntry
	sql := `select hostname,port,interval,last_check,failed,status,assignee,assigned from scouts`
	rows, err := srv.db.Query(context.Background(), sql)
	if err != nil {
		return api.ScoutList{}, err
	}

	defer rows.Close()
	for rows.Next() {
		scout := api.ScoutListEntry{}
		err := rows.Scan(
			&scout.Hostname, &scout.Port, &scout.Interval, &scout.LastCheck,
			&scout.Failed, &scout.Status, &scout.Assignee, &scout.Assigned,
		)
		if err != nil {
			return api.ScoutList{}, err
		}

		list = append(list, scout)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Hostname < list[j].Hostname
	})
	return api.ScoutList{Scouts: list}, rows.Err()
}

// GetScout from the database.
func (srv *Server) GetScout(hostname string) (*api.ScoutListEntry, error) {
	srv.Lock()
	defer srv.Unlock()
	scout := &api.ScoutListEntry{Hostname: hostname}
	sql := `select port,interval,last_check,failed,status from scouts where hostname=$1`
	err := srv.db.QueryRow(context.Background(), sql, hostname).Scan(
		&scout.Port, &scout.Interval, &scout.LastCheck, &scout.Failed,
		&scout.Status,
	)
	if err != nil {
		return nil, err
	}

	return scout, nil
}

// SetScout creates a new scout in the database or updates its interval.
func (srv *Server) SetScout(hostname string, port, interval int) error {
	srv.Lock()
	defer srv.Unlock()
	sql := `insert into scouts (hostname,port,interval) values ($1,$2,$3) on conflict(hostname) do update set interval=$3`
	_, err := srv.db.Exec(context.Background(), sql, hostname, port, interval)
	return err
}

// DeleteScout from the database.
func (srv *Server) DeleteScout(hostname string) error {
	srv.Lock()
	defer srv.Unlock()
	_, err := srv.db.Exec(context.Background(), "delete from scouts where hostname=$1", hostname)
	return err
}

// SetScoutFailure switches to a failure state.
func (srv *Server) SetScoutFailure(hostname string, status uint8) error {
	srv.Lock()
	defer srv.Unlock()
	sql := `update scouts set failed=$2,status=$3 where hostname=$1`
	_, err := srv.db.Exec(context.Background(), sql, hostname, time.Now(), status)
	return err
}

// GetCanaries returns a list of all canaries in the database.
func (srv *Server) GetCanaries() (api.CanaryList, error) {
	srv.Lock()
	defer srv.Unlock()
	var list []api.CanaryListEntry
	sql := `select c.name,interval,last_check,failed,status,assignee,assigned,key
from canaries as c left join keys as k on c.name=k.name`
	rows, err := srv.db.Query(context.Background(), sql)
	if err != nil {
		return api.CanaryList{}, err
	}

	defer rows.Close()
	for rows.Next() {
		canary := api.CanaryListEntry{}
		err := rows.Scan(
			&canary.Hostname, &canary.Interval, &canary.LastCheck, &canary.Failed,
			&canary.Status, &canary.Assignee, &canary.Assigned, &canary.Key,
		)
		if err != nil {
			return api.CanaryList{}, err
		}

		list = append(list, canary)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Hostname < list[j].Hostname
	})
	return api.CanaryList{Canaries: list}, rows.Err()
}

// GetCanary from the database.
func (srv *Server) GetCanary(name string) (*clients.Canary, error) {
	srv.Lock()
	defer srv.Unlock()
	canary := &clients.Canary{Client: clients.Client{Hostname: name}}
	sql := `select interval,last_check,failed,status,assignee,assigned,key
from canaries as c
left join keys on c.name=name`
	err := srv.db.QueryRow(context.Background(), sql, name).Scan(
		&canary.Interval, &canary.LastCheck, &canary.Failed, &canary.Status,
		&canary.Assignee, &canary.Assigned, &canary.Key,
	)
	if err != nil {
		return nil, err
	}

	return canary, err
}

// SetCanary creates a new canary in the database or updates its interval.
func (srv *Server) SetCanary(name string, interval int) error {
	srv.Lock()
	defer srv.Unlock()
	sql := `insert into canaries(name,interval) values($1,$2) on conflict(name) do update set interval=$2`
	_, err := srv.db.Exec(context.Background(), sql, name, interval)
	return err
}

// DeleteCanary from the database.
func (srv *Server) DeleteCanary(name string) error {
	srv.Lock()
	defer srv.Unlock()
	_, err := srv.db.Exec(context.Background(), "delete from canaries where name=$1", name)
	if err != nil {
		return err
	}

	_, err = srv.db.Exec(context.Background(), "delete from keys where name=$1", name)
	return err
}
