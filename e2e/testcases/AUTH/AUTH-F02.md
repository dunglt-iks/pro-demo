---
feature: AUTH-F02
spec: docs/features/AUTH-xac-thuc/FSD-xac-thuc.md
specHash: "b827bfe3abfac04c"
status: Đã duyệt
reviewer: "dunglt-iks"
approvedHash: "b60ecccfd4763623"
---

<!--
  Người dùng của AUTH-F02 chưa đăng nhập nên file không có persona mặc định. "Tài khoản T" là tài khoản riêng của
  test (tên đăng nhập, tên hiển thị "Người dùng E2E", email, mật khẩu "Test@1234"), không dùng admin / user1.
-->

## TC-AUTH-F02-001 · Quên mật khẩu, đặt mật khẩu mới rồi đăng nhập lại
- Truy vết: AUTH-F02-US01, AUTH-F02-US02
- Ưu tiên: P1
- Màn hình: AUTH-F02-SCR01, AUTH-F02-SCR02, AUTH-F01-SCR01
- Tiền điều kiện: tài khoản T có email; chưa đăng nhập

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Mở trang đăng nhập, bấm link *Quên mật khẩu?* (7) | Hiện màn "Quên mật khẩu" với hướng dẫn "Nhập tên đăng nhập và email của tài khoản để đặt lại mật khẩu." |
| 2 | Nhập tên đăng nhập (3) và email (4) của T, bấm *Tiếp tục* (6) | Hiện màn "Đặt mật khẩu mới", dòng "Tài khoản: <tên đăng nhập của T>" (2), gợi ý "Từ 8 đến 16 ký tự" dưới ô Mật khẩu mới |
| 3 | Nhập *Mật khẩu mới* (3) = "Moi@12345", *Nhập lại mật khẩu mới* (4) = "Moi@12345", bấm *Đặt lại mật khẩu* (6) | Về trang đăng nhập, hiện "Đặt lại mật khẩu thành công. Đăng nhập bằng mật khẩu mới." |
| 4 | Tải lại trang | Không còn thông báo thành công |
| 5 | Đăng nhập bằng tên đăng nhập của T và mật khẩu cũ "Test@1234" | Báo "Sai tên đăng nhập hoặc mật khẩu" |
| 6 | Đăng nhập bằng tên đăng nhập của T và "Moi@12345" | Vào trang chủ, hiện "Xin chào Người dùng E2E" |

## TC-AUTH-F02-002 · Xác minh không phân biệt hoa thường, bỏ khoảng trắng đầu / cuối
- Truy vết: AUTH-F02-US01, AUTH-BR05
- Ưu tiên: P2
- Màn hình: AUTH-F02-SCR01
- Tiền điều kiện: tài khoản T có tên đăng nhập và email viết thường; chưa đăng nhập
- Dữ liệu kiểm tra:

| Tên đăng nhập nhập vào | Email nhập vào | Kết quả mong đợi |
|---|---|---|
| tên của T viết HOA | email của T | Sang màn "Đặt mật khẩu mới" |
| tên của T | email của T viết HOA | Sang màn "Đặt mật khẩu mới" |
| "  " + tên của T + "  " | "  " + email của T + "  " | Sang màn "Đặt mật khẩu mới" |

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Mở màn "Quên mật khẩu", nhập tên đăng nhập (3) và email (4) theo dữ liệu, bấm *Tiếp tục* (6) | Như cột Kết quả mong đợi |

## TC-AUTH-F02-003 · Bỏ trống tên đăng nhập hoặc email
- Truy vết: AUTH-F02-US01-A1
- Ưu tiên: P2
- Màn hình: AUTH-F02-SCR01
- Tiền điều kiện: chưa đăng nhập
- Dữ liệu kiểm tra:

| Tên đăng nhập | Email |
|---|---|
| (trống) | (trống) |
| (trống) | email của T |
| tên của T | (trống) |
| "   " (chỉ khoảng trắng) | email của T |
| tên của T | "   " (chỉ khoảng trắng) |

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Mở màn "Quên mật khẩu", nhập theo dữ liệu, bấm *Tiếp tục* (6) | Báo "Nhập tên đăng nhập và email" (5); vẫn ở màn "Quên mật khẩu" |

## TC-AUTH-F02-004 · Tên đăng nhập và email không khớp
- Truy vết: AUTH-F02-US01-A2, AUTH-BR02
- Ưu tiên: P2
- Màn hình: AUTH-F02-SCR01
- Tiền điều kiện: tài khoản T có email; tài khoản U có email khác; tài khoản K không có email; chưa đăng nhập
- Dữ liệu kiểm tra:

