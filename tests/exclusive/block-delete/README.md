# Test Exclusive feat

The ServiceAccount `johny` is assigned Role `terminator` which allows it to delete pods.
There is a Pod with label `aura:red` which is protected by the `Lock`.

When try to delete the proteced Pod, it is blocked by the `Lock`.

When the `Lock` is updated with exception to match SA `johny` deletion is successful.
