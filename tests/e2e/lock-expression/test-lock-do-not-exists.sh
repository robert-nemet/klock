#!/usr/bin/env sh

kubectl get lock -n $NAMESPACE lock-expression
if [ $? -eq 1 ]
then 
    exit 0
fi
exit 1