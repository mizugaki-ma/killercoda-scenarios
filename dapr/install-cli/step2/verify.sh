#!/bin/bash

# Verify the Dapr Runtime is installed
dapr --version | grep Runtime | awk '{print $3}' | grep -E '*.*.*' > /dev/null 2>&1