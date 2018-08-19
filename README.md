# my ethereum

自己动手写一个简单的区块链，按照顺序实现如下功能

区块链自身的功能
1. 最基本功能
2. 工作量证明
3. 持久性：选择Redis数据库实现。

与区块链的交互方式
1. 通过http、网页交互
2. 通过命令行交互

项目结构
----
目录 | 说明 
:-: | :-:
BI.go | 新建区块
block.go | 区块的定义，序列化和反序列化函数，新建区块函数
index.html | 交互的网页
main.go | 主程序