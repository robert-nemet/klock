# Things to Do

## Docs

* [X] How to test
* [ ] How to start
* [ ] Explain example in README.md
* [ ] Explain how to configure:
  * [ ] Kustomize
  * [ ] Helm

## CI

* [ ] Setup image build and push
* [ ] Setup run of integration tests

## Testing

* [ ] Setup requirements:
  * [ ] test webhook config
  * [ ] test clusterrole config
* [ ] Setup KUTTL tests: 
  * [X] Delete pod when there is no any lock => true
  * [X] Delete pod when there is a lock and pod matches the lock => false
  * [X] Delete pod when there is a lock and pod does not match the lock => true
  * [X] Delete pod when there is a lock, two pods, no pods match
  * [X] Delete pod when there is a lock, two pods, one pod match, one pod not
  * [ ] Two locks present:
    * [ ] Delete pod when there is a lock and pod matches the lock => false
    * [ ] Delete pod when there is a lock and pod does not match the lock => true
    * [ ] Delete pod when there is a lock, two pods, no pods match
    * [ ] Delete pod when there is a lock, two pods, one pod match, one pod not
  * [ ] Update pod when no lock
  * [ ] Update pod when there is a lock
  * [ ] Update pod when there is a lock and pod matches the lock
  * [ ] Update pod when there is a lock and pod does not match the lock

## Improvements

* Setup configuration:
  * Log level
* Metrics (load test)

## Bugs
