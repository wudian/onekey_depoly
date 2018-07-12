#sh killps.sh  程序名(cmd  chain 等)


if [ "$1" = "" ]; then 
	exit
fi
echo "kill $1"

kill -9 $(ps -ef|grep $1|gawk '$0 !~/grep/ {print $2}' |tr -s '\n' ' ')