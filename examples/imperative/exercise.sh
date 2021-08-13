#!/bin/bash

# It's not currently sensible to actually run this script.
exit 1

set -e

exo new process tick ./bin/tick

exo new container echo -p 2222:80 ealen/echo-server:0.5.1

exo ls
