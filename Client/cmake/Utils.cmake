macro(GroupResources curdir pattern)
  file(GLOB _macro_children RELATIVE ${PROJECT_SOURCE_DIR}/${curdir} ${PROJECT_SOURCE_DIR}/${curdir}/${pattern})
  foreach(_loop_var ${_macro_children})
    if(IS_DIRECTORY ${PROJECT_SOURCE_DIR}/${curdir}/${_loop_var})
      GroupResources(${curdir}/${_loop_var} ${pattern})
    endif()
  endforeach()

  file(GLOB _macro_children RELATIVE ${PROJECT_SOURCE_DIR}/${curdir} ${PROJECT_SOURCE_DIR}/${curdir}/${pattern})
  foreach(_loop_var ${_macro_children})
    string(REPLACE "/" "\\" _groupname ${curdir})
    source_group(${_groupname} FILES ${PROJECT_SOURCE_DIR}/${curdir}/${_loop_var})
  endforeach()
endmacro()

macro(GroupTopResources pattern)
  file(GLOB _macro_children RELATIVE ${PROJECT_SOURCE_DIR}/${curdir} ${PROJECT_SOURCE_DIR}/${curdir}/${pattern})
  foreach(_loop_var ${_macro_children})
	source_group("" FILES ${PROJECT_SOURCE_DIR}/${_loop_var})
  endforeach()

  file(GLOB _macro_children RELATIVE ${PROJECT_SOURCE_DIR} ${PROJECT_SOURCE_DIR}/*)
  foreach(_loop_var ${_macro_children})
    if(IS_DIRECTORY ${PROJECT_SOURCE_DIR}/${_loop_var})
      GroupResources(${_loop_var} ${pattern})
    endif()
  endforeach()
endmacro()
