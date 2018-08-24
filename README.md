# my ethereum


自己动手写一个简单的区块链，按照顺序实现如下功能

区块链自身的功能
1. 最基本功能
2. 工作量证明
3. 持久性：选择Redis数据库实现。

与区块链的交互方式
1. 通过http、网页交互
2. 通过命令行交互

入门
----
查看我的[博客](https://ajinisi.github.io/2018/08/18/blockchain/)了解项目的整体


项目结构
----
目录 | 说明 
:-: | :-:
BI.go | 新建区块
block.go | 区块的定义，序列化和反序列化函数，新建区块函数
index.html | 交互的网页
main.go | 主程序


开发指南
----
* 在控制台使用go命令开启区块链
```
$ go run *.go
```
* 在浏览器中打开http://localhost:8080
* 点击新建区块，控制台会显示挖矿过程


特别声明
----
本项目在交互方面和持久化方面分别受[Coral Health](https://github.com/mycoralhealth/blockchain-tutorial.git)
和[Ivan Kuznetsov](https://github.com/Jeiwan/blockchain_go.git)
很大的启发，在此感谢开源运动和分享精神，Respect!

License
----
GPL