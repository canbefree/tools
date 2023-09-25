解决场景： 

每日，每月，每周 ,每年统计；每次通过后计数要加1； 或者只是校验，不加一

每个规则都是一个事务，需要支持回滚操作； 

支持懒加载，范围条件；

async 支持异步判断

简单的规则引擎，更自由的编辑方式，满足所有场景； 



1 重复开发， 新的需求开发工作量比较大。 对不同的业务场景需要编排不同的查询语句。

2 性能低下， 每个规则对应的是不同的查询业务。而且为了保证最终一致性，会有很多幂等锁。


简单规则： 

unique:{};ext{};rule:{}


unique_id: act=limit&uid=a&product=b; 在act=limit这个场景下，用户a购买商品b次数统计；

ext:ident=boker; // 额外信息是经纪人


范围规则；满足条件即可；

限制规则： 每次通过需要计数+1； （需要扩展，扩展格式；）





rule1: 

type: limit 

5;

then {limit++}

else{limit--}


rule2 :

type: daylimit

10;

then {daylimit++}

else{daylimit--}





{

uid=12

ident=1

}

match{

isBoker(ident)

ret (302)

}


match (

每日访问数量 

uid=+

in(indet,[]string{12,12})

async(in(in(indent,[]string{12,12})))

每日+1;

errreturn:  (301)


success()

rollback()

)


product=12&uid=12;     数据uniqkey: product=12; (如果是统计商品购买数量那么，规则是这样子的；)




step: 商品保存更新


type: minute, day, year;

AOP :

-> A -> B -> C -> D -> E

-> A

-> B

-> C

-> A -> B

-> C -> D

    -> E

验证经纪人身份 - 经纪人身份；

rules

in  identifys group [ borker ,id , df]

用户每日访问不超过 100 次；1203;notvip;limit;day;100;

用户是vip每日访问200次；   1203;vip;limit;day;200;

用户报备数量超过5，获得特权； 1203;report;limit;>5;

规则引擎：

go/tokena
