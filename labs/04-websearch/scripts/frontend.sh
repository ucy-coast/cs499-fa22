#!/bin/bash

INVENTORY_FILE=hosts

SCRIPT_HOME="$(cd "$(dirname "$0")"; pwd)"

source common.sh

function config() {
  index_servers=$(list_hosts ${INVENTORY_FILE} index 8890 ' ' | awk '{ printf "%c%s%c\n", 0x27, $0, 0x27 }' | paste -sd "," -)
  ansible-playbook ${VERBOSE} -i ${INVENTORY_FILE} -e "{'index_servers':\"[${index_servers}]\"}" ${SCRIPT_HOME}/ansible/frontend_config.yml
}

function hosts() {
  list_hosts ${INVENTORY_FILE} frontend 8080 ':'
}

function start() {
  ansible-playbook ${VERBOSE} -i ${INVENTORY_FILE} ${SCRIPT_HOME}/ansible/frontend_start.yml
}

function stop() {
  ansible-playbook ${VERBOSE} -i ${INVENTORY_FILE} ${SCRIPT_HOME}/ansible/frontend_stop.yml
}

function test() {
  for frontend in $(list_hosts ${INVENTORY_FILE} frontend 8080 ':'); do
    echo "Test frontend server instance ${frontend}..."
    curl "${frontend}/onlyHits.jsp?query=google"
  done  
}

function usage() {
  local progname=$1
  echo -e "usage: ${progname} [OPTIONS] COMMAND"
  echo -e ""
  echo -e "Manage frontend servers"
  echo -e ""
  echo -e "Options:"
  echo -e "  -i INVENTORY, --inventory=INVENTORY"
  echo -e "                        specify inventory host path "
  echo -e "  -v, --verbose         verbose mode"
  echo -e "Commands:"
  echo -e "  config    Configure frontend servers"
  echo -e "  hosts     List frontend servers host names"
  echo -e "  start     Start frontend servers"
  echo -e "  stop      Stop frontend servers"
  echo -e "  test      Test frontend servers"
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