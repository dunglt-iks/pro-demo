---
name: spec-reviewer
description: Soát độc lập spec vừa viết (tiêu chí kiểm chứng được, nhánh lỗi, mã, wireframe, Change Log) và đối chiếu spec ↔ prototype (kit check-prototype + nhãn, quy tắc hiển thị, luồng bấm được) trước khi người dùng duyệt. /spec và /prototype tự gọi. Chỉ đọc và báo cáo, không sửa file.
tools: Read, Grep, Glob, Bash
model: inherit
---
<!-- ai-kit:shell — vỏ do kit init quản lý; muốn tự viết lệnh này thì xoá dòng này -->

Chưa có `.ai-kit/` ở gốc repo thì tự gắn kit theo `kitRef` trong `ai-kit.yml` (mặc định `v2`), rồi làm tiếp:
`git clone --depth 1 -b <kitRef> https://github.com/dunglt-iks/ai-kit .ai-kit && npm ci --prefix .ai-kit --omit=dev`.

Đọc `.ai-kit/process/agents/spec-reviewer.md` rồi làm đúng theo đó với đầu vào được giao.
