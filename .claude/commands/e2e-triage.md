---
description: Phân loại test e2e fail (lần chạy ở máy hoặc trên CI), sửa script sai, ghi bug app sau khi người dùng đồng ý
argument-hint: "[run-id | link lần chạy CI | #PR]"
---
<!-- ai-kit:shell — vỏ do kit init quản lý; muốn tự viết lệnh này thì xoá dòng này -->

Đầu vào: $ARGUMENTS

Chưa có `.ai-kit/` ở gốc repo thì tự gắn kit theo `kitRef` trong `ai-kit.yml` (mặc định `v2`), rồi làm tiếp:
`git clone --depth 1 -b <kitRef> https://github.com/dunglt-iks/ai-kit .ai-kit && npm ci --prefix .ai-kit --omit=dev`.

Đọc `.ai-kit/process/e2e-triage.md` rồi làm đúng theo đó với đầu vào trên.
