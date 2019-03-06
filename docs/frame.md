比原 dapp 开发流程：

1. 编写equity合约，编译合约并实例化
    比原链是基于UTXO模型的区块链系统，equity是比原链的智能合约语言，使用equity语言可以实现许多典型的金融交易案例。以储蓄分红合约的demo为例，其说明如下：    
    
    1）合约简介
    储蓄分红合约指的是项目方发起了一个锁仓计划（即储蓄合约），用户可以在准备期内（dueBlockHeight）参与该计划，按照1：1的比例获取同等数量的储蓄票据资产，同时用户锁仓的资产将放到取现合约中，并且项目方是无法动用的，等到锁仓期限（expireBlockHeight）一到，用户便可以解锁取现合约将自己储蓄的资产连本待息一同取出来。其示意图如下:

![image](/docs/images/diagram.png)

    从上图中可以看出，项目方发布了一个利润为20%的锁仓项目，其中储蓄合约FixedLimitCollect锁定了1000个票据资产，为了提现项目方的诚意，项目方将200个利息资产锁定到利息合约中。用户1调用合约储蓄了500个利息资产，其中500个利息将被锁定在取现合约FixedLimitProfit中，同时用户1获得了500个票据资产，剩余找零的资产将继续锁定在储蓄合约FixedLimitCollect中，而用户2 和 用户3也是相同的流程，直到储蓄合约没有资产为止。取现合约跟储蓄合约的模型大致相同，只是取现合约是由多个UTXO组成的，用户在取现的时候可以并行操作，但是如果合约中的面值不能支持用户一次性取现的话，需要分多次提取，比如用户1拥有500个票据资产，计算的本息应该为600，但是取现的UTXO面值为500,那么用户1一次最多只能取500,剩下的100需要再构造一笔交易来提现。


    2）合约源代码

    ```js
    // 储蓄合约
    import "./FixedLimitProfit"
    contract FixedLimitCollect(assetDeposited: Asset,
                            totalAmountBill: Amount,
                            totalAmountCapital: Amount,
                            dueBlockHeight: Integer,
                            expireBlockHeight: Integer,
                            additionalBlockHeight: Integer,
                            banker: Program,
                            bankerKey: PublicKey) locks billAmount of billAsset {
        clause collect(amountDeposited: Amount, saver: Program) {
            verify below(dueBlockHeight)
            verify amountDeposited <= billAmount && totalAmountBill <= totalAmountCapital
            define sAmountDeposited: Integer = amountDeposited/100000000
            define sTotalAmountBill: Integer = totalAmountBill/100000000
            verify sAmountDeposited > 0 && sTotalAmountBill > 0
            if amountDeposited < billAmount {
                lock amountDeposited of assetDeposited with FixedLimitProfit(billAsset, totalAmountBill, totalAmountCapital, expireBlockHeight, additionalBlockHeight, banker, bankerKey)
                lock amountDeposited of billAsset with saver
                lock billAmount-amountDeposited of billAsset with FixedLimitCollect(assetDeposited, totalAmountBill, totalAmountCapital, dueBlockHeight, expireBlockHeight, additionalBlockHeight, banker, bankerKey)
            } else {
                lock amountDeposited of assetDeposited with FixedLimitProfit(billAsset, totalAmountBill, totalAmountCapital, expireBlockHeight, additionalBlockHeight, banker, bankerKey)
                lock billAmount of billAsset with saver
            }
        }
        clause cancel(bankerSig: Signature) {
            verify above(dueBlockHeight)
            verify checkTxSig(bankerKey, bankerSig)
            unlock billAmount of billAsset
        }
    }
    ```

    ```js
    // 取现合约(本金加利息)
    contract FixedLimitProfit(assetBill: Asset,
                            totalAmountBill: Amount,
                            totalAmountCapital: Amount,
                            expireBlockHeight: Integer,
                            additionalBlockHeight: Integer,
                            banker: Program,
                            bankerKey: PublicKey) locks capitalAmount of capitalAsset {
        clause profit(amountBill: Amount, saver: Program) {
            verify above(expireBlockHeight)
            define sAmountBill: Integer = amountBill/100000000
            define sTotalAmountBill: Integer = totalAmountBill/100000000
            verify sAmountBill > 0 && sTotalAmountBill > 0 && amountBill < totalAmountBill
            define gain: Integer = totalAmountCapital*sAmountBill/sTotalAmountBill
            verify gain > 0 && gain <= capitalAmount
            if gain < capitalAmount {
                lock amountBill of assetBill with banker
                lock gain of capitalAsset with saver
                lock capitalAmount - gain of capitalAsset with FixedLimitProfit(assetBill, totalAmountBill, totalAmountCapital, expireBlockHeight, additionalBlockHeight, banker, bankerKey)
            } else {
                lock amountBill of assetBill with banker
                lock capitalAmount of capitalAsset with saver
            }
        }
        clause cancel(bankerSig: Signature) {
            verify above(additionalBlockHeight)
            verify checkTxSig(bankerKey, bankerSig)
            unlock capitalAmount of capitalAsset
        }
    }
    ```

    注意事项：
    - 时间期限不是具体的时间，而是通过区块高度来大概估算的（平均区块时间间隔大概为2.5分钟）
    - 比原的精度是8, 即 1BTM = 100000000 neu，正常情况下参与计算都是以neu为单位的，然而虚拟机的int64类型的最大值是9223372036854775807，为了避免数值太大导致计算溢出，所以对计算的金额提出了金额限制（即amountBill/100000000）
    - 另外clause cancel是项目方的管理方法，如果储蓄或者取现没有满额，项目方也可以回收剩余的资产

    3）编译并实例化合约
    编译equity合约可以参考一下equity编译器https://github.com/Bytom/equity的介绍说明。假如储蓄合约FixedLimitCollect的参数如下：

        assetDeposited          :c6b12af8326df37b8d77c77bfa2547e083cbacde15cc48da56d4aa4e4235a3ee
        totalAmountBill         :10000000000
        totalAmountCapital      :20000000000
        dueBlockHeight          :1070
        expireBlockHeight       :1090
        additionalBlockHeight   :1100
        banker                  :0014dedfd406c591aa221a047a260107f877da92fec5
        bankerKey               :055539eb36abcaaf127c63ae20e3d049cd28d0f1fe569df84da3aedb018ca1bf

    其中bankerKey是用户的publicKey，可以通过比原链的接口list-pubkeys来获取，注意需要保存一下对应的rootXpub和Path，否则无法得到正确的签名结果。

    实例化合约命令如下：
    ./equity deposit/FixedLimitCollect --instance c6b12af8326df37b8d77c77bfa2547e083cbacde15cc48da56d4aa4e4235a3ee 10000000000 20000000000 1070 1090 1100 0014dedfd406c591aa221a047a260107f877da92fec5 055539eb36abcaaf127c63ae20e3d049cd28d0f1fe569df84da3aedb018ca1bf

