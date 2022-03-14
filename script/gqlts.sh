#!/bin/bash

set +e

extractgqlts \
  --schema ./internal/resolvers/schema.gql \
  ./gui/src/**/*.svelte \
  > ./gui/src/lib/graphql/types.generated.ts

exit_code=$?
set -e

(
  cd gui
  node ./node_modules/.bin/prettier \
    --write \
    ./src/lib/graphql/types.generated.ts
)

exit $exit_code
