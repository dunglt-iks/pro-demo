---
description: Đề xuất cấu hình ai-kit.yml và chỉ mâu thuẫn giữa hướng dẫn của repo với ai-kit (chạy khi mới nhúng kit)
---
<!-- ai-kit:shell — vỏ do kit init quản lý; muốn tự viết lệnh này thì xoá dòng này -->

Chưa có `.ai-kit/` ở gốc repo thì tự gắn kit theo `kitRef` trong `ai-kit.yml` (mặc định `v2`), rồi làm tiếp:
`git clone --depth 1 -b <kitRef> https://github.com/dunglt-iks/ai-kit .ai-kit && npm ci --prefix .ai-kit --omit=dev`.

Đọc `.ai-kit/process/adopt.md` rồi làm đúng theo đó.
