cmake_minimum_required(VERSION 3.0)

project (AFNetworking)

# subdirectory placed in 3rd
set(TOP_3rd AFNetworking)

set(AFNetworkingDIR AFNetworking/AFNetworking)

file(GLOB_RECURSE AFN_HDRS ${AFNetworkingDIR} ${AFNetworkingDIR}/*.h)
file(GLOB_RECURSE AFN_SRCS ${AFNetworkingDIR} ${AFNetworkingDIR}/*.m)

add_library(AFNetworking ${AFN_HDRS} ${AFN_SRCS})

set_property (TARGET AFNetworking APPEND_STRING PROPERTY 
  COMPILE_FLAGS "-fobjc-arc")

#################### mark as advance ####################
list(APPEND Public_3rd_HDRS ${TOP_3rd}/${AFNetworkingDIR})
set(Public_3rd_HDRS ${Public_3rd_HDRS} PARENT_SCOPE)

list(APPEND Public_3rd_LIBS AFNetworking)
set(Public_3rd_LIBS ${Public_3rd_LIBS} PARENT_SCOPE)
