
IP=$midapi_node

rm -f midapi/config.json
echo "
{
  \"ListenAddress\": \":8000\",
  \"BackendCallAddress\": \"tcp://${IP}:46657\",
  \"Timeout\":30,
  \"Public\":false,
  \"GasLimit\": 10000000,
  \"Debug\":true,
  \"TiConnEndpoint\":\"http://11382-zis-other-ti-capsule-anlinkapiserver.test.za-tech.net\",
  \"TiConnKey\":\"ZWViMTYyNWJlMTNmNDg5NDg2MTA1Mzhl\",
  \"TiConnSecret\":\"MjJkZWMzMDUyMjM5NDc3YTkxYzZlNjkxMjA5NGQ2YTQ4Mzc3YzlmMmQ1ZDc0MDdj\",

  \"redis\":\"r-tj7748e7ace28aa4.redis.rds.aliyuncs.com:6379\",
  \"redisidls\":10,
  \"redistimeout\":60,
  \"expire\":36000
}
" > midapi/config.json
