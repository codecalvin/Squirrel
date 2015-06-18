#!/bin/bash
mkdir build
cd build
#cmake -G Xcode -DCMAKE_TOOLCHAIN_FILE=../cmake/iOS.cmake ../3rd
cmake -G Xcode -DCMAKE_TOOLCHAIN_FILE=../cmake/iOS.cmake ../Squirrel
cd ..
