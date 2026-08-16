# BUG_REPRO

## Bug 是什么
AppError 的 Unwrap 被改成返回 nil；参保人和结算单的未命中错误被改成 internal error，导致 not found 语义丢失。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestNotFoundErrorsRemainNotFound 失败：查询缺失参保人/结算单不再命中 ErrNotFound，错误码也不对。
