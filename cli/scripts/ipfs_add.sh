echo "+ $1"
MFSBASEPATH="/disc3/D"
DROPID="20220401"
CID=$(ipfs add -r -Q --chunker=size-4096 --offline $1)
echo "> $CID"
ipfs files cp -p /ipfs/"$CID" $MFSBASEPATH/$DROPID/$1
