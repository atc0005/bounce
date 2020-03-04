#!/usr/bin/env python3

# Copyright 2020 Adam Chalkley
#
# https://github.com/atc0005/bounce
#
# Licensed under the MIT License. See LICENSE file in the project root for
# full license information.

# https://github.com/getify/JSON.minify/tree/python
#
# Install the JSON-minify package by running:
# "pip install JSON-minify --user"

from json_minify import json_minify
import sys

if len(sys.argv) < 2:
    sys.exit("Please provide the name of a JSON file to minify")

file = sys.argv[1]

try:
    fh = open(file, "r")
except:
    sys.exit("Unable to open %s JSON file" % file)

content = fh.read()
fh.close()

try:
    minified_content = json_minify(content)
except:
    sys.exit("Failed to minify %s" % file)

# Send to stdout so user can choose where the content goes
print(minified_content)
