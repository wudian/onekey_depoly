#sh 

sh shell/kill.sh chain
sh shell/kill.sh api

path=/home/wd/.ann_runtime
rm -rf $path

cd shell/chain
cmd=`./chain init`
echo $cmd

IP=`ifconfig eth0|grep inet|awk '{print $2}'`


rm -f ${path}/config.toml
echo "# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml
environment = \"production\"
node_laddr = \"tcp://${IP}:46656\"
rpc_laddr = \"tcp://${IP}:46657\"
moniker = \"anonymous\"
fast_sync = true
db_backend = \"leveldb\"
seeds = \"\"
signbyCA = \"\"
" > ${path}/config.toml




main_alloc=`sed -n '17,17p' /home/alloc.md|sed s/[[:space:]\t\"]//g` #主链预留账户
pub_key=`sed -n '14,14p' ${path}/priv_validator.json|sed s/[[:space:]\t\"]//g`

rm -f ${path}/genesis.json
echo "{
        \"app_hash\": \"\",
        \"chain_id\": \"delos\",
        \"genesis_time\": \"0001-01-01T00:00:00.000Z\",
        \"init_accounts\": [
                {
                        \"address\": \"${main_alloc}\",
                        \"startingbalance\": \"100000000000000\"
                }
        ],
        \"plugins\": \"specialop\",
        \"validators\": [
                {
                        \"amount\": 100,
                        \"is_ca\": true,
                        \"name\": \"\",
                        \"pub_key\": [
                                1,
                                \"${pub_key}\"
                        ],
                        \"rpc\": \"tcp://${IP}:46657\"
                }
        ]
}
" > ${path}/genesis.json

cmd=`nohup ./chain node > nohup.out &`
echo $cmd

cd ../midapi
cmd=`nohup ./midapi > nohup.out &`
echo $cmd

# cd ../midware
# cmd=`nohup ./midware > nohup.out &`
# echo $cmd
