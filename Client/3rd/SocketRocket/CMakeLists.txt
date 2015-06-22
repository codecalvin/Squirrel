cmake_minimum_required(VERSION 3.0)

project (SocketRocket)

# subdirectory placed in 3rd
set(TOP_3rd SocketRocket)

set(SocketRocketDIR SocketRocket/SocketRocket)

file(GLOB_RECURSE SR_HDRS ${SocketRocketDIR} ${SocketRocketDIR}/*.h)
file(GLOB_RECURSE SR_SRCS ${SocketRocketDIR} ${SocketRocketDIR}/*.m)

add_library(SocketRocket ${SR_HDRS} ${SR_SRCS})

set_property (TARGET SocketRocket APPEND_STRING PROPERTY 
  COMPILE_FLAGS "-fobjc-arc")

#################### mark as advance ####################
list(APPEND Public_3rd_HDRS ${TOP_3rd}/${SocketRocketDIR})
set(Public_3rd_HDRS ${Public_3rd_HDRS} PARENT_SCOPE)

list(APPEND Public_3rd_LIBS SocketRocket)
set(Public_3rd_LIBS ${Public_3rd_LIBS} PARENT_SCOPE)
