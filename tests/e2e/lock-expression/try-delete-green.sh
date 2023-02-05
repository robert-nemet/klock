#!/usr/bin/env sh

kubectl delete pod -n $NAMESPACE green
if [ $? -eq 1 ]
then 
    exit 0
fi

exit 1