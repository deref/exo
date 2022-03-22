#!/bin/bash

set +e

echo -n "extractgqlts"
time extractgqlts \
  --schema ./internal/resolvers/schema.gql \
  ./gui/src/**/*.svelte \
  > ./gui/src/lib/graphql/types.generated.ts

exit_code=$?
set -e

(
  cd gui
  echo "prettier"
  node ./node_modules/.bin/prettier \
    --write \
    ./src/lib/graphql/types.generated.ts
)

exit $exit_code
