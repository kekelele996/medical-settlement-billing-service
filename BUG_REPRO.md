# BUG_REPRO

## Bug 是什么
费用上传当日重复检查的日期范围写反导致去重失效；批次汇总条数多记一条；明细查询按 id 倒序；批次金额默认值写错。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestFeeUploadDuplicateAndSummary 失败：当天第二次上传未被拒绝、明细顺序错乱。
- service.TestFeeService_Upload 失败：批次汇总条数不正确。
