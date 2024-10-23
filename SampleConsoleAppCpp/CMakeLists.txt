cmake_minimum_required(VERSION 3.10)

# Set the project name
project(SampleConsoleAppCpp)

# Specify the C++ standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED True)

# Add the executable
add_executable(SampleConsoleAppCpp Program.cpp)

# Find and include libcurl
find_package(CURL REQUIRED)
include_directories(${CURL_INCLUDE_DIRS})
target_link_libraries(SampleConsoleAppCpp ${CURL_LIBRARIES})

# Find and include nlohmann/json
find_package(nlohmann_json REQUIRED)
target_link_libraries(SampleConsoleAppCpp nlohmann_json::nlohmann_json)
