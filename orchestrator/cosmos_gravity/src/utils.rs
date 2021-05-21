use deep_space::error::CosmosGrpcError;
use deep_space::Contact;
use std::time::{Duration, Instant};
use tokio::time::sleep;

pub async fn wait_for_cosmos_online(contact: &Contact, timeout: Duration) {
    let start = Instant::now();
    while let Err(CosmosGrpcError::NodeNotSynced) = contact.wait_for_next_block(timeout).await {
        sleep(Duration::from_secs(1)).await;
        if Instant::now() - start > timeout {
            panic!("Cosmos node has not come online during timeout!")
        }
    }
    contact.wait_for_next_block(timeout).await.unwrap();
    contact.wait_for_next_block(timeout).await.unwrap();
    contact.wait_for_next_block(timeout).await.unwrap();
}
