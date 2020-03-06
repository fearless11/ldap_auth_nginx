#!/bin/bash
# date: 2020/03/05
# desc: 启停

DIR="/data/project/ldap-auth-nginx"
APP="ldap_auth_nginx"
CFG="${DIR}/conf/app.conf"

chmod +x ${DIR}/${APP}

function stop(){
  pgrep -f ${DIR}/${APP} | xargs kill
}

function start(){
  num=`pgrep -f ${DIR}/${APP} | wc -l`
  if [[ $num < 1 ]];then
    ${DIR}/${APP}  &>  ${DIR}/log &
    echo "start ok ..."
  fi
}


case $1 in
  start)
     start;;
  stop)
     stop;;
esac

