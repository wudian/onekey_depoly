
node=$1
nodes=$2
path=$3
passwd=$4

if [ ! -d $node ]; then
mkdir $node
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
send "cd /home/wd/advertise-bin; git checkout -- .; git pull; sh step0.sh $nodes;  \r"
expect "*#"
send "exit\r"



spawn scp wd@$node:${path}priv_validator.json $node/
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

EOF
