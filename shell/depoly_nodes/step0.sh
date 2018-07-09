
nodes=$1

sh kill.sh chain

rm -f nohup.out delos_*
rm -rf angine-*
rm -rf $path
cmd=`./chain/chain init`
echo $cmd

# mv $path_local /$ssd_disk/
# ln -s $path $path_local

IP=`ifconfig eth0|grep inet|awk '{print $2}'`
IPPORT=$IP:46656,
seeds=""
for node in ${nodes[@]}; do seeds=${node}:46656,${seeds}; done
seeds=${seeds/${IPPORT}/}
seeds=${seeds:0:-1}


rm -f ${path}/config.toml
echo "# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml
environment = \"production\"
node_laddr = \"tcp://${IP}:46656\"
rpc_laddr = \"tcp://${IP}:46657\"
moniker = \"anonymous\"
fast_sync = true
db_backend = \"leveldb\"
seeds = \"${seeds}\"
signbyCA = \"\"
" > ${path}/config.toml
