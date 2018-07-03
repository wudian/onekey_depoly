# -* - coding: UTF-8 -* -  

from base import *
import time



# 
def manage_big_data(pub, pri, value, is_pub):
    op = {
        "optype": "manage_big_data",
        "is_pub": is_pub,
        "source": pub,
        # "type": "JSON",
        "value": value,  # YXNzZGZhc2RmYXNkZmFzZGY=
        # "ext": ".json",
        "memo": "manage_big_data memo"
    }    
    b, ret = transactions(pub, pri, op, True)
    time.sleep(1)
    return b, ret

    
def accounts_bigdata(pub):  
    try:
        httpClient = httplib.HTTPConnection(ip, port, timeout=300)
        httpClient.request('GET', '/v1/accounts/'+pub+'/bigdata')
        response = httpClient.getresponse()
        text = response.read()
        if 200 == response.status:
            ret = json.loads(text)
            return ret
        return False, response.reason
    except Exception, e:
            return repr(e)
    finally:
        if httpClient:
            httpClient.close()


def bigdata(pub, pri, hash1, hash2):  
    param = {
        "source": pub,
        "privkey": pri,        
        "hashs": [hash1, hash2]
    }
    try:
        headers = {"Content-type": "application/x-www-form-urlencoded", "Accept": "text/plain"}
        httpClient = httplib.HTTPConnection(ip, port, timeout=300)
        httpClient.request("POST", "/v1/bigdata", json.dumps(param), headers)
        response = httpClient.getresponse()
        text = response.read()
        if 200 == response.status:
            ret = json.loads(text)
            return ret
        return False, response.reason
    except Exception, e:
            return False, repr(e)
    finally:
        if httpClient:
            httpClient.close()



def thorizedata(pri, dest, hashs):
    param = {
        "privkey": pri,        
        "destination": dest,
        "hashs": hashs
    }
    try:
        headers = {"Content-type": "application/x-www-form-urlencoded", "Accept": "text/plain"}
        httpClient = httplib.HTTPConnection(ip, port, timeout=300)
        httpClient.request("POST", "/v1/bigdata/thorizedata", json.dumps(param), headers)
        response = httpClient.getresponse()
        text = response.read()
        if 200 == response.status:
            ret = json.loads(text)
            return ret
        return False, response.reason
    except Exception, e:
            return False, repr(e)
    finally:
        if httpClient:
            httpClient.close()


def gethorizedata(pri, hashs):
    param = {
        "privkey": pri,        
        # "destination": dest,
        "hashs": hashs
    }
    try:
        headers = {"Content-type": "application/x-www-form-urlencoded", "Accept": "text/plain"}
        httpClient = httplib.HTTPConnection(ip, port, timeout=300)
        httpClient.request("POST", "/v1/bigdata/gethorizedata", json.dumps(param), headers)
        response = httpClient.getresponse()
        text = response.read()
        if 200 == response.status:
            ret = json.loads(text)
            return ret
        return False, response.reason
    except Exception, e:
            return False, repr(e)
    finally:
        if httpClient:
            httpClient.close()