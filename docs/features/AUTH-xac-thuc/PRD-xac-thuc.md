# PRD — Xác thực (AUTH)

| Phiên bản | Cập nhật lần cuối | FSD |
|---|---|---|
| 1.0 | 2026-11-19 | [FSD-xac-thuc.md](FSD-xac-thuc.md) |

> Chuyển từ ticket `docs/features/login.md` (AUTH-12) — file đó giữ làm tham khảo, không cập nhật nữa.

## Change Log

| Phiên bản | Ngày | Mô tả thay đổi | Tác giả |
|---|---|---|---|
| 1.0 | 2026-11-19 | Khởi tạo: AUTH-F01 Đăng nhập (chuyển từ login.md), AUTH-F02 Quên mật khẩu; chốt câu hỏi mở 1–3, BR01/BR05 bỏ khoảng trắng, BR04 8–16 ký tự, BR06 điều kiện hết hiệu lực, rủi ro đã chấp nhận (§1) | BA Team |

---

## 1. Mục tiêu và phạm vi

- **Vấn đề:** người dùng nội bộ cần đăng nhập để dùng hệ thống; người quên mật khẩu hiện không có cách nào tự lấy
  lại quyền truy cập.
- **Mục tiêu:** người dùng đăng nhập bằng tên đăng nhập + mật khẩu; người quên mật khẩu tự đặt lại mật khẩu bằng cách
  xác minh tên đăng nhập + email của tài khoản, không cần nhờ người khác.
- **Ngoài phạm vi:**
  - Gửi email / link đặt lại mật khẩu, mã OTP.
  - Đăng ký tài khoản mới.
  - Khoá tài khoản hoặc giới hạn số lần khi đăng nhập sai hay xác minh sai. **Rủi ro đã chấp nhận:** ai biết tên
    đăng nhập + email của một tài khoản là đặt lại được mật khẩu tài khoản đó.
  - Màn hình nhập / sửa email của tài khoản (thuộc phân hệ Quản lý nhân viên — `docs/SRS-quan-ly-nhan-vien.md`).
  - Đổi mật khẩu khi đang đăng nhập.

## 2. Vai trò

| Vai trò | Mô tả | Quyền chính |
|---|---|---|
| Người dùng nội bộ | Người có tài khoản trong hệ thống (tạo sẵn: hiện là dữ liệu mẫu; sau này do quản trị viên tạo) | Đăng nhập, đăng xuất, tự đặt lại mật khẩu của tài khoản mình |

## 3. Feature Summary

| Code | Tên | Mô tả ngắn | Spec | Prototype |
|---|---|---|---|---|
| AUTH-F01 | Đăng nhập | Đăng nhập bằng tên đăng nhập + mật khẩu, vào trang chủ, đăng xuất | Đã duyệt (dunglt-iks, 2026-11-19) | Đã duyệt (dunglt-iks, 2026-11-19) |
| AUTH-F02 | Quên mật khẩu | Xác minh tên đăng nhập + email rồi đặt mật khẩu mới | Đã duyệt (dunglt-iks, 2026-11-19) | Đã duyệt (dunglt-iks, 2026-11-19) |

## 4. Vai trò × luồng

| Luồng | Người dùng nội bộ |
|---|---|
| AUTH-F01-US01 Đăng nhập | ✓ |
| AUTH-F01-US02 Vào trang cần đăng nhập | ✓ |
| AUTH-F01-US03 Đăng xuất | ✓ |
| AUTH-F02-US01 Xác minh tài khoản | ✓ |
| AUTH-F02-US02 Đặt mật khẩu mới | ✓ |

## 5. Quy tắc nghiệp vụ

| Mã | Quy tắc |
|---|---|
| AUTH-BR01 | Tên đăng nhập không phân biệt hoa thường (`Admin` = `admin`), bỏ khoảng trắng đầu / cuối. |
| AUTH-BR02 | Thông báo lỗi khi đăng nhập hoặc xác minh không được để lộ tên đăng nhập / email nào tồn tại, cũng không nói rõ sai ô nào. |
| AUTH-BR03 | Phiên đăng nhập có hạn 1 giờ kể từ lúc đăng nhập. |
| AUTH-BR04 | Mật khẩu mới có từ 8 đến 16 ký tự (đếm theo ký tự; khoảng trắng tính là ký tự). Không bắt buộc loại ký tự. Được trùng mật khẩu cũ. |
| AUTH-BR05 | Xác minh quên mật khẩu thành công khi tên đăng nhập và email thuộc **cùng một** tài khoản. Email so không phân biệt hoa thường, bỏ khoảng trắng đầu / cuối. Tài khoản chưa có email thì không xác minh được. |
| AUTH-BR06 | Xác minh thành công mở một phiên đặt lại mật khẩu cho đúng tài khoản đó, có hạn 10 phút, dùng được một lần. Phiên hết hiệu lực khi: đã dùng, quá hạn, người dùng xác minh lại (phiên mới thay phiên cũ), bấm "Về đăng nhập" ở bước 2, hoặc tài khoản không còn. |
| AUTH-BR07 | Đặt lại mật khẩu thành công thì mọi phiên đăng nhập cũ của tài khoản hết hiệu lực; người dùng phải đăng nhập lại bằng mật khẩu mới. |

## 6. Yêu cầu phi chức năng

Không có yêu cầu riêng ở lần này. Ràng buộc bảo mật kỹ thuật (băm mật khẩu, không log mật khẩu / token) theo
`CLAUDE.md`.

## 7. Phụ thuộc

- Module phụ thuộc: không có (chưa có module spec nào khác).
- Phụ thuộc dữ liệu: email của tài khoản — phân hệ Quản lý nhân viên (`docs/SRS-quan-ly-nhan-vien.md`, REQ-003/004) sẽ là nơi nhập /
  sửa; phân hệ này chưa triển khai, nên lần này email chỉ có ở dữ liệu mẫu và dữ liệu test.

## 8. Câu hỏi còn mở

| # | Câu hỏi | Liên quan | Trạng thái |
|---|---|---|---|
| 1 | Phiên đặt lại mật khẩu (sau khi xác minh) có hạn bao nhiêu phút? | AUTH-BR06 | Đã chốt: 10 phút |
| 2 | Có giới hạn số lần xác minh sai để chống dò email không? Hiện ai biết tên đăng nhập + email là đặt lại được mật khẩu. | AUTH-F02-US01 | Đã chốt: để sau, chấp nhận rủi ro (§1) |
| 3 | Email cho tài khoản mẫu (`admin`, `user1`) là gì? | AUTH-BR05 | Đã chốt: `admin@project1.local`, `user1@project1.local` |
