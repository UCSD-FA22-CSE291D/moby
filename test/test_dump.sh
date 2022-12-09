HOME=`pwd`
OUTPUT_FILE="${HOME}/output.txt"
SERVER_NAME="server_app"
CP_NAME="cp"

# Before starting containers
# rm -rf $HOME/tmp
# mkdir -p $HOME/tmp

# Get container id of the server by name
CONT_ID=`docker ps -qf "name=${SERVER_NAME}"`

echo -e "CONT_ID: ${CONT_ID}"
echo -e "CP_NAME: ${CP_NAME}"

# remove original checkpoint
docker checkpoint rm $CONT_ID $CP_NAME

# Start the request
echo "Start curl request..."

curl -s localhost:8000 > ${OUTPUT_FILE} &

sleep 2


# Dump the server
echo "Dumping..."
docker checkpoint create $CONT_ID $CP_NAME

# Restore the server with the same root FS
echo "Restoring..."
docker start --checkpoint $CP_NAME $CONT_ID

# Wait for curl to finish
echo "Waiting for curl to finish..."
sleep 15

echo -e "Done.\n"
echo "curl output written to ${OUTPUT_FILE}"
