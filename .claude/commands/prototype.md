---
description: Dựng/cập nhật prototype HTML tĩnh từ màn hình trong spec, đồng bộ ngược thay đổi giao diện về spec
argument-hint: <đường dẫn file spec> <mã feature hoặc module>
---
<!-- ai-kit:shell — vỏ do kit init quản lý; muốn tự viết lệnh này thì xoá dòng này -->

Làm prototype cho: $ARGUMENTS

Chưa có `.ai-kit/` ở gốc repo thì tự gắn kit theo `kitRef` trong `ai-kit.yml` (mặc định `v2`), rồi làm tiếp:
`git clone --depth 1 -b <kitRef> https://github.com/dunglt-iks/ai-kit .ai-kit && npm ci --prefix .ai-kit --omit=dev`.

Đọc `.ai-kit/process/prototype.md` rồi làm đúng theo đó với đầu vào trên.
