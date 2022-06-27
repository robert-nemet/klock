#!/usr/bin/env sh

kubectl delete pod -n $NAMESPACE hello-world
if [ $? -eq 1 ]
then 
    exit 0
fi

exit 1