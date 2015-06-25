cmake_minimum_required(VERSION 3.0)

project (Squirrel)

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${PROJECT_SOURCE_DIR}/../cmake")
set(_ProductPackage "com.dummy.dummy" CACHE STRING "Project package id")
include(Utils)

set (Third_ROOT ../3rd)
set (Third_BUILDTOP ${CMAKE_BINARY_DIR}-3rd)

add_subdirectory(${Third_ROOT} ${Third_BUILDTOP})

#################### include directory ####################
# include public 3rd libraries hdr
foreach (ph ${Public_3rd_HDRS})
  include_directories(${Third_ROOT}/${ph})
endforeach()

include_directories(${PROJECT_SOURCE_DIR})

#################### compile sources ####################
file(GLOB_RECURSE HDRS *.h)
file(GLOB_RECURSE SRCS *.m)
file(GLOB_RECURSE RESOURCES *.xib)
file(GLOB_RECURSE ICONS *.png)
file(GLOB_RECURSE STORYBOARD *.storyboard)

set(RESOURCES ${RESOURCES} ${STORYBOARD} ${ICONS})

add_executable(Squirrel MACOSX_BUNDLE ${HDRS} ${SRCS} ${RESOURCES})

#################### platform libraries ####################
set(PLATFORM_FRAMEWORKS
  CoreGraphics
  MobileCoreServices
  SystemConfiguration
  Security
  UIKit)

foreach(pf ${PLATFORM_FRAMEWORKS})
  set(PLATFORM_FRAMEWORK_FLAGS ${PLATFORM_FRAMEWORK_FLAGS} "-framework ${pf}")
endforeach()

target_link_libraries(Squirrel ${Public_3rd_LIBS} ${PLATFORM_FRAMEWORK_FLAGS})


#################### set source group ####################
GroupTopResources(*.h)
GroupTopResources(*.m)
GroupTopResources(*.xib)
GroupTopResources(*.storyboard)
#GroupTopResources(*.png)

#################### xcode property ####################
set_xcode_property(Squirrel GCC_GENERATE_DEBUGGING_SYMBOLS YES)
set_xcode_property(Squirrel CODE_SIGN_IDENTITY "iPhone Developer")
#set_xcode_property(Squirrel GCC_PREFIX_HEADER ${PREFIX_HEADER})
set_xcode_property(Squirrel GCC_PRCOMPILE_PREFIX_HEADER YES)
set_xcode_property(Squirrel DEBUG_INFORMATION_FORMAT "dwarf-with-dsym")
set_xcode_property(Squirrel SKIP_INSTALL "No")
set_xcode_property(Squirrel INSTALL_PATH "$(LOCAL_APPS_DIR)")
  
# plist property
set_xcode_property(Squirrel INFOPLIST_PREPROCESS YES)
set(PLIST_BUNDLE_GUI_IDENTIFIER "${_ProductPackage}")
set(PLIST_BUNDLE_PRODUCT_NAME "Squirrel")
set(PLIST_BUNDLE_EXECUTABLE "Squirrel")

set_target_properties(Squirrel
  PROPERTIES
  MACOSX_BUNDLE_INFO_PLIST ${CMAKE_CURRENT_SOURCE_DIR}/Info.plist.in
  )

# resources
set_source_files_properties(${ICONS} PROPERTIES MACOSX_PACKAGE_LOCATION Resources)
