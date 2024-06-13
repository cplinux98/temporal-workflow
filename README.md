## 基于temporal的workflow编排工具

目标： 能够通过解析dsl来运行workflow，通过前端来编排workflow的具体流程，最终基于temporal运行

## 计划

1. 梳理
   1. 梳理使用temporal构建workflow和从0开始构建workflow有什么区别，便于规划开发计划
   2. 梳理使用哪种DSL来运行workflow，
      1. serverlessworkflow，目前这个项目还处于沙盒阶段，没有发布正式版本，每个版本之间变化差距有点大，且需要实现的功能有很多
      2. 简单的基于DAG的DSL表达式，吸收serverlessworkflow一些思想，后续有需求再更新迭代
2. 开发计划
   1. 初步计划使用简单的基于DAG的DSL表达式来实现，基于serverlessworkflow实现的工作量很大，不是现在能够实现的，参考下面的文章和项目先实现初版
      1. https://medium.com/airbnb-engineering/journey-platform-a-low-code-tool-for-creating-interactive-user-workflows-9954f51fa3f8
      2. https://medium.com/@PhakornKiong/my-naive-implementation-of-no-code-low-code-tool-253e678f2456
      3. https://github.com/n8n-io/n8n
   2. 初版实现功能
      1. 能够接收使用event、webhook、手动、cron来运行workflow实例
      2. workflow节点：
         1. 事件接收节点
         2. shell、python代码运行节点
         3. switch分支节点
         4. 执行http节点
      3. workflow实例能够运行shell、python代码任务
      4. workflow实例失败时能够触发通知
      5. 能够查看当前workflow实例运行到了哪个节点







## 参考

国内关于temporal相关很少，要把关于temporal的相关文章归纳起来，看看别人是怎么使用这个工具的。



- [一种实用的 Temporal 架构方法](https://www.infoq.cn/article/rhm7korkk4fxcjdgfdvq)
- [Temporal: Open Source Workflows as Code](https://mikhail.io/2020/10/temporal-open-source-workflows-as-code/)
- [工作流引擎Temporal学习笔记](https://code2life.top/2023/01/23/0070-temporal-notes/)
- [基于事件溯源的任务编排](https://mp.weixin.qq.com/s?__biz=MzI3MDM1OTgxMQ==&mid=2247484094&idx=1&sn=4d45500d2e2ea8f529812a8de89f2464&chksm=ead30ed2dda487c4e36d89364e36bcdda90dfdc6f0ab5f0d01f6165dd7974e2a38967fc961a7&token=645687367&lang=zh_CN&scene=21#wechat_redirect)
- [temporal-knowledge-temporal-io-dsl-overview](https://www.restack.io/docs/temporal-knowledge-temporal-io-dsl-overview)
- [iWF](https://github.com/indeedeng/iwf)
- [Implementing DSL workflows](https://community.temporal.io/t/implementing-dsl-workflows/3413/23)
- [Defining Workflows](https://temporal.io/blog/defining-workflows)
- [Serverless Workflow](https://serverlessworkflow.io/)
- [Turning chaos into order with workflows. Introduction to Temporal](https://www.agilevision.io/blog/turning-chaos-into-order-with-workflows-introduction-to-temporal/)
- [My Naive implementation of no-code/low-code tool](https://medium.com/@PhakornKiong/my-naive-implementation-of-no-code-low-code-tool-253e678f2456)
- [TemporalDSL](https://github.com/PhakornKiong/TemporalDSL/tree/31-mar)
- [How to visualize Temporal.io workflows](https://www.reddit.com/r/ExperiencedDevs/comments/13s8kdb/opinions_about_temporalio_microservice/)
- [Using Temporal.io to build Long running Workflows](https://sachinsu.github.io/posts/temporalworkflow/)
- [why-temporal](https://www.swyx.io/why-temporal)



来自timd.cn

- [Temporal合集][1. Temporal demo](http://timd.cn/temporal/demo/)

- [Temporal合集][2. Temporal 应用程序开发-基础](http://timd.cn/temporal/application-development-foundations/)

- [Temporal合集][3. Temporal 概念-Temporal](http://timd.cn/temporal/concepts/temporal/)

- [Temporal合集]

  4. Temporal 概念-Workflow

  - [Schedule](http://timd.cn/temporal/concepts/workflow/schedule/)
  - [Schedule Demo](http://timd.cn/temporal/concepts/workflow/schedule/demo/)

- [Temporal合集][5. Temporal 概念-Activity](http://timd.cn/temporal/concepts/activity/)

- [Temporal合集][6. Temporal 概念-Worker](http://timd.cn/temporal/concepts/worker/)

- [Temporal合集][7. Temporal 概念-Task](http://timd.cn/temporal/concepts/task/)

- [Temporal进阶][1. Temporal - Autoscaling Workers 方案](http://timd.cn/temporal/advance/autoscaling-workers/)

- [Temporal进阶][2. Temporal - Temporal 平台化架构](http://timd.cn/temporal/advance/platformization-architecture/)



参考样例：

https://github.com/hatchet-dev/hatchet-workflows

https://medium.com/airbnb-engineering/journey-platform-a-low-code-tool-for-creating-interactive-user-workflows-9954f51fa3f8

https://www.linkedin.com/posts/sestegra_build-fault-tolerant-distributed-cloud-applications-activity-6952536339552595968-l_AG

Go 示例：https://github.com/temporalio/samples-go/tree/main/dsl
Java 示例： https: [//github.com/temporalio/samples-java/tree/main/src/main/java/io/temporal/samples/dsl](https://github.com/temporalio/samples-java/tree/main/src/main/java/io/temporal/samples/dsl)
演示： https: [//github.com/tsurdilo/swtemporal](https://github.com/tsurdilo/swtemporal)

希望这可以帮您朝着正确的方向前进。

