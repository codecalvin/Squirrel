cmake_minimum_required(VERSION 3.0)

project (AFNetworking)

set(AFNetworkingDIR AFNetworking/AFNetworking)

file(GLOB_RECURSE AFN_HDRS ${AFNetworkingDIR} ${AFNetworkingDIR}/*.h)
file(GLOB_RECURSE AFN_SRCS ${AFNetworkingDIR} ${AFNetworkingDIR}/*.m)

add_library(AFNetworking ${AFN_HDRS} ${AFN_SRCS})

set_property (TARGET AFNetworking APPEND_STRING PROPERTY 
  COMPILE_FLAGS "-fobjc-arc")
