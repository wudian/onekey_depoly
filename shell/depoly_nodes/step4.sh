
node=$1

/usr/bin/expect <<-EOF
set timeout 30

spawn ssh wd@$node
expect {
"*yes/no" { send "yes\r"; exp_continue }
"*password:" { send "$passwd\r" }
}
expect "*#"
send "cd ${bin}; sh kill.sh midapi; sh step5.sh $node ;  \r"
expect "*#"
send "cd midapi/; nohup ./midapi > nohup.out &  \r"
expect "*#"
send "sleep 1 \r"
expect "*#"
send "exit\r"


EOF

echo "==================$node==========midapi已经启动============="
