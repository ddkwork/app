cmake -DCMAKE_BUILD_TYPE=Debug "-DCMAKE_MAKE_PROGRAM=ninja.exe" -G Ninja -S . -B cmake-build-debug
cmake --build cmake-build-debug --target demo -j 6