2. 发布合约交易
   发布合约交易即将资产锁定到合约中。由于目前无法在比原的dashboard上构造合约交易，所以需要借助外部工具来发送合约交易，比如postman。按照上述示意图所示，项目方需要发布1000个储蓄资产的储蓄合约和200个利息资产取现合约。假设项目方需要发布1000个deposit资产（如果精度为8,所以1000个在比原链中表示为100000000000）的锁仓合约，那么他需要将对应数量的票据锁定在储蓄合约中，其交易模板如下：
   
    ```js
    {
      "base_transaction": null,
      "actions": [
        {
          "account_id": "0ILGLSTC00A02",
          "amount": 20000000,
          "asset_id": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
          "type": "spend_account"
        },
        {
          "account_id": "0ILGLSTC00A02",
          "amount": 100000000000,
          "asset_id": "13016eff73ffb7539a69e122f80f5c1cc94446773ac3f64dec290429f87e73b3",
          "type": "spend_account"
        },
        {
          "amount": 100000000000,
          "asset_id": "13016eff73ffb7539a69e122f80f5c1cc94446773ac3f64dec290429f87e73b3",
          "control_program": "20055539eb36abcaaf127c63ae20e3d049cd28d0f1fe569df84da3aedb018ca1bf160014dedfd406c591aa221a047a260107f877da92fec5024c04024204022e040500c817a8040500e40b540220c6b12af8326df37b8d77c77bfa2547e083cbacde15cc48da56d4aa4e4235a3ee4d4302597a64370200005479cda069c35b790400e1f5059600a05c797ba19a53795579a19a695a790400e1f5059653790400e1f505967800a07800a09a695c797b957c9600a069c35b797c9f9161645b010000005b79c2547951005e79895d79895c79895b7989597989587989537a894ca4587a64980000005479cd9f6959790400e1f5059653790400e1f505967800a07800a09a5c7956799f9a6955797b957c967600a069c3787c9f91616481000000005b795479515b79c1695178c2515d79c16952c3527994c251005d79895c79895b79895a79895979895879895779895679890274787e008901c07ec1696393000000005b795479515b79c16951c3c2515d79c16963a4000000557acd9f69577a577aae7cac890274787e008901c07ec169515b79c2515d79c16952c35c7994c251005d79895c79895b79895a79895979895879895779895679895579890274787e008901c07ec1696332020000005b79c2547951005e79895d79895c79895b7989597989587989537a894ca4587a64980000005479cd9f6959790400e1f5059653790400e1f505967800a07800a09a5c7956799f9a6955797b957c967600a069c3787c9f91616481000000005b795479515b79c1695178c2515d79c16952c3527994c251005d79895c79895b79895a79895979895879895779895679890274787e008901c07ec1696393000000005b795479515b79c16951c3c2515d79c16963a4000000557acd9f69577a577aae7cac890274787e008901c07ec16951c3c2515d79c1696343020000547acd9f69587a587aae7cac747800c0",
          "type": "control_program"
        }
    ],
      "ttl": 0,
      "time_range": 1521625823
    }
    ```

    合约交易成功后，合约control_program对应的utxo将会被所有用户查询到，使用比原链的接口list-unspent-outputs即可查询。

