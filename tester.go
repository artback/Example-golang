package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	client := &http.Client{Timeout: 1 * time.Second}

	for {
		time.Sleep(5 * time.Second)

		ids := make([]string, rng.Int31n(200))
		for i := range ids {
			ids[i] = strconv.Itoa(rng.Int() % 100)
		}
		body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"object_ids":[%s]}`, strings.Join(ids, ","))))
		resp, err := client.Post("http://localhost:9090/callback", "application/json", body)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_ = resp.Body.Close()
	}

}
