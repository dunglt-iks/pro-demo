# [AUTH-12] Đăng nhập bằng tên đăng nhập và mật khẩu

**Loại:** Story · **Trạng thái:** Done · **Sprint:** 3 · **Người phụ trách:** team web

## Mô tả

Người dùng nội bộ cần đăng nhập để vào hệ thống. Giai đoạn đầu chỉ có tài khoản do quản trị viên tạo sẵn
(chưa có đăng ký, chưa có quên mật khẩu).

## User story

Là **người dùng nội bộ**, tôi muốn **đăng nhập bằng tên đăng nhập và mật khẩu** để **vào trang chủ và dùng các chức
năng của hệ thống**.

## Tiêu chí chấp nhận

- Trang đăng nhập có ô *Tên đăng nhập*, ô *Mật khẩu* (che ký tự) và nút *Đăng nhập*.
- Đúng tên đăng nhập và mật khẩu → vào trang chủ, hiện "Xin chào <tên hiển thị>".
- Tên đăng nhập không phân biệt hoa thường (`Admin` = `admin`).
- Sai tên đăng nhập hoặc sai mật khẩu → hiện "Sai tên đăng nhập hoặc mật khẩu" (không nói rõ sai cái nào).
- Bỏ trống một trong hai ô → hiện "Nhập tên đăng nhập và mật khẩu".
- Trong lúc chờ, nút đổi thành "Đang đăng nhập…" và không bấm được.
- Chưa đăng nhập mà vào trang chủ → chuyển về trang đăng nhập.
- Phiên đăng nhập có hạn 1 giờ; hết hạn thì lần thao tác kế tiếp chuyển về trang đăng nhập.
- Nút *Đăng xuất* ở trang chủ → về trang đăng nhập; vào lại trang chủ phải đăng nhập lại.

## Ghi chú kỹ thuật

- API: `POST /api/auth/login` (`{username, password}` → `{token, user}`), `GET /api/me` (header `Authorization: Bearer`).
- Token JWT (HS256, 1 giờ) lưu ở `localStorage` khoá `token`.
- Tài khoản mẫu: `admin / Admin@123`, `user1 / User@123`.

## Ngoài phạm vi

- Đăng ký tài khoản mới, quên mật khẩu — làm ở ticket sau.
- Khoá tài khoản khi nhập sai nhiều lần.
