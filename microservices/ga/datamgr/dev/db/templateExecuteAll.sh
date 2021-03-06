{{define "dev/db/templateExecuteAll.sh"}}
#!/bin/bash

SQL="cockroach sql"
PORT="26257"
HOST="127.0.0.1"
USER="root"
CMD=`echo $SQL "--insecure" --host=$HOST:$PORT`

buildCommand()
{
  CMD=`echo $SQL "--insecure" --host=$HOST:$PORT`
}

all()
{

  echo "========================================"
  echo " Initializing tables"
  start=`date`
  echo " Starting at : $start"
  echo " Using: $CMD"
  echo ""

# 1 Create {{.Organization}}Admin if it doesn not already exists
  $CMD < {{.Organization}}Admin.sql

# 2 Create {{.Organization}}Web db
  $CMD < {{.Organization}}.sql

# 3 Create A{{.Organization}}dmin all on kevlarWeb db
  $CMD < {{.Organization}}GrantAdmin.sql

# 4 Create {{.Organization}}Table 
  $CMD < {{.Name}}CreateTable.sql
}

usage()
{
  echo "usage: {{.Organization}}ExecuteAll.sh -k |--k8s"
  echo "    Created database and users as needed"
  echo "    -k locates and posts to local k8s cluster"
  echo "    it will default to $host on port $port"
}

## Main
while [ "$1" != "" ]; do
  case $1 in
    -k | --k8s ) shift

      if [ -z "$COCKROACH_HOST" ]
      then
        echo "COCKROACH_HOST must be set"
        exit
      else
        HOST=$COCKROACH_HOST
      fi

      if [ -v "$COCKROACH_PORT" ]
      then
        PORT=$COCKROACH_PORT
      fi
      buildCommand
      ;;
  -h | --help ) usage
    exit
    ;;
  * ) shift
    ;;
  esac
done

# call all
all
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
