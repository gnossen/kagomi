package kagomi

import (
    "time"
)

type Post struct {
    ID          uint64
    Image       []byte
    Thumbnail   []byte
    Location    Geotag
    PosterID    uint64
    Description string
    Hearts      uint32
    Flags       uint32
    Time        time.Time
    Comments    []Comment
}

type Comment struct {
    ID          uint64
    PosterID    uint64
    ParentID    uint64
    Text        string
    Hearts      uint32
    Flags       uint32
    Time        time.Time
}

type Person struct {
    ID          uint64 `redis-primary:"person:%P"`
    Name        string `redis:"name"`
    Age         uint8  `redis:"age"`
}

type Place struct {
    ID          uint64 `redis-primary:"place:%P"`
    Name        string `redis:"name"` 
    Location    Geotag `redis-embed:"loc-%s"`
}

type Geotag struct {
    Latitude    float64 `redis:"lat"`
    Longitude   float64 `redis:"long"`
}
