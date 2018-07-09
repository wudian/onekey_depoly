sh commit.sh

address=865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5
passwd=h0pe
path=~/.ann_runtime/
midapi_node=10.253.160.117 #管控
nodes=(10.253.112.132 10.253.164.121 10.253.115.1 10.253.168.160)
bin=~/advertise-bin

cd depoly_nodes

for node in ${nodes[@]}; do sh step1.sh $node $nodes $path $passwd; done

sh step2.sh $nodes

for node in ${nodes[@]}; do sh step3.sh $node; done

sh step4.sh $midapi_node
