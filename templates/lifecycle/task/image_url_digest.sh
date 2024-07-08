#!/usr/bin/env bash
set -e
cat /layers/report.toml | grep "digest" | cut -d'"' -f2 | cut -d'"' -f2 | tr -d '\n' | tee $(results.APP_IMAGE_DIGEST.path)
cat $(results.APP_IMAGE_DIGEST.path) | tee "$(results.IMAGE_DIGEST.path)"

echo -n "$(params.APP_IMAGE)" | tee "$(results.IMAGE_URL.path)"