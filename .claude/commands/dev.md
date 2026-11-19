---
description: Triển khai một việc từ spec đã duyệt (kế hoạch, code, test, review), đồng bộ khi spec đổi, hoặc sửa bug từ issue
argument-hint: <đường dẫn file spec> <mã feature | mã việc nhỏ> | #<số issue>
---
<!-- ai-kit:shell — vỏ do kit init quản lý; muốn tự viết lệnh này thì xoá dòng này -->

Triển khai theo spec: $ARGUMENTS

Chưa có `.ai-kit/` ở gốc repo thì tự gắn kit theo `kitRef` trong `ai-kit.yml` (mặc định `v2`), rồi làm tiếp:
`git clone --depth 1 -b <kitRef> https://github.com/dunglt-iks/ai-kit .ai-kit && npm ci --prefix .ai-kit --omit=dev`.

Đọc `.ai-kit/process/dev.md` rồi làm đúng theo đó với đầu vào trên.
