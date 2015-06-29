#!/bin/bash
if [ -e build ];
then
	rm -r build
fi

if [ -e build-3rd ];
then
	rm -r build-3rd
fi

mkdir build
cd build
cmake -G Xcode -DCMAKE_TOOLCHAIN_FILE=../cmake/iOS.cmake ../Squirrel
cd ..
