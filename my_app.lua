box.cfg{listen = 3301}
box.schema.user.passwd('Gx5!')
box.schema.user.grant('guest','read,write,execute','universe', nil, {if_not_exists = true})

box.schema.space.create('WALLETS', {if_not_exists = true, engine = ‘vinyl’})

box.schema.space.create('transactions', {if_not_exists = true, engine = ‘vinyl’})


box.space.WALLETS:format({
  {name = 'ID', type ='string'},
  {name = 'HASH', type='string'},
  {name = 'MEMBER_ID', type='string'},
  {name = 'AMOUNT', type = 'unsigned'},
  {name = 'CREATED_AT', type = 'unsigned'}
  })

box.space.transactions:format({
 {name='id', type ='string'},
 {name = 'from_wallet_id', type = 'string'},
 {name = 'to_wallet_id', type = 'string'},
 {name = 'amount', type = 'number'}
 })

 box.space.WALLETS:create_index('primary', {type = 'tree', parts = {'ID'}}, {if_not_exists = true})
-- box.space.WALLETS:create_index('amount_index', {type = 'tree', unique = false, parts = {'AMOUNT'}}, {if_not_exists = true})
-- box.space.WALLETS:create_index('amount_date', {type = 'tree', unique = false, parts = {'AMOUNT', 'CREATED_AT'}}, {if_not_exists = true})

-- box.space.transactions:create_index('primary', {type = 'tree', parts = {'id'}}, {if_not_exists = true})


function sum(a, b)
    return a + b
end

function transfer(from_wallet_id, to_wallet_id, amount)
    box.begin()
    box.space.WALLETS:update(from_wallet_id, {{ '-', 4, amount}})
    box.space.WALLETS:update(to_wallet_id, {{ '+', 4, amount}})

    box.commit()
end