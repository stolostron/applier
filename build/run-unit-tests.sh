#!/bin/bash -ex
# Copyright Contributors to the Open Cluster Management project

_script_dir=$(dirname "$0")
mkdir -p test/unit/coverage

echo 'mode: atomic' > test/unit/coverage/coverage.out
echo '' > test/unit/coverage/coverage.tmp
echo -e "${GOPACKAGES// /\\n}" | xargs -n1 -I{} $_script_dir/test-package.sh {} ${GOPACKAGES// /,}

echo "Calculate coverage"
if [[ ! -f "test/unit/coverage/coverage.out" ]]; then
    echo "Coverage file test/unit/coverage/coverage.out does not exist"
    exit 0
fi

COVERAGE=$(go tool cover -func=test/unit/coverage/coverage.out | grep "total:" | awk '{ print $3 }' | sed 's/[][()><%]/ /g')
echo "-------------------------------------------------------------------------"
echo "TOTAL COVERAGE IS ${COVERAGE}%"
echo "-------------------------------------------------------------------------"

go tool cover -html=test/unit/coverage/coverage.out -o=test/unit/coverage/coverage.html
echo "test/unit/coverage/coverage.html generated"
