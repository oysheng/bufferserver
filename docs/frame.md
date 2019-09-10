# 比原链DAPP开发框架

从目前已经发布的`DAPP`来看，`DAPP`架构大致可以分成3种类型：插件钱包模式、全节点钱包模式和兼容模式。
- 插件钱包模式是借助封装了钱包的浏览器插件通过`RPC`协议与区块链节点通信，插件在运行时会将`Web3`框架注入到`DAPP`前端页面中，然后`DApp`通过`Web3`来与区块链节点通信。
- 全节点钱包模式需要项目方同步并持有一个区块链节点，并对外提供一个浏览器环境与用户进行交互。
- 兼容模式可以在插件钱包和全节点钱包下同时使用，即上述两种方式可以自由切换，安全性能相对较高。

比原链`DAPP`的架构模式跟账户模型`DAPP`的插件钱包模式有些相似，都是由`DAPP`前端、插件钱包和合约程序共同组成，其区别在于账户模型的插件钱包是直接连接区块链节点，而比原链插件钱包连接的是比原的去中心化钱包服务器`blockcenter`，该服务器主要是为了管理插件钱包的相关信息。此外，比原链是`UTXO`模型的区块链系统，合约程序存在于无状态的`UTXO`中，如果要实现这样一个具体的`DAPP`，就需要在前端和后端多做一些逻辑的处理。

比原链DAPP开发流程具体如下：

## 1. 编写、编译并实例化智能合约

### 编写智能合约
    
比原链的智能合约是由`Equity`语言编写的，使用`Equity`语言可以实现许多典型的金融模型案例。`Equity`智能合约模板结构如下：

```
contract contract_name(...) locks valueAmount of valueAsset {
  clause clause_name(...) {
    ...
    lock/unlock ...
  }
  ...
}
```

