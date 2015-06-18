#!/bin/bash
mkdir build
cd build
cmake -G Xcode -DCMAKE_TOOLCHAIN_FILE=../cmake/iOS.cmake ../3rd
cd ..
