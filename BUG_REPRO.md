# BUG_REPRO

## Bug 是什么
结算计算里甲类目录值被改坏、乙类不再按比例折算、统筹支付比例被写成自付比例、费用上传金额用单价加数量、预结算构造把明细金额填成 0。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- util.TestSettlementCalculator_Calculate 失败：分类基数、统筹支付、自费金额算错。
- service.TestPresettlementCalcChain 失败：预结算总金额和各项支付金额不等于预期。