`Equity`语法结构简单，语句意思明确，编写比原智能合约可以参考[`Equity`合约介绍](https://docs.bytom.io/mydoc_smart_contract_overview.cn.html)，文档中对`Equity`语言的语法和编译方法都做了详细的介绍。此外，文档还对一些典型的[模板合约](https://docs.bytom.io/mydoc_contract_template.cn.html)进行了介绍，开发者可以进行参考。

### 编译并实例化合约
    
编译合约目前支持两种方式：一种是使用`Equity`编译工具，另一种是调用比原链中编译合约的`RPC`接口`compile`; 而合约实例化是为了将合约脚本按照用户设定的参数进行锁定，编译并实例化合约可以参考[编译并实例化合约](https://docs.bytom.io/mydoc_smart_contract_build.cn.html)的上半部分说明，该文档不仅介绍了合约的参数构造说明，还对编译合约的步骤进行详细说明。而编译器以及相关工具位于`github`项目[`Equity`](https://github.com/Bytom/equity)中，用户可以下载源代码并编译使用。

工具编译和实例化示例如下：

```sh
// compile
./equity [contract_name] --bin

// instance
./equity [contract_name] --instance [arguments ...]
```

## 2. 部署合约

部署合约即发送合约交易，调用比原链的`build-transaction`接口将指定数量的资产发送到合约`program`中，只需将输出`output`中接收方`control_program`设置为指定合约即可。用户可以参考[合约交易说明](https://docs.bytom.io/mydoc_smart_contract_build.cn.html)中的锁定合约章节，交易的构造按照文档中介绍进行参考即可。如果合约交易发送成功，并且交易已经成功上链，便可以通过调用`API`接口`list-unspent-outputs`来查找该合约的`UTXO`。

部署合约交易模板大致如下：

```js
{
  "actions": [
    // inputs
    {
	// btm fee
    },
    {
	amount, asset, spend_account
	// spend user asset
    },

    // outputs
    {
	amount, asset, contract_program
	// receive contract program with instantiated result
    }
  ],
  ...
}
```

## 3. 搭建DAPP架构

比原链`DAPP`主要由前端和后端两部分组成：`DAPP`前端负责的功能包括页面的展示、与插件钱包的连接、以及与后端的交互；`DAPP`后端需要根据合约的实际情况做不同的业务逻辑处理，但是都需要包含同步合约交易和`UTXO`信息等功能。因此，比原`DAPP`的总体框架模型大致如下：

![image](/docs/images/dapp_frame.png)
   
### DAPP前端
    
`DAPP`前端除了页面展示之外，还包含三个重要的业务模块：第一个是与插件钱包的交互，第二个是对输入合约参数进行验证逻辑的处理，第三个是与后端服务器的交互。插件钱包是用户跟区块链节点服务器通信的窗口，比原的插件钱包除了与后台服务器进行交互之外，还包含一些本地业务逻辑处理的接口`API`，具体内容可以参考一下[DAPP开发者向导](https://github.com/Bytom/Bystore/wiki/Dapp-Developer-Guide)。由于比原链是基于`UTXO`模型的区块链系统，交易是由多输入和多输出构成的结构，并且交易输入或输出的位置也需要按照顺序来排列，因此开发`DAPP`需要前端处理一些构建交易的逻辑。除此之外，合约中的`lock-unlock`语句中涉及到数量的计算需要根据抽象语法树来进行预计算，计算的结果将用于构建交易，而`verify`、`if-else`等其他语句类型也需要进行相关的预校验，从而防止用户在执行合约的时候报错。

从功能层面来说，前端主要包含页面的设计、插件的调用、合约交易逻辑的处理、后端服务器的交互等。接下来对这几个重要的部分展开说明： 

- 1）前端页面的设计主要是网页界面的设计，这个部分开发者可以自己的喜好来进行设计
    
- 2）插件钱包已经进行了结构化的封装，并且提供了外部接口给`DAPP`开发者调用，开发者只需要将插件的参数按照规则进行填充，具体请参考[DAPP开发者向导](https://github.com/Bytom/Bystore/wiki/Dapp-Developer-Guide)

- 3）比原链的合约交易是多输入多输出的交易结构，前端需要进行一些预判断逻辑的处理，然后再选择合适的合约交易模板结构。

- 4）DAPP的插件连接的是去中心化的`blockcenter`服务器，`blockcenter`从比原节点服务器上同步所有区块信息和交易信息，该部分`RPC`调用在插件钱包层进行了高度的封装，用户只需按照`API`接口调用即可。除此之外，需要开发者搭建一个后端服务器，主要用于处理`DAPP`的业务逻辑，同时实时同步合约交易和`UTXO`等状态信息。

前端逻辑处理流程大致如下：

- 调用插件，比原的`chrome`插件源码位于[Bytom-JS-SDK](https://github.com/Bytom/Bytom-JS-SDK)，开发比原`DAPP`时调用插件的说明可以参考[Dapp Developer Guide](https://github.com/Bytom/Bystore/wiki/Dapp-Developer-Guide)，其网络配置如下：

    ```js
    window.addEventListener('load', async function() {

      if (typeof window.bytom !== 'undefined') {
        let networks = {
            solonet: ... // solonet bycoin url 
            testnet: ... // testnet bycoin url 
            mainnet: ... // mainnet bycoin url 
        };

        ...

        startApp();
    });
    ```

- 配置合约参数，可以采用文件配置的方式，该步骤是为了让前端得到需要用到的一些已经固定化的合约参数，其前端配置文件为[`configure.json.js`](https://github.com/Bytom/Bytom-Dapp-Demo/blob/master/contracts/configure.json.js)，其示例模型如下:
  
    ```js
    var config = {
        "solonet": {
            ...         // contract arguments
            "gas": 0.4  // btm fee
        },
        "testnet":{
            ...
        },
        "mainnet":{
            ...
        }
    }
    ```

- 前端预计算处理，如果合约中包含`lock-unlock`语句，并且`Amount`是一个数值表达式，那么前端来提取计算表达式并进行相应的预计算。此外，前端还需要预判下所有可验证的`verify`语句，从而判定交易是否可行，因为一旦前端对这些验证失败，合约将必然验证失败。此外，如果`define`或`assign`语句涉及的变量，前端也需预计算这些变量的值。

- 构建合约交易模板，由于解锁合约是解锁`lock`语句条件，构造交易需要根据`lock`语句或`unlock`语句来进行变换。解锁合约交易是由`inputs`和`outputs`构成，交易的第一个`input`输入一般都是是固定的，即合约`UTXO`的哈希值，除此之外，其他输入输出都需要根据`DAPP`中的实际合约来进行变更，其模型大致如下：

    ```js
    const input = []
    input.push(spendUTXOAction(utxohash))
    ... // other input

    const output = []
    output.push(controlProgramAction(amount, asset, program))
    ... // other output
    ```

- 启动前端服务

    编译前端命令如下：

    ```sh
    npm run build
    ```

    启动之前需要先启动`DAPP`后端服务器，然后再启动前端服务，其前端启动命令如下：

    ```sh
    npm start
    ```

### DAPP后端
  
`DAPP`后端主要负责根据合约实现对应的业务应用，同时提供给前端调用的应用服务接口，并且实时同步合约交易和`UTXO`等状态信息。`blockcenter`服务器是比原链的去中心化钱包服务器，管理着所有的区块、交易和`UTXO`等信息，`DAPP`后端向它定时请求并保持状态的同步更新，比原官方插件钱包默认连接的就是该服务器。此外，`DAPP`开发者也可以根据应用需求搭建自己的去中心化钱包服务器以及相关的钱包插件。`DAPP`后端业务模型是多变的，一般向`blockcenter`同步合约交易和`UTXO`等状态信息是共同必备的操作，除此之外，`DAPP`的业务模型、数据库模型、性能等都需要根据实际情况来进行设计。

`DAPP`后端服务器架构可以参考一下[储蓄分红DAPP后端服务器bufferserver](https://github.com/oysheng/bufferserver)的源代码及其说明。



    
