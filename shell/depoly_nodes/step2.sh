
nodes=$1
address=$2
bin=$3

str0="
{
        \"app_hash\": \"\",
        \"chain_id\": \"annchain-CaJIFw\",
        \"genesis_time\": \"0001-01-01T00:00:00.000Z\",
        \"init_accounts\": [
                {
                        \"address\": \"$address\",
                        \"startingbalance\": \"100000000000000000000000000000000\"
                }
        ],
        \"plugins\": \"specialop\",
        \"validators\": [
"

str1=""
for node in ${nodes[@]}; 
do 
	pub_key=`sed -n '14,14p' ${node}/priv_validator.json|sed s/[[:space:]\t\"]//g`;
	str1="
	{
                        \"amount\": 100,
                        \"is_ca\": true,
                        \"name\": \"\",
                        \"pub_key\": [
                                1,
                                \"${pub_key}\"
                        ],
                        \"rpc\": \"tcp://${node}:46657\"
                },
	"$str1
done
str1=${str1:0:-1}

str2="
	]
}
"

echo ${str0}${str1}${str2} > ${bin}/genesis.json