3. dapp架构
   dapp demo 的架构大致如下：

   ![image](/docs/images/frame.png)

    - dapp前端
    参考前端架构说明


    - dapp缓冲服务器
  
    缓冲服务器主要是为了在管理合约utxo层面做一些效率方面的处理，此外对dapp的相关交易记录也进行了存储。缓冲服务器的数据来源于bycoin服务器，而bycoin服务器是比原链的去中心化钱包服务器，比原官方插件钱包默认连接的就是该服务器。尽管bycoin服务器的也对比原链的所有utxo进行了管理，但是由于utxo数量比较大，如果直接在该层面处理会导致DAPP性能不佳，所以建议用户自己构建自己的缓冲服务器做进一步处理。此外，DAPP开发者也可以搭建了自己的去中心化钱包服务器，并且自己开发相关的插件。

    以储蓄分红合约dapp的demo版为例，其架构说明如下：

    1）缓冲服务器构成，目前设计了3张数据表：base、utxo和balance表。其中base表用于初始化该dapp关注的合约program，即在查找utxo集合的时候，仅仅只需过滤出对应的program和资产即可; utxo表是该dapp合约的utxo集合，其数据是从bycoin服务器中实时同步过来的，主要是为了提高dapp的并发性; balance表是为了记录用户参与该合约的交易列表。

    2）后端服务由API进程和同步进程组成，其中API服务用于管理对外的用户请求，而同步进程包含了两个方面：一个是从bycoin服务器同步utxo，另一个是则是通过区块链浏览器查询交易状态

    3）项目管理员调用update-base接口更新dapp关注的合约program和asset。而utxo同步进程会根据base表的记录来定时扫描并更新本地的utxo表中的信息，并且根据超时时间定期解锁被锁定的utxo

    4）用户在调用储蓄或取现之前需要查询合约的utxo是否可用，可用的utxo集合中包含了未确认的utxo。倘若所有合约utxo都被锁定了，则会缩短第一个utxo的锁定时间为60s，设置该时间间隔是为了保证未确认的交易被成功验证并生成未确认的utxo。如果该时间间隔并没有产生新的utxo，则认为前面一个用户并没有产生交易，则60s后可以再次花费该utxo。

    5）用户发送交易成功后会生成两条balance记录表，默认状态是失败的，其中交易ID用于向区块链浏览器查询交易状态，如果交易成功则会更新balance的交易状态。此外，前端页面的balance列表表只显示交易成功的记录。


    
