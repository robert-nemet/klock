#!/usr/bin/env sh

kubectl delete pod -n $NAMESPACE hello-new-world
if [ $? -eq 1 ]
then 
    exit 0
fi

exit 1