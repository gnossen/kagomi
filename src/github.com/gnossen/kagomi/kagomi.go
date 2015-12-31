package kagomi

import (
    "fmt"
    "reflect"
    "strings"
    "strconv"
    "gopkg.in/redis.v3"
)

func PrintHello() {
    fmt.Printf("カゴミへは、ようこそ！\n")
}

type Encoder struct {
    Hash        map[int]string
    PrimaryKey  int
    KeyPrefix   string  
    NumFields   int
}

func NewEncoder(m interface{}) (Encoder, error) {

    typ := reflect.TypeOf(m)
    if typ.Kind() == reflect.Ptr {
        typ = typ.Elem()
    } 

    e := Encoder{
        Hash: make(map[int]string),
        PrimaryKey: -1,
        KeyPrefix: "",
        NumFields: typ.NumField(),
    }

    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        if prefix := field.Tag.Get("redis-primary"); prefix != "" {
            e.PrimaryKey = i
            e.KeyPrefix = prefix
        } else if key := field.Tag.Get("redis"); key != "" {
            e.Hash[i] = key
        } else {
            e.Hash[i] = field.Name
        }
    }

    return e, nil
}

func valueToString(v reflect.Value) string {
    k := v.Kind()
    if k >= reflect.Int && k <= reflect.Int64 {
        return strconv.FormatInt(v.Int(), 10)
    } else if k >= reflect.Uint && k <= reflect.Uint64 {
        return strconv.FormatUint(v.Uint(), 10)
    } else {
        return v.String()
    }
}

func (e *Encoder) makeKey(prefix, name string) string {
    return fmt.Sprintf("%s:%s", prefix, name)
}

func (e *Encoder) getPrefix(v reflect.Value) string {
    primaryStr := valueToString(v.Field(e.PrimaryKey))
    return strings.Replace(e.KeyPrefix, "%P", primaryStr, -1)    
}

func (e *Encoder) genMainHash(val reflect.Value)  map[string]string {
    hash := make(map[string]string)
    for i, name := range e.Hash {
       hash[name] = valueToString(val.Field(i))
    }
    return hash
}

func (e *Encoder) Store(c *redis.Client, m interface{}) error {
    val := reflect.ValueOf(m)
    prefix := e.getPrefix(val)
    hash := e.genMainHash(val)
    err := writeHash(c, prefix, hash) 
    return err
}

func writeHash(c *redis.Client, prefix string, hash map[string]string) error {
    for key, value := range hash {
        err := c.HSet(prefix, key, value).Err()
        if err != nil {
            return err
        }
    }
    return nil
}
