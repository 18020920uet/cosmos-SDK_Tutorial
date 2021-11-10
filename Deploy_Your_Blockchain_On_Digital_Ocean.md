# Deploy Your Blockchain on Digital Ocean

---

## Publish Your Starport Blockchain App on DigitalOcean

Trong hướng dẫn này, bạn sẽ học cách công bố ứng dụng blockchain trên hệ điều hành Ubuntu 20.04 _droplet_ trên DigitalOcean.

### Goals
Bạn sẽ làm theo các bước dưới đây để xây và chạy ugnws dụng blockchain với các bước từ thấp đến cao sau.
- Tạo 1 ứng dụng blockchain Cosmos SDK với Starport
- Tạo 1 Droplet trên DigitalOcean
- Đấy ứng dụng lên Droplet
- Sử dụng ứng dụng trên Droplet

### Introduction

Bài viết này viết cho DigitalOcean.
Với các thông tin này có thể rất hữu ích với các dịch vụ đám mây khác, phần chi tiết chỉ áp dụng cho DigitalOcean Droplets.
Vì DigitalOcean Droplet được dựa trên các máy áo dùng hệ điều hành Linux, các hệ điều hành trong hướng dẫn này chỉ áp dụng cho MacOS và Linux.

Các bước có thể khác với hệ điều hành Window.

**Quan trọng!. Giá trị của cac placeholder được sử dụng trong suốt hướng dẫn này**:

    Khi bạn thấy [username] hãy đổi thành tên của bạn
    Bạn có thể thay đổi tên của thư mục ứng dụng: [digitalocean] bằng tên của ứng dụng của bạn.


### Requirements
- Một tài khoản DigitalOcean
- Starport
- Go

### Step 1 — Create Your Starport App

Trong bước này, bạn sẽ tạo 1 ứng dụng blockchain với Starport trên máy tính của bạn.
Tiếp theo bạn sẽ đẩy phần ứng dụng này lên DigitalOcean Droplet

Tạo một thư mục và trong thư mục đó dùng Starport để tạo 1 ứng dụng blockchain
```bash
mkdir blockchain && cd blockchain
starport app github.com/<username>/<application>
```

