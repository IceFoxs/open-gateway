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
if [ ! -d "${BASE_DIR}/logs" ]; then
  mkdir ${BASE_DIR}/logs
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