package kagomi

import (
    "testing"
    "gopkg.in/redis.v3"
)

func TestHello(t *testing.T) {
    PrintHello()
}

func TestEncoder(t *testing.T) {
    e, err := NewEncoder(&Person{})
    if err != nil {
        t.Errorf(err.Error())
    }
    person := Person{
        ID: 1,
        Name: "Bobby Hill",
        Age: 13,
    }

    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
    })
    err = e.Store(client, person)
    if err != nil {
        t.Errorf(err.Error())
    }

    res, err := client.HGetAllMap("person:1").Result()
    if err != nil {
        t.Errorf(err.Error())
    }

    if len(res) != 2 {
        t.Errorf("Hash did not have expected size.")
    }

    if res["name"] != "Bobby Hill" {
        t.Errorf("Hash did not have expected contents.")
    }

    if res["age"] != "13" {
        t.Errorf("Hash did not have expected contents.")
    }
}