| Tên đăng nhập | Email |
|---|---|
| tên của T | "khong-co@project1.local" |
| tên không tồn tại | email của T |
| tên của T | email của U |
| tên của K | "k@project1.local" |

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Mở màn "Quên mật khẩu", nhập theo dữ liệu, bấm *Tiếp tục* (6) | Báo "Tên đăng nhập hoặc email không đúng" (5) — cùng một câu cho mọi dòng; vẫn ở màn "Quên mật khẩu" |

## TC-AUTH-F02-005 · Mất kết nối máy chủ khi xác minh
- Truy vết: AUTH-F02-US01-E1
- Ưu tiên: P3
- Màn hình: AUTH-F02-SCR01
- Tiền điều kiện: tài khoản T có email; máy chủ không phản hồi (mất kết nối)

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Nhập tên đăng nhập và email của T, bấm *Tiếp tục* (6) | Báo "Không kết nối được máy chủ" (5); vẫn ở màn "Quên mật khẩu" |

## TC-AUTH-F02-006 · Lỗi hệ thống khi xác minh
- Truy vết: AUTH-F02-US01-E2
- Ưu tiên: P3
- Màn hình: AUTH-F02-SCR01
- Tiền điều kiện: tài khoản T có email; máy chủ trả lỗi hệ thống

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Nhập tên đăng nhập và email của T, bấm *Tiếp tục* (6) | Báo "Lỗi hệ thống" (5); vẫn ở màn "Quên mật khẩu" |

## TC-AUTH-F02-007 · Nút Tiếp tục khi đang chờ phản hồi
- Truy vết: AUTH-F02-US01
- Ưu tiên: P3
- Màn hình: AUTH-F02-SCR01
- Tiền điều kiện: tài khoản T có email; máy chủ phản hồi chậm

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Nhập tên đăng nhập và email của T, bấm *Tiếp tục* (6) | Trong lúc chờ, nút 6 hiện "Đang kiểm tra…" và không bấm được |
| 2 | Chờ máy chủ phản hồi | Sang màn "Đặt mật khẩu mới" |

## TC-AUTH-F02-008 · Bỏ trống ô mật khẩu
- Truy vết: AUTH-F02-US02-A1
- Ưu tiên: P2
- Màn hình: AUTH-F02-SCR02
- Tiền điều kiện: tài khoản T đã xác minh xong, đang ở màn "Đặt mật khẩu mới"
- Dữ liệu kiểm tra:

| Mật khẩu mới | Nhập lại |
|---|---|
| (trống) | (trống) |
| (trống) | "Moi@12345" |
| "Moi@12345" | (trống) |

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Nhập theo dữ liệu, bấm *Đặt lại mật khẩu* (6) | Báo "Nhập mật khẩu mới và nhập lại mật khẩu" (5); vẫn ở màn "Đặt mật khẩu mới" |
| 2 | Đăng nhập bằng tên đăng nhập của T và "Test@1234" | Vào trang chủ (mật khẩu không đổi) |

## TC-AUTH-F02-009 · Mật khẩu ngắn hơn 8 ký tự
- Truy vết: AUTH-F02-US02-A2, AUTH-BR04
- Ưu tiên: P2
- Màn hình: AUTH-F02-SCR02
- Tiền điều kiện: tài khoản T đã xác minh xong, đang ở màn "Đặt mật khẩu mới"

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Nhập "Ab@1234" (7 ký tự) vào cả hai ô (3), (4), bấm *Đặt lại mật khẩu* (6) | Báo "Mật khẩu phải có ít nhất 8 ký tự" (5); vẫn ở màn "Đặt mật khẩu mới" |

## TC-AUTH-F02-010 · Mật khẩu dài hơn 16 ký tự
- Truy vết: AUTH-F02-US02-A5, AUTH-BR04
- Ưu tiên: P2
- Màn hình: AUTH-F02-SCR02
- Tiền điều kiện: tài khoản T đã xác minh xong, đang ở màn "Đặt mật khẩu mới"

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Nhập "Abcdefgh@12345678" (17 ký tự) vào cả hai ô (3), (4), bấm *Đặt lại mật khẩu* (6) | Báo "Mật khẩu tối đa 16 ký tự" (5); vẫn ở màn "Đặt mật khẩu mới" |

