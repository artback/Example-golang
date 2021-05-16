package repository

import (
	"bitburst/pkg/online"
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"log"
	"net"
	"net/url"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func startDatabase(tb testing.TB) *url.URL {
	tb.Helper()

	pgURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("myuser", "mypass"),
		Path:   "mydatabase",
	}
	q := pgURL.Query()
	q.Add("sslmode", "disable")
	pgURL.RawQuery = q.Encode()

	pool, err := dockertest.NewPool("")
	if err != nil {
		tb.Fatalf("Could not connect to docker: %v", err)
	}

	pw, _ := pgURL.User.Password()
	env := []string{
		"POSTGRES_USER=" + pgURL.User.Username(),
		"POSTGRES_PASSWORD=" + pw,
		"POSTGRES_DB=" + pgURL.Path,
	}

	resource, err := pool.Run("postgres", "13-alpine", env)
	if err != nil {
		tb.Fatalf("Could not start postgres container: %v", err)
	}
	tb.Cleanup(func() {
		err = pool.Purge(resource)
		if err != nil {
			tb.Fatalf("Could not purge container: %v", err)
		}
	})

	pgURL.Host = resource.Container.NetworkSettings.IPAddress

	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		pgURL.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}

	logWaiter, err := pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    resource.Container.ID,
		OutputStream: log.Writer(),
		ErrorStream:  log.Writer(),
		Stderr:       true,
		Stdout:       true,
		Stream:       true,
	})
	if err != nil {
		tb.Fatalf("Could not connect to postgres container logging output: %v", err)
	}

	tb.Cleanup(func() {
		err = logWaiter.Close()
		if err != nil {
			tb.Fatalf("Could not close container logging: %v", err)
		}
		err = logWaiter.Wait()
		if err != nil {
			tb.Fatalf("Could not wait for container logging to close: %v", err)
		}
	})

	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() (err error) {
		db, err := sql.Open("pgx", pgURL.String())
		if err != nil {
			return err
		}
		defer func() {
			cerr := db.Close()
			if err == nil {
				err = cerr
			}
		}()

		return db.Ping()
	})
	if err != nil {
		tb.Fatalf("Could not connect to postgres container: %v", err)
	}

	return pgURL
}

func Test_postgresRepository_DeleteOlder(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	tests := []struct {
		name    string
		time    time.Time
		setup   func(repo online.Repository)
		want    int64
		wantErr bool
	}{
		{
			name:  "empty database",
			time:  time.Now(),
			setup: func(repo online.Repository) {},
			want:  0,
		},
		{
			name: "delete 2",
			time: time.Now().Add(30 * time.Minute),
			setup: func(repo online.Repository) {
				repo.UpsertAll(
					context.Background(),
					[]int{1, 5},
					time.Now(),
				)
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewPostgresRepository(startDatabase(t))
			tt.setup(s)
			err = s.DeleteOlder(context.Background(), tt.time)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_postgresRepository_UpsertAll(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	tests := []struct {
		name    string
		ids     []int
		context context.Context
		wantErr bool
		want    int64
	}{
		{
			name: "insert 2",
			ids:  []int{1, 2},
			want: 2,
		},
		{
			name: "insert empty",
			ids:  []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := NewPostgresRepository(startDatabase(t))
			if err != nil {
				t.Fatal(err)
			}
			err = repo.UpsertAll(context.Background(), tt.ids, time.Now())
			if (err != nil) != tt.wantErr {
				t.Errorf("UpsertAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				r := repo.(*postgresRepository)
				res, err := r.ExecContext(context.Background(), "select * from status")
				if err != nil {
					t.Fatal(err)
				}
				got, _ := res.RowsAffected()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("UpsertAll() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestNewPostgresRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	type args struct {
		url *url.URL
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "new with schema",
			args: args{
				startDatabase(t),
			},
			want: true,
		},
		{
			name: "error url",
			args: args{
				&url.URL{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPostgresRepository(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPostgresRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.want {
				t.Errorf("NewPostgresRepository() got = %v, want %v", got, tt.want)
			}
		})
	}
}
