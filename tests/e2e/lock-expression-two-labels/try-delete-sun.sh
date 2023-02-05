#!/usr/bin/env sh

kubectl delete pod -n $NAMESPACE sun
if [ $? -eq 1 ]
then 
    exit 0
fi

exit 1