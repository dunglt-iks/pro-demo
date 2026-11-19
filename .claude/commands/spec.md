---
description: Phỏng vấn rồi ghi yêu cầu vào đúng tài liệu spec của repo (trạng thái nháp)
argument-hint: <mô tả ngắn việc cần làm>
---
<!-- ai-kit:shell — vỏ do kit init quản lý; muốn tự viết lệnh này thì xoá dòng này -->

Người dùng muốn làm: $ARGUMENTS

Chưa có `.ai-kit/` ở gốc repo thì tự gắn kit theo `kitRef` trong `ai-kit.yml` (mặc định `v2`), rồi làm tiếp:
`git clone --depth 1 -b <kitRef> https://github.com/dunglt-iks/ai-kit .ai-kit && npm ci --prefix .ai-kit --omit=dev`.

Đọc `.ai-kit/process/spec.md` rồi làm đúng theo đó với đầu vào trên.