## TC-AUTH-F02-011 · Độ dài mật khẩu hợp lệ ở biên
- Truy vết: AUTH-F02-US02, AUTH-BR04
- Ưu tiên: P2
- Màn hình: AUTH-F02-SCR02, AUTH-F01-SCR01
- Tiền điều kiện: tài khoản T đã xác minh xong, đang ở màn "Đặt mật khẩu mới"
- Dữ liệu kiểm tra:

| Mật khẩu mới (nhập vào cả hai ô) | Ghi chú |
|---|---|
| "Abcd@123" | đúng 8 ký tự |
| "Abcdefgh@1234567" | đúng 16 ký tự |
| "Mật khẩu 12" | 11 ký tự, có dấu tiếng Việt và khoảng trắng |
| "Test@1234" | trùng mật khẩu hiện tại |

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Nhập mật khẩu theo dữ liệu vào ô (3) và (4), bấm *Đặt lại mật khẩu* (6) | Về trang đăng nhập, hiện "Đặt lại mật khẩu thành công. Đăng nhập bằng mật khẩu mới." |
| 2 | Đăng nhập bằng tên đăng nhập của T và mật khẩu vừa đặt | Vào trang chủ |

## TC-AUTH-F02-012 · Nhập lại mật khẩu không khớp
- Truy vết: AUTH-F02-US02-A3
- Ưu tiên: P2
- Màn hình: AUTH-F02-SCR02
- Tiền điều kiện: tài khoản T đã xác minh xong, đang ở màn "Đặt mật khẩu mới"

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Nhập *Mật khẩu mới* (3) = "Moi@12345", *Nhập lại* (4) = "Moi@12346", bấm *Đặt lại mật khẩu* (6) | Báo "Mật khẩu nhập lại không khớp" (5); vẫn ở màn "Đặt mật khẩu mới" |

## TC-AUTH-F02-013 · Phiên đặt lại hết hiệu lực khi bấm Đặt lại (xác minh lại ở tab khác)
- Truy vết: AUTH-F02-US02-A4, AUTH-BR06
- Ưu tiên: P2
- Màn hình: AUTH-F02-SCR02, AUTH-F02-SCR01
- Tiền điều kiện: tài khoản T có email; chưa đăng nhập

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Tab A: xác minh T ở màn "Quên mật khẩu" | Tab A ở màn "Đặt mật khẩu mới" |
| 2 | Tab B (cùng trình duyệt): xác minh T lần nữa | Tab B ở màn "Đặt mật khẩu mới" |
| 3 | Tab A: nhập "Ab@1" vào cả hai ô (quá ngắn), bấm *Đặt lại mật khẩu* (6) | Tab A về màn "Quên mật khẩu", báo "Phiên xác minh đã hết hạn, vui lòng xác minh lại" — không báo lỗi độ dài |
| 4 | Đăng nhập bằng tên đăng nhập của T và "Test@1234" | Vào trang chủ (mật khẩu không đổi) |

## TC-AUTH-F02-014 · Phiên đặt lại hết hiệu lực khi tài khoản bị xoá
- Truy vết: AUTH-F02-US02-A4, AUTH-BR06
- Ưu tiên: P3
- Màn hình: AUTH-F02-SCR02, AUTH-F02-SCR01
- Tiền điều kiện: tài khoản T đã xác minh xong, đang ở màn "Đặt mật khẩu mới"; sau đó T bị xoá

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Nhập "Moi@12345" vào cả hai ô, bấm *Đặt lại mật khẩu* (6) | Về màn "Quên mật khẩu", báo "Phiên xác minh đã hết hạn, vui lòng xác minh lại" |

## TC-AUTH-F02-015 · Bấm Về đăng nhập ở bước 2 thì huỷ phiên đặt lại
- Truy vết: AUTH-F02-US02-A6, AUTH-BR06
- Ưu tiên: P2
- Màn hình: AUTH-F02-SCR02, AUTH-F01-SCR01, AUTH-F02-SCR01
- Tiền điều kiện: tài khoản T đã xác minh xong, đang ở màn "Đặt mật khẩu mới"

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Bấm *< Về đăng nhập* (7) | Về trang đăng nhập, không có thông báo thành công |
| 2 | Bấm nút Back của trình duyệt (quay lại màn "Đặt mật khẩu mới") | Về màn "Quên mật khẩu", báo "Phiên xác minh đã hết hạn, vui lòng xác minh lại" |
| 3 | Đăng nhập bằng tên đăng nhập của T và "Test@1234" | Vào trang chủ (mật khẩu không đổi) |

