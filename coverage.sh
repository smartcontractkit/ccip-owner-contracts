#!/bin/sh

forge coverage --ffi --report lcov && genhtml lcov.info --branch-coverage --legend --output-dir coverage && open coverage/index.html
