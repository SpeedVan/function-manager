#!/bin/bash

$GOPATH/src/k8s.io/code-generator/generate-groups.sh all \
github.com/SpeedVan/function-manager/k8s/group/client \
github.com/SpeedVan/function-manager/k8s/group/apis \
faas:v1