# BUG_REPRO

## Bug 是什么
冲正时把状态写成了 failed；重复冲正检查被移除；结算单查询在未命中时返回 nil,nil 导致上层解引用 panic；结算单状态默认值写错。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestReverseSettlementChain 失败：冲正后状态不是 reversed、重复冲正未报错、缺失单号 panic。
- service.TestSettlementService_SubmitAndReverse 失败：冲正状态断言不通过。
