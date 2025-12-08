package transaction

import (
	"fmt"

	"gorm.io/gorm"
)

/**

	假设有两个表：
	accounts 表（包含字段 id 主键， balance 账户余额）和
	transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
	要求 ：
	编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，
	向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。

**/

type Account struct {
	ID       int `gorm:"primary_key"`
	Balance  int
	UserName string
}

type Transaction struct {
	ID            int `gorm:"primary_key"`
	FromAccountID int // 转出账户ID
	ToAccountID   int // 转入账户ID
	Amount        int // 转账金额
}

func Run(db *gorm.DB) {
	db.AutoMigrate(&Account{}, &Transaction{})

	// 创建账户 A B
	//db.Create(&Account{UserName: "A", Balance: 1000})
	//db.Create(&Account{UserName: "B", Balance: 1000})

	//转出账户 ID
	fromID := 1
	//转入账户 ID
	toID := 2
	//转账金额
	amount := 100
	db.Transaction(func(tx *gorm.DB) error {
		var fromAccount Account
		tx.First(&fromAccount, fromID)
		if fromAccount.Balance < amount {

			return fmt.Errorf("余额不足")
		}
		var toAccount Account
		if err := tx.First(&toAccount, toID).Error; err != nil {
			return fmt.Errorf("获取转入账户信息失败", err)
		}
		if err := tx.Model(&fromAccount).Update("balance", fromAccount.Balance-amount).Error; err != nil {
			return fmt.Errorf("更新转出账户信息失败", err)
		}
		if err := tx.Model(&toAccount).Update("balance", toAccount.Balance+amount).Error; err != nil {
			return fmt.Errorf("更新转入账户信息失败", err)
		}
		//记录交易
		t1 := Transaction{FromAccountID: fromID, ToAccountID: toID, Amount: amount}
		if err := tx.Create(&t1).Error; err != nil {
			return fmt.Errorf("记录交易失败", err)
		}
		return nil

	})

}
