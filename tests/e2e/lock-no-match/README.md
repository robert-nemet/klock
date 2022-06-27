# Test description

```yaml
when:
    lock:
        exists: true
    pod:
        exists: true
        match: false
    operation: delete
expect:
    result: success
```
