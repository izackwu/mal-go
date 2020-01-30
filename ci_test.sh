#!/bin/bash

BINARY_FILE=./mal-go
TESTS_FOLDER=./tests
TEST_RUNNER=./runtest.py

# echo
echo "BINARY_FILE   $BINARY_FILE"
echo "TESTS_FOLDER  $TESTS_FOLDER"
echo "TEST_RUNNER   $TEST_RUNNER"

# make sure the  binary file has been correctly generated
if [[ ! -f $BINARY_FILE ]]; then
    echo "$BINARY_FILE doesn't exist!"
    exit 1
fi

# run all tests in the folder
FAIL=0
for file in $TESTS_FOLDER/*
do
    if [[ -f $file ]]; then
        echo "Run test cases in $file"
        $TEST_RUNNER $file $BINARY_FILE
        if [[ ! 0 -eq $? ]]; then
            let FAIL=1
        fi
    fi
done
exit $FAIL


