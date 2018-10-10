'use strict';

/**
* Place order
* @param{org.workshop.escrow.OrderItem} orderItem
* @transaction
*/

async function orderItem(tx){
  console.log('orderItem');

  const NS ='org.workshop.escrow';
  const factory = getFactory();
  
  let item = tx.item;
  let buyer = getCurrentParticipant();
  
  const order = factory.newResource(NS, 'Order', tx.orderId);
  order.item = factory.newRelationship(NS,'Item', item.getIdentifier());
  order.buyer = factory.newRelationship(NS, 'User', buyer.getIdentifier());
  
  let price = item.price;
    
  if(buyer.tokenBalance >= price){
    buyer.tokenBalance = buyer.tokenBalance - price;   
          	
    const escrowAccount = factory.newResource(NS,'EscrowAccount', tx.escrowAccount.getIdentifier());
    order.escrowAccount = escrowAccount;
    order.status = 'PENDING';
    escrowAccount.tokenBalance = price;
    escrowAccount.buyer = factory.newRelationship(NS, 'User', buyer.getIdentifier());


    const orderRegistry = await getAssetRegistry(order.getFullyQualifiedType());
    await orderRegistry.add(order);

    const escrowRegistry = await getAssetRegistry(escrowAccount.getFullyQualifiedType());
    await escrowRegistry.add(escrowAccount);

    const participantRegistry = await getParticipantRegistry('org.workshop.escrow.User');
    await participantRegistry.update(buyer);
    
  }else{
    throw new Error("insufficient funds");
  }   	
}   

/**
* complete order
* @param{org.workshop.escrow.CompleteOrder} completeOrder
* @transaction
*/

async function completeOrder(tx){
  console.log('Completing Order');
  
  let order = tx.order;
  let balance = order.escrowAccount.tokenBalance;
  let seller = null;
  
                       
  
  if (order.status = 'PENDING'){
  	order.status = 'COMPLETED';
    seller = order.item.owner;
    seller.tokenBalance += balance;
    
    order.item.owner = order.buyer;
   
    const orderRegistry = await getAssetRegistry(order.getFullyQualifiedType());
    await orderRegistry.update(order);
    
    const itemRegistry = await getAssetRegistry(order.item.getFullyQualifiedType());
    await itemRegistry.update(order.item);
    
    const sellerRegistry = await getParticipantRegistry('org.workshop.escrow.User');
    await sellerRegistry.update(seller);
    
    const escrowRegistry = await getAssetRegistry(order.escrowAccount.getFullyQualifiedType());
    await escrowRegistry.remove(order.escrowAccount);
    
  } else {
  	throw new Error('No orders found');
  }
}

/**
* complete order
* @param{org.workshop.escrow.CancelOrder} cancelOrder
* @transaction
*/

async function cancelOrder(tx){
  console.log('Cancelling Order');
  
  let order = tx.order;
  let balance = order.escrowAccount.tokenBalance;
  let buyer = null;
  
                       
  
  if (order.status = 'PENDING'){
  	order.status = 'CANCELLED';
    buyer = order.buyer;
    buyer.tokenBalance += balance;
    
  
   
    const orderRegistry = await getAssetRegistry(order.getFullyQualifiedType());
    await orderRegistry.update(order);
    
    const sellerRegistry = await getParticipantRegistry('org.workshop.escrow.User');
    await sellerRegistry.update(buyer);
    
    const escrowRegistry = await getAssetRegistry(order.escrowAccount.getFullyQualifiedType());
    await escrowRegistry.remove(order.escrowAccount);
    
  } else {
  	throw new Error('No orders found');
  }
}

/**
* Buy tokens
* @param{org.workshop.escrow.BuyTokens} buyTokens
* @transaction
*/
 
async function buyTokens(tx){
  	console.log('buyTokens');
     
    let amount = tx.amount;
    let bank = null;
    let buyer = getCurrentParticipant();
    
    if (tx.amount > 0){
      	bank = tx.bank;
      	
      	if(bank.tokenBalance >= amount){
          bank.tokenBalance -= amount;
          buyer.tokenBalance += amount;
		
          //update buyer balance
          const buyerRegistry = await getParticipantRegistry('org.workshop.escrow.User');
          await buyerRegistry.update(buyer);
          
		  //update bank balance
          const bankRegistry = await getParticipantRegistry('org.workshop.escrow.Bank');
          await bankRegistry.update(bank);
        }
    } else {
    	throw new Error("No sufficient funds");
    }
  }