Một ứng dụng blockchain với phần backend và frontend được tạo ra. Việc dựng mặc định sẽ tạo 1 ứng dụng blockcahin với các chức năng tương đồng với [Trung tâm Cosmos](https://hub.cosmos.network/main/hub-overview/overview.html) [($ATOM)](https://cosmos.network/features/).

Để thay đổi bất cứ tham số nào cảu blockchain thì có xem thể cài đặt ở [đây](https://docs.starport.network/kb)

### Step 2 — Host Your App on DigitalOcean

DigitalOcean cung cấp cho các nhà phát triển dịch vụ đám mây có khả năng deploy và mở rộng ứng dụng chạy đồng thời trên các máy tính

Sử dụng cơ sở hạ tầng của DigitalOcean để công bố 1 ứng dụng blockchain Cosmos SDK để đảm bảo rằng ứng dụng có thể truy cập từ web. Ứng dụng của bạn phải có thể truy cập từ web để có thể cho phép các ứng dụng khác hoặc các node khác kết nối đến mạng blockchain của bạn.

#### Create a Droplet on DigitalOcean

1. Đăng nhập vào DigitalOcean và vào [Bảng điều kiên](https://cloud.digitalocean.com/login)

2. Từ nút **Create** ở trên cùng bên phải của **Bảng điều kiển**, chọn **Droplet**

3. Chọn Ubuntu 20.04 (LTS) x64 Gói CPU dùng chung cơ bản

  ![Ảnh](https://tutorials.cosmos.network/assets/img/do-ubuntu-image.59b08bdf.png)

4. Chọn bộ nhớ và các tài nguyên để tính toán.
Máy ảo cung cấp 1 số chọn lựa trộn lẫn. Với 1 ứng dụng blockchain điển hình mà không có dùng với mục đích tính toán nhiều, sử dụng 2 GB Ram / 1 CPU là môi trường thể phù hợp để phát triển và test.

  ![Ảnh](https://tutorials.cosmos.network/assets/img/do-droplet-size.eb59960c.png)

5. Xác nhận khu vực lưu trữ dữ liệu và lựa chọn mạng VPC

6. Chúng tôi gợi ý sử dụng SSH keys đẻ cho việc xác thực khi truy cập vào Droplet.

7. Với câu hỏi **How many Droplets?** Chọn mặc định 1 Droplet

8. Vơi câu hởi **Choose a hostname** bạn có thể tự lựa chọn và cài đặt tên.

9. Với **Tags** bạn có thể tùy chọn chọn lưa thêm các thẻ tag

10. **Select Project** bạn có thể sử dụng project mặc định

11. Không chọn **Add backups**

12. Xác nhận việc tạo Droplet **Create Droplet**

Droplet có thể tìm thấy ở danh sách Droplets của bạn

#### Access Your New Node

1. Trong danh sách Droplet của bạn trong **Control Panel**, bạn có thể di chuột vào để thấy thông tin địa chỉ IP và sao chép lại.

2. Mở terminal và nhập như dưới để có thể truy cập vào node. Có thể xem hướng dẫn tại [đây](https://docs.digitalocean.com/products/droplets/how-to/connect-with-ssh/openssh/)
  ```bash
  ssh root@<IP address>
  ```

#### Create a User to Run the Blockchain

Đầu tiên hãy tạo user để chạy ứng dụng blockchain.
Trong ví dụ này bạn có thể sử dụng **user** như là người dùng.
Thêm người dụng bằng cách sau:

```bash
adduser <user>
```

Sau khi nó chạy xong, nhập và xác nhận mật khẩu cho **user**

Sau khi mật khẩu được cập nhật, bạn có thê yêu cầu thay đổi thông tin cho **user**.
Ấn ENTER để sử dụng cá thông tin mặc định.
Ấn Y để xác nhận rằng các thông tin đó chính xác.

Đẻ có thể thêm quyền phù hợp với người dùng, thêm **user** và group sudo

```bash
usermod -aG sudo <user>
```

Để có thể đăng nhập bằng tài khoản **user**, sao chép _ssh public key_ đến thư mục _home_ mới của **user**

```bash
rsync --archive --chown=<user>:<user> ~/.ssh /home/<user>
```

Chuyển sang tài khoản **user**

```bash
su - <user>
```

#### Set Up Your Droplet

1. Cài đặt Go

  Phiên bản sử dụng là v1.16.5

  ```bash
  cd ~
  wget https://golang.org/dl/go1.16.5.linux-amd64.tar.gz
  sha256sum go1.16.5.linux-amd64.tar.gz
  # b12c23023b68de22f74c0524f10b753e7b08b1504cb7e417eccebdd3fae49061  go1.16.5.linux-amd64.tar.gz
  ```

2. Giải nén và lưu trữ nó trong thư mục _/usr/local_:

  ```bash
  sudo tar -xvf go1.16.5.linux-amd64.tar.gz -C /usr/local_
  ```

  Sẽ có 1 thư mục _go_ trong _/usr/local_: **_/usr/local/go_**

3. Chuyển quyền sở hữu cho **user**
  ```bash
  sudo chown -R <user>:<user> /usr/local/go
  ```

4. Chuyển đường dẫn của Go vào _~/.**profile**_

  Mở file _~/.**profile**_
  ```bash
  nano ~/.profile
  ```
  Thêm 2 dòng này vào cuối file

  ```bash
  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin
  ```

5. Cập nhật lại profile
  ```bash
  source ~/.profile
  ```

6. Kiểm tra lại file go

  ```bash
  go version
  # Result: go version go1.16.5 linux/amd64
  ```

7. Cài đặt Starport phiên bản mới nhất

  ```bash
  curl https://get.starport.network/starport! | sudo bash
  ```

### Upload Your App to the Droplet

Có 1 vài cách để có thể đẩy ứng dụng của bạn lên Droplet

Một các tiên lọi nhất là sử dụng Github và tải nó về ở Droplet.
Starport tạo 1 Github Repo mà có thể giúp tiện lợi để thêm ứng dụng vào kho lưu trữ của bạn.

Để upload các file lên thì cbạn có thể dùng **scp** để bỏ qua 1 số bước không cần thiết.

Trong terminal của máy tính của bạn, đảm bảo rằng bạn đang ở thư mục cha của thưc mục repo mà được tạo bởi Starport tư lúc nãy.

Ví dụ thư mục /home/blockchain là thự mục cha của ứng dụng của bạn
Hãy sử dụng câu lệnh dưới đây từ thư mục đó để upload phần dữ liệu lên:

```bash
# /home/blockchain
scp -r digitalocean <user>@<IP address>:/home/<user>/
```

Câu lệnh này sẽ upload folder ứng dụng của bạn lên đến Droplet của DigitalOcen ở thư mục
/home/[user]

### Install your app

1. Dùng terminal của Droplet chuyển sang thư mục /home/[user]/[application] sau khi đã upload thành công
  ```bash
  cd <application>
  ```
2. Để chạy và cài đặt ứng dụng trên Droplet, sử dụng lệnh CLI
  ```bash
  starport serve
  ```

Câu lệnh trễn sẽ giúp build ứng dụng Starport và tạo file genesis và file cài đặt.

Đây có thể là một ví dụ về kết quả của quá trình này:

```
Cosmos SDK's version is: Stargate v0.40.0 (or later)

🛠️  Building proto...
📦 Installing dependencies...
🛠️  Building the blockchain...
💿 Initializing the app...
🙂 Created account "alice" with address "cosmos1sucepwrvgkud7fc6ftne8s2glzjpqx4zsl6zxa" with   mnemonic: "expect goddess business detail loud know broom trial deliver board victory   despair tackle ripple body weapon runway lawn roast cactus attitude midnight town fox"
🙂 Created account "bob" with address "cosmos15jayucugyfnlj4tu7xhutnufxr4qd8j6qjmekk" with   mnemonic: "mouse lamp excuse young top century empower afford oven grass pass heavy evil   sample lake trick leisure aisle bird dumb radio learn ecology stamp"
Genesis transaction written to "/home/appuser/.digitalocean/config/gentx/gentx-  85ef06cb86aa501cceb9ed0a497a02503f5aa57f.json"
🌍 Tendermint node: http://0.0.0.0:26657
🌍 Blockchain API: http://0.0.0.0:1317
🌍 Token faucet: http://0.0.0.0:4500
```

Để tắt ứng dụng có thể sử dụng **Ctrl + C**



Để xác thực việc đã cài đặt thành công dùng lệnh applicationd ("application+d" ), và phần kết quả đã được mô tả ở dưới:

```bash
<applicationd>

# ----

# Stargate CosmosHub App
#
# Usage:
#   <applicationd> [command]
#
# Available Commands:
#
#
#   add-genesis-account Add a genesis account to genesis.json
#   collect-gentxs      Collect genesis txs and output a genesis.json file
#   debug               Tool for helping with debugging your application
#   export              Export state to JSON
#   gentx               Generate a genesis tx carrying a self delegation
#   help                Help about any command
#   init                Initialize private validator, p2p, genesis, and application # configuration files
#   keys                Manage your application's keys
#   migrate             Migrate genesis to a specified target version
#   query               Querying subcommands
#   start               Run the full node
#   status              Query remote node for status
#   tendermint          Tendermint subcommands
#   tx                  Transactions subcommands
#   unsafe-reset-all    Resets the blockchain database, removes address book files, and resets # data/priv_validator_state.json to the genesis state
#   validate-genesis    validates the genesis file at the default location or at the location # passed as an arg
#   version             Print the application binary version information
#
# Flags:
#   -h, --help                help for <applicationd>
#       --home string         directory for config and data (default "/home/appuser/# .digitalocean")
#       --log_format string   The logging format (json|plain) (default "plain")
#       --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) # (default "info")
#       --trace               print out full stack trace on errors
#
# Use "<applicationd> [command] --help" for more information about a command.
```

### Create a Systemd process

Câu lệnh **starport serve** giúp để cài đặt, xậy dựng và cấu hình ứng dụng của bạn cho một blockchain mới. Bạn sẽ nhận _token_ khi 1 tài khoản được tạo mới.

Ứng dụng mới này sẽ chỉ giúp ích cho mục đích phát triển và không hướng tới mục đích chạy như 1 ứng dụng thực tế.

Để chạy ứng dụng với 1 thời gian dài hơn, bạn sẽ phải sử dụng tiến trình **systemd** trên Droplet của bạn.

Tạo các file log và file **systemd** để quản lý tiến trình mà từ đó bạn có thể tắt terminal đi và có thể tụ tin rằng ứng dụng sẽ chạy kể cả khi server khởi động lại.

Sẽ gọi file sau khi được build là [applcationd]

```bash
sudo mkdir -p /var/log/<applicationd>
sudo touch /var/log/<applicationd>/<applicationd>.log
sudo touch /var/log/<applicationd>/<applicationd>_error.log
sudo touch /etc/systemd/system/<applicationd>.service
```

Dùng **nano** để sửa file _system sevice_:

```bash
sudo nano /etc/systemd/system/<applicationd>.service
```

Copy phần cấu hình dưới đây vào

```txt
[Unit]
Description=digitalocean daemon
After=network-online.target
[Service]
User=appuser
ExecStart=/home/appuser/go/bin/<applicationd> start --home=/home/appuser/.digitalocean
WorkingDirectory=/home/appuser/go/bin
StandardOutput=file:/var/log/<applicationd>/<applicationd>.log
StandardError=file:/var/log/<applicationd>/<applicationd>_error.log
Restart=always
RestartSec=3
LimitNOFILE=4096
[Install]
WantedBy=multi-user.target
```

Cho phép tiến trình hoạt động để ứng dụng có thể chạy sau khi server bị khởi động lại

```bash
sudo systemctl enable <applicationd>.service
```

Khởi chạy tiến trình

```bash
sudo systemctl start <applicationd>.service
```

Kiểm tra nhật ký với lệnh

```bash
sudo journalctl -u <applicationd> -f
```

### Connect your local running chain to the published chain

Để kết nối với chuỗi nội bộ của bạn với ứng đụng đã được phát hành thì có một số yêu cầu sau:
- Tìm node mà bạn vừa tạo với địa chỉ IP
- Tải file genesis
- Cài đặt ứng dụng của bạn để có thể kết nối với node đã được phát hành

Trên DigitalOcean Droplet của bạn, tìm mã hóa của node bằng cách:

```bash
<applicationd> tendermint show-node-id
# For example, 85ef06cb86aa501cceb9ed0a497a02503f5aa57f
```

Mỗi node có 1 id khác nhau.

Tại terminal ở máy nội bộ của bạn, khởi tạo ứng dụng blockchain với câu lệnh:

```bash
<applicationd> init localnode
```

Tiếp theo là việc tải file genesis.json về

```bash
scp <user>@<IP address>:/home/<user>/.<applcation>/config/genesis.json ~/.<application>/config/
```

Bây giờ hãy thêm ID của node đã được công bố và IP server vào file cấu hình trong máy nội bộ của bạn
```bash
nano ~/.<application>/config/config.toml
````

Tạo 1 liên bền. Sử dụng node-id từ Droplet với địa chỉ IP và cổng mặc định:

```bash
# Comma separated list of nodes to keep persistent connections to
persistent_peers = "<node-id>@<IP address>:<port>"

# Ví dụ
# persistent_peers = "85ef06cb86aa501cceb9ed0a497a02503f5aa57f@<IP address>:26656"
```

Khởi chạy node:
```bash
<applicationd> start
```

Bạn có thể sử dụng một công cụ giống _gex_ để kiểm tra xem bạn đã kết nối với 1 node chưa.
Bạn có thể cài đặt gẽ trên cả remote và máy local. Chạy câu lệnh dưới đây trên cả 2 máy:

```bash
go get -u github.com/cosmos/gex
```
Sau khi cài đặt, bắt đầu sử dụng gex với câu lệnh

```bash
gex
```

Khi có ít nhất 1 kết quả ở phần _Connected Peers_ thì mọi thứ đã được cài đặt thành công


### Access API and RPC

Bây giờ bạn có thể truy cập API bằng dường dẫn

> http://<IP address>:1317

Vả RPC ở cổng 26657:

> http://<IP address>:26657

### Optional: Set Up a Firewall on Your Droplet

Đẻ giảm tải việc tải và tấn công đến droplet bạn có thể sử dụng tường lửa.

Mặc dù DigitalOcean đã cung câp 1 giao điện UI cho việc tạo tường lửa thì hướng dẫn này sẽ chỉ ra 1 cách phổ thông để tiếp cần bằng câu lệnh **ufw**

Sử dụng các câu lệnh **ufw** để chặn mọi thứ mà không liên quan đến ứng dụng blockchain.
Hãy làm theo các bước sau:

```bash
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 60000:61000/udp
sudo ufw allow 26657/tcp
sudo ufw allow 1317/tcp
sudo ufw allow 8080/tcp
sudo ufw --force enable
```