## TC-AUTH-F02-016 · Mất kết nối máy chủ khi đặt lại
- Truy vết: AUTH-F02-US02-E1
- Ưu tiên: P3
- Màn hình: AUTH-F02-SCR02
- Tiền điều kiện: tài khoản T đã xác minh xong, đang ở màn "Đặt mật khẩu mới"; sau đó máy chủ không phản hồi

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Nhập "Moi@12345" vào cả hai ô, bấm *Đặt lại mật khẩu* (6) | Báo "Không kết nối được máy chủ" (5); vẫn ở màn "Đặt mật khẩu mới" |

## TC-AUTH-F02-017 · Lỗi hệ thống khi đặt lại
- Truy vết: AUTH-F02-US02-E2
- Ưu tiên: P3
- Màn hình: AUTH-F02-SCR02
- Tiền điều kiện: tài khoản T đã xác minh xong, đang ở màn "Đặt mật khẩu mới"; sau đó máy chủ trả lỗi hệ thống

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Nhập "Moi@12345" vào cả hai ô, bấm *Đặt lại mật khẩu* (6) | Báo "Lỗi hệ thống" (5); vẫn ở màn "Đặt mật khẩu mới" |

## TC-AUTH-F02-018 · Đặt lại mật khẩu làm phiên đăng nhập cũ hết hiệu lực
- Truy vết: AUTH-F02-US02, AUTH-BR07
- Ưu tiên: P2
- Màn hình: AUTH-F02-SCR02, AUTH-F01-SCR02, AUTH-F01-SCR01
- Tiền điều kiện: tài khoản T có email; T đang đăng nhập ở trình duyệt X (đang ở trang chủ)

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Ở trình duyệt Y (chưa đăng nhập): xác minh T, đặt mật khẩu mới "Moi@12345" | Trình duyệt Y về trang đăng nhập, hiện thông báo thành công |
| 2 | Ở trình duyệt X: tải lại trang chủ | Chuyển về trang đăng nhập |

## TC-AUTH-F02-019 · Mở màn Đặt mật khẩu mới khi chưa xác minh
- Truy vết: AUTH-F02-SCR02
- Ưu tiên: P3
- Màn hình: AUTH-F02-SCR02, AUTH-F02-SCR01
- Tiền điều kiện: chưa đăng nhập; chưa xác minh

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Mở thẳng địa chỉ của màn "Đặt mật khẩu mới" | Về màn "Quên mật khẩu", không có thông báo lỗi |

## TC-AUTH-F02-020 · Mở lại màn Đặt mật khẩu mới sau khi đã đặt lại xong
- Truy vết: AUTH-F02-SCR02, AUTH-BR06
- Ưu tiên: P3
- Màn hình: AUTH-F02-SCR02, AUTH-F02-SCR01
- Tiền điều kiện: tài khoản T vừa đặt lại mật khẩu thành công (phiên đặt lại đã dùng)

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Mở lại địa chỉ của màn "Đặt mật khẩu mới" | Về màn "Quên mật khẩu", báo "Phiên xác minh đã hết hạn, vui lòng xác minh lại" |

## TC-AUTH-F02-021 · Người đang đăng nhập mở màn Quên mật khẩu
- Truy vết: AUTH-F02-SCR01
- Ưu tiên: P3
- Persona: user1
- Màn hình: AUTH-F02-SCR01, AUTH-F01-SCR02
- Tiền điều kiện: user1 đang đăng nhập

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Mở thẳng địa chỉ của màn "Quên mật khẩu" | Chuyển về trang chủ, hiện "Xin chào Nguyễn Văn A" |
| 2 | Mở thẳng địa chỉ của màn "Đặt mật khẩu mới" | Chuyển về trang chủ |

## TC-AUTH-F02-022 · Xác minh sai không để lộ tài khoản (API)
- Truy vết: AUTH-F02-US01-A2, AUTH-BR02
- Ưu tiên: P2
- Tầng: api
- Tiền điều kiện: tài khoản T có email

| # | Bước | Kết quả mong đợi |
|---|---|---|
| 1 | Gửi xác minh với tên không tồn tại + email bất kỳ | Bị từ chối với thông báo "Tên đăng nhập hoặc email không đúng" |
| 2 | Gửi xác minh với tên của T + email sai | Bị từ chối, cùng mã trạng thái và cùng nội dung lỗi như bước 1 |
