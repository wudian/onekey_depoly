# -* - coding: UTF-8 -* -  

from base import *


def manage_offer(pub, pri, price, amount):
    op = {
        "optype": "manage_offer",
        "source": pub,
        
        "buying": {
            "issuer": "02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5",
            "type": 1,
            "code": "MTB"
        },
        "selling": {
            "issuer": "02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5",
            "type":1,
            "code": "NTB"
        },
        "amount": amount,
        "price": price
    }    
    return transactions(issuer, pri, op)

    def create_passive_offer(pub, pri, price, amount):
        op = {
            "optype": "create_passive_offer",
            "source": pub,
            
            "buying": {
                "issuer": "02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5",
                "type": 1,
                "code": "MTB"
            },
            "selling": {
                "issuer": "02865c395bfd104394b786a264662d02177897391aba1155f854cb1065b6a444e5",
                "type":1,
                "code": "NTB"
            },
            "amount": amount,
            "price": price
        }    
        return transactions(issuer, pri, op)