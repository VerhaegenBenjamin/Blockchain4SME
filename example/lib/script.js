/**
 * This is the transaction processor function description
 * @param {org.example.workshop.TradeVehicle} tradeVehicle - the TradeVehicle transaction
 * @transaction
 */

async function tradeVehicle(tradeVehicle) {
  console.log("TradeVehicle transaction");
  
  let vehicle = tradeVehicle.vehicle;
  let newOwner = tradeVehicle.newOwner;
  
  const factory = getFactory();
  
  if(vehicle.vehicleStatus == 'NEW'){
  	vehicle.owner = newOwner;
  } else {
  	throw new Error('Car must be new');
  }
  
  let assetRegistry = await getAssetRegistry('org.example.workshop.Vehicle');
  await assetRegistry.update(vehicle);
  
  let event = factory.newEvent('org.example.workshop', 'TradeVehicleEvent');
  event.vehicle = vehicle;
  event.newOwner = newOwner;
  emit(event);
}

/**
 * Remove used vehicles
 * @param {org.example.workshop.RemoveUsedVehicles} remove - the RemoveUsedVehicles transaction
 * @transaction
 */

async function removeUsedVehicles(remove) {
  let results = await query('selectUsedVehicles');
  const factory = getFactory();
  
  for (let n = 0; n < results.length; n++){
  	let vehicle = results[n];
    
    let assetRegistry = await getAssetRegistry('org.example.workshop.Vehicle');
    await assetRegistry.remove(vehicle);
    
    let removeNotification = factory.newEvent('org.example.workshop', 'RemoveNotification');
    removeNotification.vehicle = vehicle;
    emit(removeNotification);
  }
}




