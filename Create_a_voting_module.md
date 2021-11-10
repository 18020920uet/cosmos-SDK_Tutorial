# Create a voting module
Hướng dẫn tạo ứng dụng blockchain thăm dò ý kiến

---
## Polling app

Hướng dẫn này hướng dẫn từ việc tạo một ứng dụng blockchain, thêm và sửa types cho transaction, sửa message và thiết kế front-end

Nội dung:
+ Tạo một ứng dụng blockchain thăm dò ý kiến đơn giản (poll)
+ Thiết kế ứng dụng front-end mà người dùng có thể đăng nhập, tạo poll, vote, và xem các kết quả vote
+ Thêm logic để yêu cầu tiền để có thể tạo poll transaction
+ Tùy chỉnh REST endpoint
+ Tùy chỉnh giao dịch CLI
+ Thêm các components cho front-end
+ Thêm bank keeper cho module
+ Tùy chỉnh message

Yêu cầu:
- Starport v0.17.0
```bash
curl https://get.starport.network/starport@v0.17.0! | bash
```

### Voting App Goals
Tạo một ứng dụng blockchain thăm dò ý kiến với các module bình chọn. Người dùng có thể:
+ Đăng nhập
+ Tạo polls
+ Vote
+ Xem kết quả vote

Thiết kể ứng dụng để:
- Thu phí yêu cầu
  - Tạo giao dịch thăm dò ý kiến: 200 tokens
  - Bình chọn: free

- Hạn chế bình chọn cho người dùng đã đăng nhập

### Build your Blockchain App

Dùng Starport để scaffold(tạo khung) ứng dụng blockchain và module bình chọn

#### Build the new blockchain

```bash
starport scaffold chain github.com/cosmonaut/voter
```

**app**: Bao gồm tất cả cá file mà liên kết giữa các phần của ứng dụng lại với nhau

**cmd**: Có trách nghiệm cho voterd (server, chạy ứng dụng) và phần tương tác với ứng dụng

**proto**: chứ các loại protobuf

**vue**: bao gồm web UI

**x**: là Cosmos-SDK cho ```voter```

#### Launch the voter app

Khởi động ứng dụng bằng:
```bash
starport chain serve
```

### Add a Poll Transaction

Bước đầu tiên là thêm các type và thêm các chức năng CRUD cho các type đó.

Có 2 type là: **polls** và **votes**

#### Add the poll type
**Poll**: bao gồm ```title``` và ```options```

Dùng scaffold list để tạo type và CRUD của **poll**

```bash
starport scaffold list poll title options
```

![Starport scaffold list poll](/images/Create-a-voting-module/Create-poll.png)

#### View the Front-end User Interface

Khởi dộng front-end trong _vue_

```bash
npm install
npm run serve
```
Kiểm tra port 8080 xem tình trạng

> http:/{server}:{port}

Ở đây dùng localhost nên sẽ vào localhost:8080 để kiểm tra
> htpp:/localhost:8080

![Starport vue](/images/Create-a-voting-module/vue.png)

#### Sign in as Alice

Sau khi chạy ```starport chain serve```, nnemonic của 2 tài khoản Bob và Alice sẽ được show ra như dưới đây

![Starport Alice and Bob](/images/Create-a-voting-module/Alice-and-Bob.png)

Trên màn hình Access wallet rồi "Import existing wallet", sau đó nhập mnemonic vào. Rồi chỉnh thông tin tài khoản với Wallet Name và Password

![Starport Import Existing Wallet](/images/Create-a-voting-module/Mnemonic.png)

Nhập thông tin wallet sau khi nhập mnemonic

![Nhập thông tin wallet](/images/Create-a-voting-module/Sign-in-information.png)

Sau khi đăng nhập

![Signed in](/images/Create-a-voting-module/Signed-in.png)

### View the Poll Type

Để xem các giao dịch poll thì ấn vào **Custome type**

Thử tạo 1 poll và lưu lại
![Create example poll](/images/Create-a-voting-module/Create-example-poll.png)

Cấu trúc của poll chưa thể vote được, phải chuyển hóa ```options``` thành array để có thể chọn được.
Phần này có thể tìm hiểu thêm, chỉnh sửa qua ở ```starport scaffold list```

### Modify the Protobuffer Types

