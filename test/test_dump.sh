HOME=`pwd`
OUTPUT_FILE="${HOME}/dump_output.txt"
APP_NAME="test_app"
CP_NAME="cp"

# Before starting containers
> ${OUTPUT_FILE}

# Get container id
CONT_ID=`docker ps -qf "name=${APP_NAME}"`

echo -e "CONT_ID: ${CONT_ID}"
echo -e "CP_NAME: ${CP_NAME}"

# remove original checkpoint
docker checkpoint rm $CONT_ID $CP_NAME

sleep 2

echo "Start Dumping Timestamp:" >> ${OUTPUT_FILE}
date +"%T.%N" >> ${OUTPUT_FILE}

echo "Dumping..."
docker checkpoint create $CONT_ID $CP_NAME

echo "Start Restoring Timestamp:" >> ${OUTPUT_FILE}
date +"%T.%N" >> ${OUTPUT_FILE}

echo "Restoring..."
docker start --checkpoint $CP_NAME $CONT_ID

echo "End Timestamp:" >> ${OUTPUT_FILE}
date +"%T.%N" >> ${OUTPUT_FILE}

sleep 2

echo -e "Done.\n"
echo "output written to ${OUTPUT_FILE}"
