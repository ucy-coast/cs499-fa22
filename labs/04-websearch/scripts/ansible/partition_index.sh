#!/bin/bash

#######################################
# Find all index parts in a directory.
# Arguments:
#   Directory that contains index parts
# Outputs:
#   Writes full path to the directory of each index part
#######################################
function find_parts() {
  local src_base_dir=$1
  echo $(ls ${src_base_dir} | sort | while read line; do echo ${line}; done)
}

#######################################
# Produce a list of index parts that correspond to a partition.
# Arguments:
#   String that contains a list of all index parts.
#   Partition identifier for which to produce the list of index parts
#   Number of index parts per partition  
# Outputs:
#   Writes a list of index parts.
#######################################
function partition_parts() {
  local parts=$1
  local partition_id=$2
  local partition_parts_count=$3
  local p=0
  local low_part_id=$((partition_id * partition_parts_count))
  local high_part_id=$(((partition_id+1) * partition_parts_count))
  for part in $parts; do
    if ((p >= low_part_id && p < high_part_id)); then
      echo $part
    fi
    (( p++))
  done
}

function usage() {
  local progname=$1
  echo -e "usage: ${progname} <parts_count> <partitions_count> <partition_id> <src_base_dir>"
  echo -e ""
  echo -e "Produce a list of index parts that correspond to a partition"
  echo -e ""
  echo -e "  <parts_count>\t\ttotal number of index parts to distribute equally among partitions"
  echo -e "  <partitions_count>\ttotal number of partitions"
  echo -e "  <partition_id>\tidentifier of the partition to produce list for"
  echo -e "  <src_base_dir>\tsource directory containing all the index parts"
  exit 1
}

function main() {
  if (( $# != 4 )); then
    usage
  fi

  local total_parts_count=$1
  local partitions_count=$2
  local partition_id=$3
  local src_base_dir=$4

  local all_parts=$(find_parts ${src_base_dir})
  local partitions_parts_count=$((total_parts_count / partitions_count))
  local parts=$(partition_parts "$all_parts" ${partition_id} ${partitions_parts_count})
  for part in $parts; do
    echo $part
  done
}

main "$@"
