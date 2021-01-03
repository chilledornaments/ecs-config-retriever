#!/usr/bin/env sh

sleep 10
# TEST_FILE_PATH is an env var
file_contents=$(cat ${TEST_FILE_PATH})

# EXPECTED_FILE_CONTENTS is an env var
if [ "${file_contents}" != "${EXPECTED_FILE_CONTENTS}" ]
then
  echo "Expected file contents to be ${EXPECTED_FILE_CONTENTS} but received ${file_contents}"
  echo ""
  exit 1
else
  echo "File contents matched what was expected"
  echo ""
fi

