#!/bin/bash

INVENTORY_FILE=hosts

SCRIPT_HOME="$(cd "$(dirname "$0")"; pwd)"

source common.sh

function deploy() {
  ansible-playbook ${VERBOSE} -i ${INVENTORY_FILE} ${SCRIPT_HOME}/ansible/index_mount_partition.yml
}

function clear() {
  ansible-playbook ${VERBOSE} -i ${INVENTORY_FILE} ${SCRIPT_HOME}/ansible/index_umount_partition.yml
}

function hosts() {
  list_hosts ${INVENTORY_FILE} index 8890 ':'
}

function start() {
  ansible-playbook ${VERBOSE} -i ${INVENTORY_FILE} ${SCRIPT_HOME}/ansible/index_start.yml
}

function stop() {
  ansible-playbook ${VERBOSE} -i ${INVENTORY_FILE} ${SCRIPT_HOME}/ansible/index_stop.yml
}

function usage() {
  local progname=$1
  echo -e "usage: ${progname} [OPTIONS] COMMAND"
  echo -e ""
  echo -e "Manage index servers"
  echo -e ""
  echo -e "Options:"
  echo -e "  -i INVENTORY, --inventory=INVENTORY"
  echo -e "                        specify inventory host path "
  echo -e "  -v, --verbose         verbose mode"
  echo -e "Commands:"
  echo -e "  clear     Clear index partitions"
  echo -e "  deploy    Deploy index partitions"
  echo -e "  hosts     List index servers host names"
  echo -e "  start     Start index servers"
  echo -e "  stop      Stop index servers"
  exit 1
}

POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
  case $1 in
    -i|--inventory)
      INVENTORY_FILE="$2"
      shift # past argument
      shift # past value
      ;;
    -v|--verbose)
      VERBOSE=-v
      shift # past argument
      ;;
    -*|--*)
      echo "Unknown option $1"
      usage 
      ;;
    *)
      POSITIONAL_ARGS+=("$1") # save positional arg
      shift # past argument
      ;;
  esac
done

$POSITIONAL_ARGS
