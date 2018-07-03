# -* - coding: UTF-8 -* -  

from qry import *
from bigdata import *
import threading
import time

lock = threading.Lock() # 引入锁


# 循环开户
def cre_accs():
    for num in range(0, 5):
        pub, pri = genkey()
        b, ret = create_account(pub)
        if b:
            print qryAccount(pub)
        else:
            print 'create_account fail:', ret
        
        b, ret = manage_data(pub, pri)
        if b:
            print qryAccount(pub)['data']
        else:
            print 'manage_data fail:', ret

# 循环转账原生资产
def pay():
    pub1, pri1 = genkey()
    create_account(pub1)
    pub2, pri2 = genkey()
    create_account(pub2)
    print "init: ", qryAccount(pub1)['balance'], qryAccount(pub2)['balance']
    for n in range(0, 5):        
        payment(pub1, pri1, pub2, None, 0, None, 5)
        print n,": ", qryAccount(pub1)['balance'], qryAccount(pub2)['balance']


# 添加信任后，再转账发行资产
def chg_trust_and_pay():
    pub1, pri1 = genkey()
    create_account(pub1)
    pub2, pri2 = genkey()
    create_account(pub2)
    
    type = 1
    code = 'BTC2'
    b, ret = change_trust(pub1, pri1, pub2, type, code, 100)
    if b:
        print pub1, "trust", pub2, code
    else:
        print pub1, "trust", pub2, 'fail', ret

    b, ret = allow_trust(pub2, pri2, pub1, type, code, True)
    if b:
        print pub2, 'allow_trust', pub1
    else:
        print pub2, 'allow_trust', pub1, 'fail', ret

    print qryAccount(pub1)['trustLines'][0]
    payment(pub2, pri2, pub1, pub2, type, code, 5)
    print qryAccount(pub1)['trustLines'][0]

# 管理数据
def bigdata_test():
    pub1, pri1 = genkey()
    b, ret = create_account(pub1)
    if not b:
        print ret
        return
    b, ret1 = manage_big_data(pub1, pri1, "888", True)
    if not b:
        print ret1
        return
    b, ret2 = manage_big_data(pub1, pri1, "999", False)
    if not b:
        print ret2
        return
    b, ret3 = manage_big_data(pub1, pri1, "000", False)
    if not b:
        print ret3
        return
    else:        
        print 'accounts_bigdata', accounts_bigdata(pub1)
        print 'bigdata', bigdata(pub1, pri1, ret1, ret2)

    pub2, pri2 = genkey()
    b, ret = create_account(pub2)
    print thorizedata(pri1, pub2, [ret2, ret3])
    print 'gethorizedata', gethorizedata(pri2, [ret2, ret3])
    

if __name__ == '__main__':
    cre_accs()
    pay()
    chg_trust_and_pay()
    bigdata_test()


    








class myThread (threading.Thread):   #继承父类threading.Thread
    def __init__(self):
        threading.Thread.__init__(self)
        
    def run(self):                   #把要执行的代码写到run函数里面 线程在创建后会直接运行run函数 
        cre_accs()
        

# thread1 = myThread()
# thread1.start()

# thread2 = myThread()
# thread2.start()

# thread3 = myThread()
# thread3.start()

# thread1.join()
# thread2.join()
# thread3.join()
