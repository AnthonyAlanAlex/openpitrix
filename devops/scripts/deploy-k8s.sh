#!/bin/bash

PILOT="on"
[ "$1" == "-pilot=off" ] && PILOT="off"
echo "pilot is $PILOT"

# Back to the root of the project
cd $(dirname $0)
cd ../..

kubectl create namespace openpitrix-system
kubectl create secret generic mysql-pass --from-file=./devops/kubernetes/password.txt -n openpitrix-system
kubectl apply -f ./devops/kubernetes/db
kubectl apply -f ./devops/kubernetes/etcd
kubectl apply -f ./devops/kubernetes/openpitrix
[ "$PILOT" == "on" ] && kubectl apply -f ./devops/kubernetes/pilot
