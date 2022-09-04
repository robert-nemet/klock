# Deadlock

To avoid deadlock( when two locks lock each other ), ignore locks when it comes to `Lock` kind.

__Lock should not target lock.__

Ignore such cases. Deadlock detections is do not needed as result.