Để có thể có nhiều ```options``` trong 1 poll thì phải chuyển kiểu dữ liệu ```string options``` trong cấu trúc của ```Poll```

```proto3
message Poll {
  string creator = 1;
  uint64 id = 2;
  string title = 3;
  repeated string options = 4;
}
```

Trong file _proto/voter/**tx.proto**_ cập nhật type CRUD của poll bằng cách chuyển cấu trúc của **MsgCreatePoll** từ ```string options``` thành ```repeated string options``` tương tụ **MsgUpdatePoll**


```proto3
message MsgCreatePoll {
  string creator = 1;
  string title = 2;
  repeated string options = 4;
}

```

```proto3
message MsgUpdatePoll {
  string creator = 1;
  uint64 id = 2;
  string title = 3;
  repeated string options = 4;
}
```

### Modify the Poll Transaction Message

File _x/voter/types/**messages_poll.go**_ có định nghĩa message khi tạo poll là hàm **NewMsgCreatePoll**

chuyển kiểu dữ liệu ```options string``` thành ```options []string``` ở **NewMsgCreatePoll** và **NewMsgUpdatePoll**

```go
func NewMsgCreatePoll(creator string, title string, options []string) *MsgCreatePoll {
  ...
}
...
```

```go
func NewMsgUpdatePoll(creator string, id uint64, title string, options []string) *MsgUpdatePoll {
  ...
}
...
```

### About the Poll Keeper

Để ghi bất cứ thứ gì vào một blockchain hoặc tác động vào các giao dịch trạng thái, client sẽ gửi HTTP POST request để có thể tạo ra sự thay đổi đó.

Hanlder tạo ra một unsigned transaction (hay có thể gọi là giao dịch chưa được đánh dấu) bao gôm 1 mảng các **message**.
Client sẽ xác nhận các giao dịch đó và gửi đến ```http://localhost:1317/txs```. Ứng dụng sẽ xử lý giao dịch bằng cách gửi các **message** đến đúng handler, trong trường hợp này là _x/voter/**handler.go**_.

```go
// x/voter/handler.go
func NewHandler(k keeper.Keeper) sdk.Handler {
  msgServer := keeper.NewMsgServerImpl(k)

  return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
    ctx = ctx.WithEventManager(sdk.NewEventManager())

    switch msg := msg.(type) {
    // this line is used by starport scaffolding # 1

    case *types.MsgCreatePoll:
      res, err := msgServer.CreatePoll(sdk.WrapSDKContext(ctx), msg)
      return sdk.WrapServiceResult(ctx, res, err)

    case *types.MsgUpdatePoll:
      res, err := msgServer.UpdatePoll(sdk.WrapSDKContext(ctx), msg)
      return sdk.WrapServiceResult(ctx, res, err)
    ...
  }
}
```

Handler sẽ gọi đến hàm **CreatePoll** trong _x/voter/keeper/**poll.go**_ để ghi dữ liệu vào trong **store**.

```go
// x/voter/keeper/msg_server_poll.go
...
func (k msgServer) CreatePoll(goCtx context.Context, msg *types.MsgCreatePoll) (*types.MsgCreatePollResponse, error) {
  ctx := sdk.UnwrapSDKContext(goCtx)

  var poll = types.Poll{
    Creator: msg.Creator,
    Title:   msg.Title,
    Options: msg.Options,
  }

  id := k.AppendPoll(
    ctx,
    poll,
  )

  return &types.MsgCreatePollResponse{
    Id: id,
  }, nil
}
```

```go
// x/voter/keeper/poll.go
...
func (k Keeper) AppendPoll(
  ctx sdk.Context,
  poll types.Poll,
) uint64 {
  // Create the poll
  count := k.GetPollCount(ctx)

  // Set the ID of the appended value
  poll.Id = count

  store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PollKey))
  appendedValue := k.cdc.MustMarshalBinaryBare(&poll)
  store.Set(GetPollIDBytes(poll.Id), appendedValue)

  // Update poll count
  k.SetPollCount(ctx, count+1)

  return count
}
```

#### Modify the CLI Transaction

Một người sử dụng có thể tương tác với ứng dụng thông qua CLI, CLI đuọc định nghĩa ở _x/voter/client/cli/**tx_poll.go**_

Ví dụ về câu lệnh CLI sẽ sử dụng

```bash
voterd tx voter create-poll "Text editors" "Emacs" "Vim" --from alice
```

Trong hàm **CmdCreatePoll** thay đổi thành

```go
func CmdCreatePoll() *cobra.Command {
  cmd := &cobra.Command{
    Use:   "create-poll [title] [options]",
    Short: "Create a new poll",
    Args:  cobra.MinimumNArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
      argsTitle, err := cast.ToStringE(args[0])
      if err != nil {
        return err
      }
      argsOptions := args[1:len(args)]

      clientCtx, err := client.GetClientTxContext(cmd)
      if err != nil {
        return err
      }

      msg := types.NewMsgCreatePoll(clientCtx.GetFromAddress().String(), argsTitle, argsOptions)
      if err := msg.ValidateBasic(); err != nil {
        return err
      }
      return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
    },
  }

  flags.AddTxFlagsToCmd(cmd)

  return cmd
}
```

tương tụ hàm **CmdUpdatePoll**

```go
func CmdUpdatePoll() *cobra.Command {
  cmd := &cobra.Command{
    Use:   "update-poll [id] [title] [options]",
    Short: "Update a poll",
    Args:  cobra.MinimumNArgs(3),
    RunE: func(cmd *cobra.Command, args []string) error {
      id, err := strconv.ParseUint(args[0], 10, 64)
      if err != nil {
        return err
      }

      argsTitle, err := cast.ToStringE(args[1])
      if err != nil {
        return err
      }

      argsOptions := args[2:len(args)]

      clientCtx, err := client.GetClientTxContext(cmd)
      if err != nil {
        return err
      }

      msg := types.NewMsgUpdatePoll(clientCtx.GetFromAddress().String(), id, argsTitle, argsOptions)
      if err := msg.ValidateBasic(); err != nil {
        return err
      }
      return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
    },
  }

  flags.AddTxFlagsToCmd(cmd)

  return cmd
}
```

Rest lại ứng dụng với lệnh

```bash
starport chain serve --reset-once
```

![Cập nhật về types của Poll](/images/Create-a-voting-module/Client-after-update-poll-type.png)

### Add the Votes

Tạo type **vote**, với pollID là ID của poll, option là lựa chọn của người dùng

```bash
starport scaffold list vote pollID option
```

sau khi chạy xong, reset lại ứng dụng để cập nhật thay đổi

```bash
starport serve --reset-once
```

![Sau khi thêm vote](/images/Create-a-voting-module/Client-after-add-vote.png)

Chú ý: mỗi lần rest trạng thái của ứng dụng, dữ liệu được tạo từ trạng thái trước sẽ không được lưu

### Front-end Application

Starpot sẽ tự động tạo một ứng dụng front-end đơn giản sử dụng framework Vue và Vuex để quản lý state. Vì ứng dụng blockchain hoạt động quai HTTP API nên có thể dùng ngôn ngữ, framework khác để viết lại client tùy ý.

Một số folder của ứng dụng

- _vue/src/views_ và _vue/src/components_: Template của các page
- _vue/src/store_: Xử lý gửi yêu cầu giao dịch và nhận dữ liệu từ ứng dụng blockchain
- _@tendermint/vue_ chứa các phần về UI (form, button, ...), đồng thời thư mục này bao gồm các file định nghĩa protobuffer mà được định nghĩa trong thư _vue/src/store/generated/cosmonaut/voter/cosmonaut.voter.voter_
- _vue/src/store/generated/cosmonaut/voter/cosmonaut.voter.voter/index.js_: đã tự động sinh ra các transactions (giao dịch) **MsgCreatePoll**, **MsgUpdatePoll**, **MsgDeletePoll** mà sử dụng thư viện [ComJS](https://github.com/cosmwasm/cosmjs) để xử lý wallets (các ví), creating (tạo), signing (đăng nhập), và gửi transactions (các giao dịch) và định nghĩa Vuex store


### Add the Voter Module Component to the Front End

Sửa file _vue/src/views/**Types.vue**_

```html
<template>
  <div>
    <div class="container">
      <!-- this line is used by starport scaffolding # 4 -->
      <!-- <SpType modulePath="cosmonaut.voter.voter" moduleType="Vote"  /> -->
      <!-- <SpType modulePath="cosmonaut.voter.voter" moduleType="Poll"  /> -->
    <h2>
    Voter Module
    </h2>
    <PollForm />
    <PollList />
    </div>
  </div>
</template>

<script>
    import PollForm from "../components/PollForm";
    import PollList from "../components/PollList";
    export default {
      name: 'Types',
      components: { PollForm, PollList },
    }
</script>
```

tiếp tục sẽ tạo 2 components là **PollForm** và **PollList**

#### Create the PollForm Component

Tạo file _vue/src/componets/**PollForm.vue**

```html
<template>
  <div>
    <div class="sp-voter__main sp-box sp-shadow sp-form-group">
        <form class="sp-voter__main__form">
          <div class="sp-voter__main__rcpt__header sp-box-header">
            Create a Poll
          </div>

          <input class="sp-input" placeholder="Title" v-model="title" />
          <div v-for="(option, index) in options" v-bind:key="'option' + index">
            <input class="sp-input" placeholder="Option" v-model="option.title" />
          </div>
          <sp-button @click="add">+ Add option</sp-button>
          <sp-button @click="submit">Create poll</sp-button>
        </form>
    </div>
  </div>
</template>
<script>
export default {
  name: "PollForm",
  data() {
    return {
      title: "",
      options: [{
        title: "",
      }],
    };
  },
  computed: {

    currentAccount() {
      if (this._depsLoaded) {
        if (this.loggedIn) {
          return this.$store.getters['common/wallet/address']
        } else {
          return null
        }
      } else {
        return null
      }
    },
    loggedIn() {
      if (this._depsLoaded) {
        return this.$store.getters['common/wallet/loggedIn']
      } else {
        return false
      }
    }
  },
  methods: {
    add() {
      this.options = [...this.options, { title: "" }];
    },
    async submit() {
      const value = {
        creator: this.currentAccount,
        title: this.title,
        options: this.options.map((o) => o.title),
      };
      await this.$store.dispatch("cosmonaut.voter.voter/sendMsgCreatePoll", {
        value,
        fee: [],
      });
    },
  },
};
</script>
```

Sau đó có thể thử bằng cách refresh + tạo poll mới vào **http://localhost:1317/voter/poll** để kiểm tra kết quả

#### Create the Poll List Component

Tạo file _vue/src/componets/**PollList.vue**

```html
<template>
  <div>
    <h2> List of Polls </h2>
    <div v-for="poll in polls" v-bind:key="'poll' + poll.id">
      <b> {{poll.id}}. {{ poll.title }} </b>
      <AppRadioItem
        @click="submit(poll.id, option)"
        v-for="option in poll.options"
        v-bind:key="option"
        :value="option"
      />
      <AppText type="subtitle">Results: {{ results(poll.id) }}</AppText>
    </div>
  </div>
</template>
<style>
.option-radio > .button {
  height: 40px;
  width: 50%;
}
</style>
<script>
import AppRadioItem from "./AppRadioItem";
import AppText from "./AppText";
import { countBy } from "lodash";

export default {
  components: { AppText, AppRadioItem },
  data() {
    return {
      selected: "",
    };
  },
  computed: {

    currentAccount() {
      if (this._depsLoaded) {
        if (this.loggedIn) {
          return this.$store.getters['common/wallet/address']
        } else {
          return null
        }
      } else {
        return null
      }
    },
    loggedIn() {
      if (this._depsLoaded) {
        return this.$store.getters['common/wallet/loggedIn']
      } else {
        return false
      }
    },
    polls() {
      return (
        this.$store.getters["cosmonaut.voter.voter/getPollAll"]({
          params: {}
        })?.Poll ?? []
      );
    },
    votes() {
      return (
        this.$store.getters["cosmonaut.voter.voter/getVoteAll"]({
          params: {}
        })?.Vote ?? []
      );
    },
  },
  methods: {
    results(id) {
      const results = this.votes.filter((v) => v.pollID === id);
      return countBy(results, "option");
    },
    async submit(pollID, option) {

      const value = { creator: this.currentAccount, pollID, option };
      await this.$store.dispatch("cosmonaut.voter.voter/sendMsgCreateVote", {
        value,
        fee: [],
      });
      await this.$store.dispatch("cosmonaut.voter.voter/QueryPollAll", {
        options: { subscribe: true, all: true },
        params: {},
      });
    },
  },
};
</script>
```

#### Add the Options Component

Thêm vào file _vue/src/components/**AppRadioItem.vue**_

```html
<template>
  <div>
    <button class="button">{{ value }}</button>
  </div>
</template>

<style scoped>
.button {
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.1);
  padding: 0.75rem 1rem;
  margin-bottom: 0.5rem;
  border-radius: 6px;
  user-select: none;
  cursor: pointer;
  transition: all 0.1s;
  font-weight: 500;
  outline: none;
  border: none;
  background: rgba(0, 0, 0, 0.01);
  width: 100%;
  font-family: inherit;
  text-align: left;
  font-size: 1rem;
  font-weight: inherit;
}
.button:hover {
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.1),
    0 1px 5px -1px rgba(0, 0, 0, 0.1);
}
.button:focus {
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.2),
    0 1px 5px -1px rgba(0, 0, 0, 0.1);
}
.button:active {
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.07);
  color: rgba(0, 0, 0, 0.7);
  transform: scale(0.998);
}
</style>
<script>
export default {
  props: {
    value: "",
  },
};
</script>
```

#### Add the Poll List Text Component

Sửa file _vue/src/components/**AppText.vue**_

```html
<template>
  <div>
    <div :class="[`${type}`]">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.h1 {
  font-size: 2rem;
  font-weight: 800;
  letter-spacing: 0.03em;
  margin-bottom: 2rem;
}
.h2 {
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.5rem;
  margin-top: 2rem;
}
.subtitle {
  color: rgba(0, 0, 0, 0.5);
  font-size: 0.85rem;
}
</style>
<script>
export default {
  props: {
    type: {
      default: "p1",
    },
  },
};
</script>
```
#### Update the Front-end App

Cập nhật file _vue/src/**App.vue**_ để lấy dữ liệu vote

```html
<template>
	<div v-if="initialized">
		<SpWallet ref="wallet" v-on:dropdown-opened="$refs.menu.closeDropdown()" />
		<SpLayout>
			<template v-slot:sidebar>
				<Sidebar />
			</template>
			<template v-slot:content>
				<router-view />
			</template>
		</SpLayout>
	</div>
</template>

<style>
body {
	margin: 0;
}
</style>

<script>
import './scss/app.scss'
import '@starport/vue/lib/starport-vue.css'
import Sidebar from './components/Sidebar'

export default {
    components: {
        Sidebar
    },
    data() {
        return {
            initialized: false
        }
    },
    computed: {
        hasWallet() {
            return this.$store.hasModule([ 'common', 'wallet'])
        }
    },
    async created() {
        await this.$store.dispatch('common/env/init')
        this.initialized = true
        await this.$store.dispatch("cosmonaut.voter.voter/QueryPollAll",{options:{subscribe:true, all:true},params:{}})
        await this.$store.dispatch("cosmonaut.voter.voter/QueryVoteAll",{options:{subscribe:true, all:true},params:{}})
    },
    errorCaptured(err) {
        console.log(err)
        return false
    }
}
</script>
```

![Sau khi cập nhật front-end](/images/Create-a-voting-module/After-update-fe.png)
Ở đây có thể thấy số lượng vote, dưới mỗi poll nhưng hiện tại mỗi người có thể vote nhiều lần, sau khi cập nhật phần dưới thì sẽ làm đươc mỗi người chỉ bình chọn được 1 lần

### Access the API

Phần này nói về cách truy cập API và lấy dữ liệu từ server, để tìm hiểu thêm có thể đọc tại [đây](https://tutorials.cosmos.network/voter/#access-the-api)

### Limit to One Vote per User

Tùy chỉnh mỗi người dùng chỉ bình chọn một vote trong 1 poll

Phần logic để truy cập các giao dịch được để trong phần **keeper**. Để xem các giao dịch về việc bình chọn xem file _x/voter/keeper/**msg_server_vote.go**_. Tùy chỉnh hàm **CreateVote** để kiểm tra xem tài khoản đã bình chọn chưa và trả lại lỗi nếu người đó bỏ phiếu lại

```go
func (k msgServer) CreateVote(goCtx context.Context, msg *types.MsgCreateVote) (*types.MsgCreateVoteResponse, error) {
  ctx := sdk.UnwrapSDKContext(goCtx)

  var vote = types.Vote{
    Creator: msg.Creator,
    PollID:  msg.PollID,
    Option:  msg.Option,
  }

  // Get all existing votes
  voteList := k.GetAllVote(ctx)
  for _, existingVote := range voteList {
    // Check if the account has already voted on this PollID
    if existingVote.Creator == msg.Creator && existingVote.PollID == msg.PollID {
      // Return an error when a vote has been cast by this account on this PollID
      return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Vote already casted.")
    }
  }

  id := k.AppendVote(
    ctx,
    vote,
  )

  return &types.MsgCreateVoteResponse{
    Id: id,
  }, nil
}
```

Đã kiểm tra và xác nhận cùng 1 người dùng chỉ vote được poll 1 lần duy nhất, không có khả năng vote lại nữa

![Cập nhật vote 1 lần](/images/Create-a-voting-module/Vote-one-time.png)

### Introducing a Fee for Creating Polls

Yêu cầu đề ra là mỗi lần tạo poll sẽ mất 200 tokens

Mỗi người dùng cuối có tokens dư trong tài khoản. Việc phải làm là gửi coins từ tài khoản của người dùng đến module **account** trước khi tạo poll

### Add the Bank Keeper to the Voter Module

Load expected_keepers trong file _x/voter/types/**expected_keepers.go**_ để định nghĩa các chức năng của **bank** mà muốn sử dụng trong module

```go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
    SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
    MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
    BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
    SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
    SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}
```

Thêm **bankKeeper** trong _x/voter/keeper/**keeper.go**_, thêm vào type và hàm NewKeeper như sau

```go
type (
    Keeper struct {
        cdc        codec.Marshaler
        storeKey   sdk.StoreKey
        memKey     sdk.StoreKey
        bankKeeper types.BankKeeper
    }
)

func NewKeeper(cdc codec.Marshaler, storeKey, memKey sdk.StoreKey, bankKeeper types.BankKeeper) *Keeper {
    return &Keeper{
        cdc:        cdc,
        storeKey:   storeKey,
        memKey:     memKey,
        bankKeeper: bankKeeper,
    }
}
```

Cuối cùng là thêm module **bank** vào trong hàm khai báo NewKeeper trong _app/**app.go**_

```go
app.VoterKeeper = *votermodulekeeper.NewKeeper(
  appCodec,
  keys[votermoduletypes.StoreKey],
  keys[votermoduletypes.MemStoreKey],
  app.BankKeeper,
)
```

Bước tiếp theo là định nghĩa lại hàm lúc tạo poll sao cho yêu cầu tiền để có thể tạo

### Modify the Create Poll Message with the Price

Sửa file _x/voter/keeper/**msg_server_poll.go**_

```go
package keeper

import (
    "context"
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "github.com/tendermint/tendermint/crypto"
    "github.com/cosmonaut/voter/x/voter/types"
)

func (k msgServer) CreatePoll(goCtx context.Context, msg *types.MsgCreatePoll) (*types.MsgCreatePollResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
    feeCoins, err := sdk.ParseCoinsNormalized("200token")
    if err != nil {
        return nil, err
    }

    creatorAddress, err := sdk.AccAddressFromBech32(msg.Creator)
    if err != nil {
        return nil, err
    }
    if err := k.bankKeeper.SendCoins(ctx, creatorAddress, moduleAcct, feeCoins); err != nil {
        return nil, err
    }

    var poll = types.Poll{
    Creator: msg.Creator,
    Title:   msg.Title,
    Options: msg.Options,
  }

  id := k.AppendPoll(
    ctx,
    poll,
  )

  return &types.MsgCreatePollResponse{
    Id: id,
  }, nil
}
```

![Token trước](/images/Create-a-voting-module/Previous-tokens.png)

Token sau khi tạo 1 Poll

![Token sau](/images/Create-a-voting-module/After-tokens.png)

Tính phí và thu phí trước khi gọi hàm **k.AppendPoll** nên nếu người dùng không có đủ tiền để tạo ứng dụng thì sẽ không tạ được poll và sẽ thông báo lại lỗi.
