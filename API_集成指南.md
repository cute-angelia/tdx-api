# 📡 API 集成指南

## 概述

这份文档只回答一件事：`怎么把 tdx-api 接进你的系统`。

接口定义、参数细节、返回结构、错误码、版本更新，请统一查：

- [API_接口文档.md](API_接口文档.md)

如果文档与代码看起来不一致，建议以当前服务代码为准，并优先检查：

- `web/server.go`
- `web/server_api_extended.go`
- `web/server_ws.go`

---

## 适用场景

适合以下接入方式：

- 后端服务拉取行情、K 线、分时、成交数据
- 前端页面轮询行情或通过 WebSocket 订阅实时行情
- 量化/分析脚本批量获取股票、ETF、交易日、收益区间数据
- 本地任务系统触发 K 线或成交入库任务

不建议用这份文档替代接口文档。它强调的是“接入策略”，不是“字段字典”。

---

## 快速接入

### 1. 启动服务

```bash
cd web
go run .
```

默认地址：

- HTTP: `http://localhost:8080`
- WebSocket: `ws://localhost:8080`

### 2. 先做健康检查

```bash
curl "http://localhost:8080/api/health"
curl "http://localhost:8080/api/server-status"
```

如果这两步不通，先不要继续接行情接口。

### 3. 选择你要用的能力

常见组合：

- 看盘概览：`/api/quote` + `/api/minute` + `/api/kline?type=day`
- 批量监控：`/api/batch-quote`
- 实时推送：`/ws/quote`
- 历史分析：`/api/kline-history`、`/api/kline-all`、`/api/index/all`
- 成交明细：`/api/trade`、`/api/trade-history`、`/api/trade-history/full`
- 数据准备：`/api/codes`、`/api/stock-codes`、`/api/etf-codes`、`/api/workday`
- 后台入库：`/api/tasks/pull-kline`、`/api/tasks/pull-trade`

---

## 接入原则

### 1. 接口定义只看一个地方

不要在多个文档里分别抄接口清单。接入时：

1. 在本文件确认“该用哪类接口”
2. 去 [API_接口文档.md](API_接口文档.md) 看准确参数和返回结构

### 2. 注意返回单位

这个项目不是所有字段都用常见的“元/股”：

- 价格经常是 `厘`
- 成交量经常是 `手`
- 五档挂单量通常是 `股`

对接 UI、风控或指标计算时，先做单位换算，不要直接把原始值当元或股。

### 3. 先分清“结构体返回”还是“自定义 map 返回”

这会直接影响字段名大小写：

- 直接返回 Go 结构体的接口，通常字段是 `Count`、`List`、`Time`、`Price`
- 手工组装的接口，通常字段是 `count`、`list`、`date`、`meta`

拿不准时，直接查 [API_接口文档.md](API_接口文档.md)。

---

## 推荐超时与重试

### HTTP 请求

推荐客户端超时：

- 普通查询接口：`5-10s`
- `kline-all` / `index/all`：`10-20s`
- 同花顺前复权全量接口：建议 `>=10s`

### 重试策略

建议只对这类场景做有限重试：

- 网络抖动
- 上游瞬时不可用
- 客户端超时

不建议无脑重试：

- 参数错误
- 股票代码错误
- 明确返回业务错误

建议策略：

- 最多重试 `2-3` 次
- 指数退避，例如 `0.5s / 1s / 2s`

---

## 常见接入方案

### 方案一：前端行情看板

推荐：

- 首屏：用 `/api/stock-info` 或 `/api/quote` + `/api/minute`
- 自选列表：用 `/api/batch-quote`
- 高频刷新：优先用 `/ws/quote`

不推荐：

- 前端每只股票分别轮询 `/api/quote`

### 方案二：历史分析脚本

推荐顺序：

1. 用 `/api/stock-codes` 或 `/api/codes` 获取标的池
2. 用 `/api/workday` 或 `/api/workday/range` 确认交易日
3. 用 `/api/kline-history` 做区间拉取
4. 需要全量历史时再用 `/api/kline-all`

