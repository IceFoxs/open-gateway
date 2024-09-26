export GO_ENV=test
export BASE_DIR=$(
  cd $(dirname $0)
  pwd
)
export OUTLOG=""
export OUT_LOG_DIR=""
while getopts "o" opt; do
  case $opt in
  o)
    OUTLOG="out.log"
    ;;
  esac
done
export SERVER="opengateway"
export CONF_PATH=${BASE_DIR}/config/conf.yaml
if [ ! -d "${BASE_DIR}/logs" ]; then
  mkdir ${BASE_DIR}/logs
fi
if [ ! -d "${BASE_DIR}/logs/nacos/cache" ]; then
  mkdir -p ${BASE_DIR}/logs/nacos/cache
fi
if [ ! -d "${BASE_DIR}/logs/nacos/log" ]; then
  mkdir -p ${BASE_DIR}/logs/nacos/log
fi
if [ -n "${OUTLOG}" ]; then
    OUT_LOG_DIR=${BASE_DIR}/${OUTLOG}
else
    OUT_LOG_DIR=/dev/null
fi
export APP_IDENTITY="opengateway.opengateway"
chmod u+x ${BASE_DIR}/${SERVER}
echo "${BASE_DIR}/${SERVER} ${APP_IDENTITY}"
nohup ${BASE_DIR}/${SERVER} ${APP_IDENTITY} >> $OUT_LOG_DIR 2>&1 &
sleep 1s
rm -Rf ${BASE_DIR}/log