#!/bin/bash

#######################################
# Lists all host names that are part of a group 
# Arguments:
#   Path to inventory file
#   Group name
# Outputs:
#   Writes list of host names
#######################################
function list_hosts() {
  local inventory_file=$1
  local group_name=$2
  local host_port=$3
  local separator=$4
  awk -v TARGET=${group_name} -v HOST_PORT=${host_port} -v SEPARATOR="${separator}" -F ' *= *' '
    {
      if ($0 ~ /^\[.*\]$/) {
        gsub(/^\[|\]$/, "", $0)
        SECTION=$0
      } else if (SECTION==TARGET) {
        host_port=HOST_PORT
        where_host = match($0, /^([A-Za-z0-9\-]*)/, arr_host)
        where_ansible_host = match($0, /ansible_host=([A-Za-z0-9\-]*)/, arr_ansible_host)
        where_socket = match($0, /socket=([0-9]*)/, arr_socket)
        if (where_socket != 0)
          host_port += arr_socket[1]
        if (where_ansible_host != 0)
          host = arr_ansible_host[1]
        else
          host = arr_host[1]
          if (length(host) > 0)
            printf "%s%c%d\n", host, SEPARATOR, host_port
      }
    }
    ' ${inventory_file}
}
