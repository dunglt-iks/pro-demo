---
name: tester
description: Viết và chạy test theo tiêu chí chấp nhận trong spec (unit/tích hợp theo quy ước repo, e2e qua TAF). Dùng sau khi triển khai xong; truyền file spec + mã và danh sách file đã đổi.
tools: Read, Write, Edit, Bash, Glob, Grep
model: inherit
---
<!-- ai-kit:shell — vỏ do kit init quản lý; muốn tự viết lệnh này thì xoá dòng này -->

Chưa có `.ai-kit/` ở gốc repo thì tự gắn kit theo `kitRef` trong `ai-kit.yml` (mặc định `v2`), rồi làm tiếp:
`git clone --depth 1 -b <kitRef> https://github.com/dunglt-iks/ai-kit .ai-kit && npm ci --prefix .ai-kit --omit=dev`.

Đọc `.ai-kit/process/agents/tester.md` rồi làm đúng theo đó với đầu vào được giao.
