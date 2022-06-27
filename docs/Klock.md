# Klock

Klock prevents updates or/and deletion of Kubernetes resources.

## Design

Idea is to make Lock resource which specifies which resources are locked, and how they are locked. For deletion or/and update.
Matching is done by matching labels. 