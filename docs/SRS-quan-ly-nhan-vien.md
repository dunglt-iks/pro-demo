# ĐẶC TẢ YÊU CẦU PHẦN MỀM — QUẢN LÝ NHÂN VIÊN

| Phiên bản | Ngày | Người viết | Ghi chú |
|---|---|---|---|
| 0.1 | 2026-10-02 | BA | Bản đầu |
| 0.3 | 2026-10-20 | BA | Thêm khoá tài khoản, phân quyền |

Trạng thái: **chưa triển khai**.

## 1. Giới thiệu

### 1.1 Mục đích

Tài liệu mô tả yêu cầu của phân hệ Quản lý nhân viên: quản trị viên xem, tìm, thêm, sửa, khoá nhân viên.

### 1.2 Phạm vi

Gồm danh sách, tìm kiếm, thêm, sửa, khoá / mở khoá nhân viên. Không gồm chấm công, lương, nghỉ phép.

### 1.3 Thuật ngữ

| Thuật ngữ | Nghĩa |
|---|---|
| Nhân viên | Người có tài khoản trong hệ thống |
| QTV | Quản trị viên |
| Khoá | Nhân viên không đăng nhập được nữa, dữ liệu giữ nguyên |

## 2. Mô tả tổng quan

### 2.1 Người dùng

- **QTV:** toàn quyền trên phân hệ.
- **Nhân viên:** chỉ xem được thông tin của chính mình (không thuộc phạm vi tài liệu này).

### 2.2 Ràng buộc

- Dùng chung tài khoản đăng nhập hiện có (xem `docs/features/login.md`).
- Giao diện tiếng Việt.

## 3. Yêu cầu chức năng

**REQ-001 — Danh sách nhân viên.** QTV xem danh sách nhân viên gồm mã, họ tên, tên đăng nhập, phòng ban, trạng thái
(Đang làm / Đã khoá). Mặc định sắp theo họ tên. Phân trang 20 dòng.

**REQ-002 — Tìm kiếm.** QTV tìm theo họ tên hoặc tên đăng nhập (chứa chuỗi, không phân biệt hoa thường), lọc theo
phòng ban và trạng thái. Không có kết quả thì hiện "Không tìm thấy nhân viên".

**REQ-003 — Thêm nhân viên.** QTV nhập họ tên, tên đăng nhập, email, phòng ban. Hệ thống tạo mật khẩu tạm và hiện cho
QTV một lần. Tên đăng nhập và email không được trùng. Họ tên, tên đăng nhập bắt buộc.

**REQ-004 — Sửa nhân viên.** QTV sửa họ tên, email, phòng ban. Không sửa tên đăng nhập.

**REQ-005 — Khoá / mở khoá.** QTV khoá nhân viên sau khi xác nhận. Nhân viên bị khoá không đăng nhập được; nếu đang
đăng nhập thì lần thao tác kế tiếp bị đưa về trang đăng nhập. QTV không khoá được chính mình.

**REQ-006 — Phân quyền.** Chỉ QTV vào được phân hệ. Nhân viên thường vào đường dẫn của phân hệ thì hiện
"Bạn không có quyền truy cập".

## 4. Yêu cầu phi chức năng

**NFR-001.** Danh sách 1.000 nhân viên tải dưới 1 giây.

**NFR-002.** Mọi thao tác thêm / sửa / khoá ghi nhật ký: ai, lúc nào, thay đổi gì.

## 5. Dữ liệu

| Trường | Kiểu | Bắt buộc | Ghi chú |
|---|---|---|---|
| Mã nhân viên | Chuỗi | Có | Tự sinh, dạng NV0001 |
| Họ tên | Chuỗi | Có | Tối đa 100 ký tự |
| Tên đăng nhập | Chuỗi | Có | Duy nhất, không đổi |
| Email | Chuỗi | Không | Duy nhất nếu có |
| Phòng ban | Danh mục | Không | Kinh doanh, Kỹ thuật, Kế toán, Nhân sự |
| Trạng thái | Danh mục | Có | Đang làm / Đã khoá |

## 6. Câu hỏi mở

1. Mật khẩu tạm gửi qua email hay QTV tự báo cho nhân viên?
2. Có cần xoá hẳn nhân viên không, hay chỉ khoá?
