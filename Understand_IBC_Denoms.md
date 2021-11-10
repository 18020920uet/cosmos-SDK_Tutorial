# Understand IBC Denoms

---

## Understand IBC Denoms with Gaia

Một trong những công nghệ mạnh mẽ nhất khi sử dụng Cosmos SDK là Giao thức truyền thông chuỗi liên kết (IBC).
Trong hệ sinh thái Cosmos, mọi blockchain đều hướng tới việc có chủ quyền và đối với từng ứng dụng cụ thể.
Với IBC, mọi blockchain có thể kết nối với một blockchain khác bằng giao thức IBC.
Giao thức truyền thông này cuối cùng sẽ tạo ra một hệ thống các blockchain có chủ quyền và được kết nối.

### Introduction

Tính năng được sử dụng nhiều nhất của IBC là gửi _token_ từ blockchain này sang blockchain khác.
Khi 1 _token_ được gửi từ blockchain khác, một _token chứng từ_ được tạo ra trên blockchain nhận.

Hãy tưởng tượng hai blockchain, blockchain A và blockchain B. Ban đầu, bạn có _token_ của mình trên blockchain A.

![Gưi token từ blockchain A -> blockchain B](https://tutorials.cosmos.network/assets/img/ibc_token.161a464a.png)

Giá trị mà token đại diện có thể được chuyển qua các blockchain, nhưng bản thân token thì không thể. Khi gửi các token bằng IBC đến một blockchain khác:

1. Blockchain A locks the tokens and relays proof to blockchain B

  (Blockchain A khóa các token và chuyển tiếp bằng chứng cho blockchain B)

2. Blockchain B mints its own representative tokens in the form of voucher replacement tokens

  (Blockchain B sẽ đúc _token_ của riêng nó dưới dạng _token chứng từ_ thay cho token của A)
3. Blockchain B sends the voucher tokens back to blockchain A

  (Blockchain B sẽ gửi các _token chứng từ_ trở lại blockchain A)
4. The voucher tokens are destroyed (burned) on blockchain B

  (_Token chứng từ_ ở blockchain B sẽ bị phá hủy)

5. The locked tokens on blockchain A are unlocked

  (_Token_ bị khóa ở blockchain A sẽ được mở khóa)

Cách duy nhất để mở khóa _token_ bị khóa trên blockchain A là gửi lại _token chứng từ_ từ blockchain B.
Kết quả là _token chứng từ_ trên blockchain B bị phá hủy.
Việc hủy _token chứng từ_ có mục đích chính là đẩy _token_ đó ra khỏi việc lưu thông (trao đổi)

Trong hướng dẫn này, bạn sẽ học định dạng của _token chứng từ_ ở blockchain B.
Bạn sẽ tìm hiểu về phần thông tin mà _token chứng từ_ chứa và _token chứng từ_ có cấu trúc như thế nào và sẽ hiểu thêm về chúng.
Thông tin của _token_ được mô tả như những mệnh giá IBC.
Bạn có thể phân tích mệnh giá IBC này để nhận thông tin về chứng từ và tìm hiểu được _token chứng từ_ này đến từ blockchain nào.

Bạn sẽ học:
- Trao đổi các mệnh giá IBC
- Hiểu cách hoạt động của mệnh giá
- Tìm hiểu được từ blockchain nào mà _token_ được tạo ra

### Requirements

Cài đặt file nhị phân **gaia**

```bash
git clone https://github.com/cosmos/gaia.git
cd gaia
git checkout v5.0.0
make install

gaiad version #-> v5.0.0
```

### What Is This IBC Denom

_Token chứng từ_ được gới thiệu trong quá trình chuyển đồi tài nguyên là **các mệnh giá IBC** (IBC denom).
_Token chứng từ_ là kết quả của việc chuyển đôi 1 _token_ bằng IBC từ blockchain này đến blockchain khác.
Cấu trúc của 1 _token chứng từ_ như sau:

> ibc/DENOMHASH

Tưởng tượng rằng bạn nhận 1 _/ibc token_ trên blockchain B nơi mà giữ **samoleans token** và **stake token**

Số dư hiện tại sẽ như sau:

> 1000000ibc/CDC4587874B85BEA4FCEC3CEA5A1195139799A1FEE711A07D972537E18FDA39D,100000000000samoleans,99999977256stake

Giống như **samoleans** và **stakte**, _ibc/CDC..._ là mệnh giá của _token_ nhận từ IBC.
Sau khi _ibc/CDC.._ là mã hóa của mệnh giá, phần cổng IBC, và phần kênh.

Tại sao CDC4... là một mã hóa

- Phần mã hóa bao gồm đường dẫn để có thể truy vết _token_ trên nhiều bước nhảy từ các blockchain khác đến tài khoản của bạn
- Phần đường dẫn này có thể cực dài khi trực tiếp thể hiện ra
- Cosmos SDK có max là 64 kí tự trên phần mệnh giá của _token_

Việc đánh đổi bằng cách sử dụng mã hóa là bạn phải query 1 node để tìm ra phần đường dẫn và mệnh giá của 1 token.
Việc này gọi là **denomtrace** (truy vết mệnh giá).

Cùng với câu lệnh **gaiad** dưới đây để query mệnh giá và học thêm về phần kênh mà _token_ đến

```bash
gaiad query ibc-transfer denom-trace CDC4587874B85BEA4FCEC3CEA5A1195139799A1FEE711A07D972537E18FDA39D --node https://rpc.testnet.cosmos.network:443

---

## Phản hồi:
denom_trace:
  base_denom: moon
  path: transfer/channel-14
```

Với kết quả phản hồi của câu lệnh, bạn sẽ thấy rằng có 1 cổng chuyển đồi IBC và kênh
**channel-14**.
Để biết tại sao **IBC light client** đằng sau cồng và kênh, bạn phải thực hiện 1 query khác

Tại sao lại gọi là **IBC light client**? Bời vì nó là 1 **light client** của blockchain khác, theo dõi liên tục phần khối mã hóa của nó.
Lệnh **ibc channel client-state transfer** giải thích chi tiết phần đường dẫn của mệnh giá

```bash
gaiad query ibc channel client-state transfer channel-14 --node https://rpc.cosmos.network:443

---

## Phản hồi:
client_id: 07-tendermint-18
client_state:
  '@type': /ibc.lightclients.tendermint.v1.ClientState
  allow_update_after_expiry: true
  allow_update_after_misbehaviour: true
  chain_id: mars
  frozen_height:
    revision_height: "0"
    revision_number: "0"
  latest_height:
    revision_height: "2207"
    revision_number: "0"
  max_clock_drift: 600s
  proof_specs:
  - inner_spec:
      child_order:
      - 0
      - 1
      child_size: 33
      empty_child: null
      hash: SHA256
      max_prefix_length: 12
      min_prefix_length: 4
    leaf_spec:
      hash: SHA256
      length: VAR_PROTO
      prefix: AA==
      prehash_key: NO_HASH
      prehash_value: SHA256
    max_depth: 0
    min_depth: 0
  - inner_spec:
      child_order:
      - 0
      - 1
      child_size: 32
      empty_child: null
      hash: SHA256
      max_prefix_length: 1
      min_prefix_length: 1
    leaf_spec:
      hash: SHA256
      length: VAR_PROTO
      prefix: AA==
      prehash_key: NO_HASH
      prehash_value: SHA256
    max_depth: 0
    min_depth: 0
  trust_level:
    denominator: "3"
    numerator: "1"
  trusting_period: 1209600s
  unbonding_period: 1814400s
  upgrade_path:
  - upgrade
  - upgradedIBCState
```

Có rất nhiều thông tin nhưng nó không trả lời được câu hỏi: Làm sao bạn biết rằng có thể tin tường khách hàng IBC này?


#### The Chain ID and the Client ID

Bất kỳ ai cũng có thể chạy một blockchain với cùng chain ID. Nhưng IBC client ID được tạo ra bởi [mô-đun Cosmos SDK IBC Keeper](https://github.com/cosmos/ibc-go/blob/e012a4af5614f8774bcb595962012455667db2cf/modules/core/02-client/keeper/keeper.go#L56) (ISC-02 không xác định 1 chuẩn cho IBC client IDs).
Các dịch vụ tên blockchain và kho dằng ký blockchain Github không phân cấp có thể xác minh sự kết hợp của chain ID và client ID.
Tên dịch vụ blockchain và repo công ty đăng ký blockchain đều đang được phát triển và được coi là thử nghiệm.  

#### Ensure the IBC Client Isn't Expired

Trong trường hợp không có sự đồng thuận của Tendermint (nếu > 1/3 trình xác thực tạo ên 1 khối xung đột), và bằng chứng về sự thât bại của sự đồng thuận này được gửi trục tiếp, IBC client sẽ bị đóng bằng với **frozen_height** khác không.
Trong ví dụ trước, kết quả của truy vấn **gaiad** kênh ibc trạng thái máy khách xác nhận rằng trạng thái máy khác và biết được rằng máy khách IBC chưa hết thời hạn.

Last_height.revision_height là chiều cao khối khi IBC client được cập nhật lần cuối.
Để đảm bảo rằng chiều cao khối vẫn được cập nhật, bạn sẽ phải truy vấn chính blockchain về chiều cao khối 2207 và đảm bảo rằng dấu thời gian của khối đó + thời gian tin cậy của 1209600s/336h/14d là sau thời điểm hiện tại.


Ví dụ về xác thực trạng thái của IBC client bằng truy vấn:

```bash
gaiad query block 5200792 --node https://rpc.cosmos.network:443
```

### Find the Path of Another Blockchain

Việc có thể liệt kê hết các đường dẫn blockchain vẫn là một vấn đề chưa được giải quyết

Bất kỳ blockchain được tạo nào cũng có thể tạo một kênh đối với blockchain khác mà không để lộ quá nhiều thông tin.

Hiện tại, các kênh phải trao đổi với nhau bằng giao thức giữa người với người để có thể coi là đáng tin cậy.
Giao thức giao tiếp giữa người với người này sử dụng mệnh giá IBC đề người dùng có thể xác định _token_ nào có thể chấp nhận bởi ứng dụng.

Một cách tiếp cận để giải quyết vấn đề này là sử dụng cơ sở dữ liệu tập trung hoặc phi tập chung của các blockchain ID và các node của nó.
Có 2 giải pháp đang được phát triển:

- Tên của các dịch vụ blockchain (Chain Name Service), phi tập trung:

  [CNS](https://github.com/tendermint/cns) với mục tiêu trở thành 1 mô-đun của Cosmos SDK mà Trung tâm Cosmos một ngày nào đó sẽ chạy.
  Là một trung tâm thông qua đó các giao dịch chuỗi chéo sẽ diễn ra, nó chỉ hợp lý đối với việc một ngày nào đó Trung tâm Cosmos lưu giữ thông tin quan trọng về cách tiếp cận các blockchain ID.
  CNS vẫn đang mới và đang được phát triển.

- Cơ quan đăng ký Blockchain (Cosmos Registry), bán phi tập trung:

  Repo **github.com/cosmos/registry** là 1 giải pháp tạm thời.
  Mỗi một  chain ID có 1 thư mục mô tả genesis và 1 danh sách các đối tác.
  Để có thể lấy chain ID, ứng dụng blockchain sẽ phải 'fork' repo đăng ký đó, tạo 1 nhánh với chain ID, và gửi pull request để để thêm chain ID vào trong danh sách ID chính thức của **cosmos/registry**

  Với mỗi chain ID được thể hiện bới 1 thư mục, trong thư mục đó có file **peers.json** bao gồm danh sách các node mà người dùng có thể kế nối dến

- Công cụ đăng ký Cosmos (Cosmos Registrar), bán phi tập trung:
  [Cosmos-registrar](https://github.com/apeunit/cosmos-registrar) là một công cụ được khơi mào bởi Jack Zampolin và được phát triển tiếp bới Ape Unit.
  **Cosmos-registrar** tự động xác nhận và cập nhật chain ID.
  Trong trường hợp này, việc cập nhật chain ID có nghĩa là 'commit' một danh sách các đối tác mới cho kho lưu trữ Github.
  'Commit' này sẽ chạy với 1 **cronjob**.
