# -* - coding: UTF-8 -* -  

from base import *

# 四.
def qryAccount(pub):  
    try:
        httpClient = httplib.HTTPConnection(ip, port, timeout=300)
        httpClient.request('GET', '/v1/accounts/'+pub)
        response = httpClient.getresponse()
        text = response.read()
        if 200 == response.status:
            ret = json.loads(text)
            if ret['isSuccess']:
                return ret['result']
            else:
                return ret['message']
        return ""
    except Exception, e:
            print repr(e)
    finally:
        if httpClient:
            httpClient.close()


def qryTx(txhash):  
    try:
        httpClient = httplib.HTTPConnection(ip, port, timeout=300)
        httpClient.request('GET', '/v1/transactions/'+txhash)
        response = httpClient.getresponse()
        text = response.read()
        if 200 == response.status:
            ret = json.loads(text)
            if ret['isSuccess']:
                return ret['result']
        return ""
    except Exception, e:
            return repr(e)
    finally:
        if httpClient:
            httpClient.close()