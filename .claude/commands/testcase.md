---
description: Sinh hoặc đồng bộ test case e2e (markdown) từ spec đã duyệt của một feature
argument-hint: <mã feature>
---
<!-- ai-kit:shell — vỏ do kit init quản lý; muốn tự viết lệnh này thì xoá dòng này -->

Mã feature: $ARGUMENTS

Chưa có `.ai-kit/` ở gốc repo thì tự gắn kit theo `kitRef` trong `ai-kit.yml` (mặc định `v2`), rồi làm tiếp:
`git clone --depth 1 -b <kitRef> https://github.com/dunglt-iks/ai-kit .ai-kit && npm ci --prefix .ai-kit --omit=dev`.

Đọc `.ai-kit/process/testcase.md` rồi làm đúng theo đó với đầu vào trên.
