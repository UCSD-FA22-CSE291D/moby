HOME=`pwd`
OUTPUT_FILE="${HOME}/predump_output.txt"
APP_NAME="test_app"
CP_PRE_NAME="cp_pre"
CP_NAME="cp"

# Before starting containers
> ${OUTPUT_FILE}


# Get container id
CONT_ID=`docker ps -qf "name=${APP_NAME}"`

echo -e "CONT_ID: ${CONT_ID}"
echo -e "CP_PRE_NAME: ${CP_PRE_NAME}"
echo -e "CP_NAME: ${CP_NAME}"

# remove original checkpoints
echo "Removing old checkpoint..."
docker checkpoint rm $CONT_ID $CP_PRE_NAME
docker checkpoint rm $CONT_ID $CP_NAME

# wait for removal
sleep 5

echo "Start Pre-dumping Timestamp:" >> ${OUTPUT_FILE}
date +"%T.%N" >> ${OUTPUT_FILE}

echo "Pre-dumping..."
docker checkpoint create $CONT_ID $CP_PRE_NAME --predump

echo "Start Dumping Timestamp:" >> ${OUTPUT_FILE}
date +"%T.%N" >> ${OUTPUT_FILE}

echo "Dumping..."
docker checkpoint create $CONT_ID $CP_NAME --parent-checkpoint=${CP_PRE_NAME}

echo "Start Restoring Timestamp:" >> ${OUTPUT_FILE}
date +"%T.%N" >> ${OUTPUT_FILE}

echo "Restoring..."
docker start --checkpoint $CP_NAME $CONT_ID

echo "End Timestamp:" >> ${OUTPUT_FILE}
date +"%T.%N" >> ${OUTPUT_FILE}

echo -e "Done.\n"
echo "output written to ${OUTPUT_FILE}"
