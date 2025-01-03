SCRIPT_DIR=$(dirname ${BASH_SOURCE})

IMAGE_NAME=$(cat ${SCRIPT_DIR}/IMAGE_NAME)
IMAGE_VERSION=$(cat ${SCRIPT_DIR}/IMAGE_VERSION)

CONTAINER_NAME=${IMAGE_NAME}
CONTAINER_ID=
CONTAINER_ID_FILE=$(pwd)/CONTAINER_ID

if [[ -f ${CONTAINER_ID_FILE} ]]; then
	CONTAINER_ID=$(cat ${CONTAINER_ID_FILE})
fi

build() {
	docker build -t ${IMAGE_NAME}:${IMAGE_VERSION} . && docker image prune -f
}

run() {
	if [[ -n ${CONTAINER_ID} ]]; then
		docker start ${CONTAINER_ID}
	else
		WORK_DIR=/go/src/neo-inu
		docker run --name ${CONTAINER_NAME} --detach --env-file .env \
			--volume ./internal:${WORK_DIR}/internal \
			--volume ./pkg:${WORK_DIR}/pkg \
			--volume ./cmd:${WORK_DIR}/cmd ${IMAGE_NAME}:${IMAGE_VERSION} >>${CONTAINER_ID_FILE}
	fi
}

restart() {
	docker restart ${CONTAINER_ID}
}

stop() {
	docker stop ${CONTAINER_ID}
	return 0
}

clean_container() {
	docker container rm ${CONTAINER_ID} && rm ${CONTAINER_ID_FILE}
	return 0
}

clean_image() {
	docker image rm ${IMAGE_NAME}:${IMAGE_VERSION}
	return 0
}

case $1 in
build | run | stop | restart | clean_container | clean_image)
	$1
	;;
"")
	echo No argument provided
	exit 1
	;;
*)
	echo No option matched
	exit 1
	;;
esac