### 方案三：逐笔成交研究

推荐：

- 单日分页：`/api/trade-history`
- 单日全量：`/api/minute-trade-all`
- 长时间跨度：`/api/trade-history/full`

注意：

- `/api/trade-history/full` 返回的是自定义小写字段结构
- `price` 在这个接口里是浮点元值，不是文档里很多结构体返回那种原始 `厘`

### 方案四：本地入库任务

流程：

1. `POST /api/tasks/pull-kline` 或 `POST /api/tasks/pull-trade`
2. 记录返回的 `task_id`
3. 轮询 `/api/tasks/{task_id}`
4. 必要时 `POST /api/tasks/{task_id}/cancel`

适合离线准备数据库，不适合给实时前台直接调用。

---

## 代码示例

### Python：查询日 K 线区间

```python
import requests

BASE_URL = "http://localhost:8080"

resp = requests.get(
    f"{BASE_URL}/api/kline-history",
    params={
        "code": "000001",
        "type": "day",
        "start_date": "2024-11-01",
        "end_date": "2024-11-30",
        "limit": 100,
    },
    timeout=10,
)
resp.raise_for_status()
data = resp.json()["data"]

for item in data["List"]:
    print(item["Time"], item["Close"])
```

### JavaScript：批量获取行情

```javascript
const response = await fetch('http://localhost:8080/api/batch-quote', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    codes: ['000001', '600519', '601318']
  })
});

const payload = await response.json();
console.log(payload.data);
```

### JavaScript：订阅实时行情

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/quote?code=000001,600519&interval=3');

ws.onmessage = (event) => {
  const payload = JSON.parse(event.data);
  console.log(payload.data);
};
```

### Python：发起入库任务

```python
import requests

BASE_URL = "http://localhost:8080"

task = requests.post(
    f"{BASE_URL}/api/tasks/pull-kline",
    json={
        "codes": ["000001", "600519"],
        "tables": ["day", "week", "month"],
        "limit": 4,
        "start_date": "2020-01-01",
    },
    timeout=10,
).json()["data"]

task_id = task["task_id"]
print("task:", task_id)
```

---

## 对接排坑

### 1. `/api/kline-history` 现在真的支持日期范围

当前版本里，`start_date` / `end_date` 已接入过滤逻辑。

如果你发现返回结果和预期区间不一致，先检查：

- 日期格式是否为 `YYYYMMDD` 或 `YYYY-MM-DD`
- `limit` 是否把结果又截短了

### 2. `/api/minute` 不会自动回退交易日

如果指定日期没有数据，会返回：

- 原请求日期
- `Count = 0`
- `List = []`

调用方自己决定是否换日期重试。

### 3. `/api/search` 没有结果时不是报错

当前实现是返回空数组，不是“未找到相关股票”的错误。

### 4. 全量接口不要默认给前端直接用

以下接口可能返回很多数据：

- `/api/kline-all`
- `/api/kline-all/tdx`
- `/api/kline-all/ths`
- `/api/index/all`
- `/api/trade-history/full`

更适合作为：

- 离线分析
- 后端预处理
- 定时任务

### 5. 任务接口是异步的

`pull-kline` / `pull-trade` 只是创建任务，不代表任务已经完成。

---

## 生产建议

上线前建议至少补这些能力：

- 认证：如 API Token 或网关鉴权
- 限流：防止高频轮询压垮服务
- 超时：客户端和反向代理都设置
- 监控：记录接口耗时、错误率、任务失败率
- 缓存：对股票代码、交易日、ETF 列表等低频变化接口做缓存

---

## 文档分工

建议长期保持下面这个边界：

- [API_接口文档.md](API_接口文档.md)
  - 定义接口
  - 定义参数
  - 定义返回结构
  - 定义错误码

- [API_集成指南.md](API_集成指南.md)
  - 说明怎么接
  - 说明先调哪些接口
  - 说明超时、重试、缓存和排坑
  - 给出典型调用方案

这样两份文档不容易互相漂移。
