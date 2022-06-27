# Test description

```yaml
when:
    lock:
        exists: false
    pod:
        exists: true
        match: false
    operation: delete
expect:
    result: success
```
