from=$midapi_node
node=$1   #ip

if [ -d $node ]; then
rm -rf $node
fi

/usr/bin/expect <<-EOF
set timeout 30

spawn ssh wd@$node
expect {
"*yes/no" { send "yes\r"; exp_continue }
"*password:" { send "$passwd\r" }
}
expect "*#"
send "\r"
expect "*#"
send "scp wd@$from:${bin}/genesis.json $path  \r"
expect {
 "(yes/no)?"
  {
  send "yes\n"
  expect "*assword:" { send "$passwd\n"}
 }
 "*assword:"
{
 send "$passwd\n"
}
}
expect "100%"
expect eof

send "\r"
expect "*#"
send "cd ${bin}/chain; nohup ./chain node  > nohup.out &  \r"
expect "*#"

send "exit\r"




EOF

echo "==================$node==========chain已经启动============="
