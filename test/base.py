# -* - coding: UTF-8 -* -  

import urllib2
import json
import httplib, urllib
import copy
import time

ip = "127.0.0.1"
# ip = "10.253.7.159"
# ip = "10.253.163.57"
port = 8000

init_account_pri = "a8971729fbc199fb3459529cebcd8704791fc699d88ac89284f23ff8e7fca7d6"
init_account_pub = "865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5"



param_tpl = {
        "nonce":"0",
        # "createtime":"1527824140",
        "basefee": "0",
        "memo": "xxx",
        "lowertime": 0,
        "uppertime": 0,
        "privkey": [],
        "operations": []
}

# 一.
def genkey():  
    try:
        httpClient = httplib.HTTPConnection(ip, port, timeout=300)
        httpClient.request('GET', '/v1/genkey')
        response = httpClient.getresponse()
        text = response.read()
        if 200 == response.status:
            ret = json.loads(text)
            if ret['isSuccess']:
                return ret['result']['address'], ret['result']['privkey']
        else:
            print response.reason, text
        return "", ""
    except Exception, e:
            print e
    finally:
        if httpClient:
            httpClient.close()


def nonce(pub):  
    try:
        httpClient = httplib.HTTPConnection(ip, port, timeout=300)
        httpClient.request('GET', '/v1/nonce/'+pub)
        response = httpClient.getresponse()
        text = response.read()
        if 200 == response.status:
            ret = json.loads(text)
            if ret['isSuccess']:
                return ret['result']
        return -1
    except Exception, e:
            print e
    finally:
        if httpClient:
            httpClient.close()

# 二十五.
# 成功的话，返回交易哈希、交易结果
# 失败的话，返回错误信息
def transactions(pub, pri, op, isbigdata=False):  
    def tx(param):
        try:
            headers = {"Content-type": "application/x-www-form-urlencoded", "Accept": "text/plain"}
            httpClient = httplib.HTTPConnection(ip, port, timeout=300)
            httpClient.request("POST", "/v1/transactions", json.dumps(param), headers)
            response = httpClient.getresponse()
            text = response.read()
            if 200 == response.status:
                ret = json.loads(text)
                if ret['isSuccess']:
                    if isbigdata:
                        return True, ret['result']['bigDatas'][0]['hash']
                    else:
                        return True, ret['result']['tx']
                else:
                    return False, ret['message']
            return False, response.reason
        except Exception, e:
                return False, repr(e)
        finally:
            if httpClient:
                httpClient.close()
    param = copy.deepcopy(param_tpl)
    nc = nonce(pub)
    if nc == -1:
        return False, "nonce == -1"
    param['nonce'] = str(nc)
    param['privkey'].append(pri)
    param['operations'].append(op)
    return tx(param)

############################################################################

# 二十八.
def create_account(pub):
    op = {
        "optype": "create_account",
        "source": init_account_pub,
        "destination": pub,
        "startingBalance":"2000000000"
    }

    b, ret = transactions(init_account_pub, init_account_pri, op)
    time.sleep(1)
    return b, ret

# 修改账号
def manage_data(pub, pri):
    op = {
        "optype": "manage_data",
        "source": pub,
        "keyPair":[
            {
                "name":"name",
                "value":"wudian"
            },
            {
                "name":"age",
                "value":"18"
            }
        ]
    }    
    return transactions(pub, pri, op)

# inflationDest：通货膨胀地址
# setFlags：信任处理标记
# masterWeight：当前账户的签名权重
# lowThreshold：低级安全阈值
# medThreshold：中级安全阈值
# highThreshold：高级安全阈值
# signer：签名者
# signerAccount：签名账户
# weight：签名账户权重
# type：签名类型
def set_option(pub, pri, inflationDest):
    param = copy.deepcopy(param_tpl)
    op = {
        "optype": "set_option",
        "source": pub,
        "setFlags":3,
        "masterWeight":240,
        "lowThreshold":10,
        "medThreshold":110,
        "highThreshold":210,
        "signer":{
            "signerAccount":"0x02dcbb90da7ac2d010b62f9069c3a814c642b0d22efbd7d6e4552c535170b0eaa3",
            "weight":140,
            "type":"25519"
        }
	}
    if inflationDest != None:
        op['inflationDest'] = inflationDest
    nc = nonce(pub)
    param['nonce'] = str(nc)
    sucess, ret = transactions(pub, pri ,op)
    if sucess:
        print 'set_option success', ret
    else:
        print 'set_option false', ret

# 三十四
# type: 0：原生资产；1：长度小于4，例如BTC； 2：长度大于4，例如股票代码
def change_trust(source, pri, issuer, type, code, limit):
    op = {
        "optype": "change_trust",
        "source": source,
        "asset": {
            "issuer": issuer,
            "type": type,
            "code": code
        },
        "limit": str(limit)
    }
    
    return transactions(source, pri, op)

# 三十五
# authorize 信任true  不信任false
def allow_trust(issuer, pri, trustor, type, code, authorize):
    op = {
        "optype": "allow_trust",
        "source": issuer,
        "trustor": trustor,
        "asset": {
            "issuer": issuer,
            "type": type,
            "code": code
        },
        "authorize": authorize
    }    
    return transactions(issuer, pri, op)
    




# 二十九.
# type: 0：原生资产；1：长度小于4，例如BTC； 2：长度大于4，例如股票代码
def payment(pub, pri, to, issuer, type, code, amount):
    asset = {
        "type": type
    }
    if  type != 0 :
        asset['issuer'] = issuer
        asset['code'] = code
    op = {
        "optype": "payment",
        "source": pub,
        "destination": to,
        "asset": asset,
        "amount": str(amount)
    }
    
    return transactions(pub, pri, op)


def path_payment(pub, pri, destination):
    send_asset = {
        "issuer":"0x02de2697632c83ca487632c8092facee01d076987c1d451a36f431775f0bba0419",
        "type":2,
        "code":"Sheep"
    }

    dest_asset = {
        "issuer":"0x02de2697632c83ca487632c8092facee01d076987c1d451a36f431775f0bba0419",
        "type":2,
        "code":"Wheat"
    }
    
    op = {
        "optype":"path_payment",
        "source":pub,
        "destination": destination,
        "send_asset": send_asset,
        "send_max":"10",
        "dest_asset": dest_asset,
        "dest_amount":"100",
        "path":[
            {
                "issuer":"0x02de2697632c83ca487632c8092facee01d076987c1d451a36f431775f0bba0419",
                "type":1,
                "code":"Fish"
            }
        ]
    }
    
    return transactions(pub, pri, op)

