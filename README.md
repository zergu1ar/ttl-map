# TTL Map

```go
expireMap := ttlmap.New()
expireMap.Add("foo", 2018, time.Millisecond)
if !expireMap.Exists("foo") {
    log.Fatal("entry not exists")
}

time.Sleep(time.Millisecond)

if expireMap.Exists("foo") {
    log.Fatal("entry not exists")
}
```