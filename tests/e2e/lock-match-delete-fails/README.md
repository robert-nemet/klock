# Test description

```yaml
when:
    lock:
        exists: true
    pod:
        exists: true
        match: true
        operation: delete
expect:
    result: delete fails
and:
    when:
        lock:
            exists: true
            operation: delete
        pod:
            exists: true
            match: true
            operation: delete
    expect:
        result: delete success

```
